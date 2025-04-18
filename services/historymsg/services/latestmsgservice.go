package services

import (
	"context"
	"fmt"
	"im-server/commons/bases"
	"im-server/commons/caches"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/historymsg/storages"
	"time"
)

var latestMsgCache *caches.LruCache
var latestMsgLocks *tools.SegmentatedLocks

func init() {
	latestMsgCache = caches.NewLruCacheWithAddReadTimeout("latestmsg_cache", 10000, nil, 5*time.Minute, 5*time.Minute)
	latestMsgLocks = tools.NewSegmentatedLocks(128)
}

type LatestMsgItem struct {
	key           string
	LatestMsgId   string
	LatestMsgSeq  int64
	LatestMsgTime int64
}

func (item *LatestMsgItem) Update(msg *pbobjs.DownMsg) {
	if msg.MsgTime > item.LatestMsgTime {
		item.LatestMsgId = msg.MsgId
		item.LatestMsgSeq = msg.MsgSeqNo
		item.LatestMsgTime = msg.MsgTime
	}
}

func getLatestMsgCacheKey(appkey, converId string, channelType pbobjs.ChannelType) string {
	return fmt.Sprintf("%s_%s_%d", appkey, converId, channelType)
}
func GetLatestMsg(ctx context.Context, converId string, channelType pbobjs.ChannelType) *LatestMsgItem {
	appkey := bases.GetAppKeyFromCtx(ctx)
	key := getLatestMsgCacheKey(appkey, converId, channelType)
	if val, exist := latestMsgCache.Get(key); exist {
		return val.(*LatestMsgItem)
	} else {
		lock := latestMsgLocks.GetLocks(key)
		lock.Lock()
		defer lock.Unlock()

		if val, exist := latestMsgCache.Get(key); exist {
			return val.(*LatestMsgItem)
		} else {
			item := &LatestMsgItem{
				key:           key,
				LatestMsgId:   "",
				LatestMsgSeq:  0,
				LatestMsgTime: 0,
			}
			if channelType == pbobjs.ChannelType_Private {
				storage := storages.NewPrivateHisMsgStorage()
				latestMsg, err := storage.QryLatestMsg(appkey, converId)
				if err == nil && latestMsg != nil {
					item.LatestMsgId = latestMsg.MsgId
					item.LatestMsgSeq = latestMsg.MsgSeqNo
					item.LatestMsgTime = latestMsg.SendTime
				}
			} else if channelType == pbobjs.ChannelType_Group {
				storage := storages.NewGroupHisMsgStorage()
				latestMsg, err := storage.QryLatestMsg(appkey, converId)
				if err == nil && latestMsg != nil {
					item.LatestMsgId = latestMsg.MsgId
					item.LatestMsgSeq = latestMsg.MsgSeqNo
					item.LatestMsgTime = latestMsg.SendTime
				}
			} else if channelType == pbobjs.ChannelType_System {
				storage := storages.NewSystemHisMsgStorage()
				latestMsg, err := storage.QryLatestMsg(appkey, converId)
				if err == nil && latestMsg != nil {
					item.LatestMsgId = latestMsg.MsgId
					item.LatestMsgSeq = latestMsg.MsgSeqNo
					item.LatestMsgTime = latestMsg.SendTime
				}
			} else if channelType == pbobjs.ChannelType_GroupCast {
				storage := storages.NewGrpCastHisMsgStorage()
				latestMsg, err := storage.QryLatestMsg(appkey, converId)
				if err == nil && latestMsg != nil {
					item.LatestMsgId = latestMsg.MsgId
					item.LatestMsgSeq = latestMsg.MsgSeqNo
					item.LatestMsgTime = latestMsg.SendTime
				}
			} else if channelType == pbobjs.ChannelType_BroadCast {
				storage := storages.NewBrdCastHisMsgStorage()
				latestMsg, err := storage.QryLatestMsg(appkey, converId)
				if err == nil && latestMsg != nil {
					item.LatestMsgId = latestMsg.MsgId
					item.LatestMsgSeq = latestMsg.MsgSeqNo
					item.LatestMsgTime = latestMsg.SendTime
				}
			}
			latestMsgCache.Add(key, item)
			return item
		}
	}
}

func IsLatestMsg(ctx context.Context, converId string, channelType pbobjs.ChannelType, msgId string, msgTime, msgSeq int64) bool {
	latestMsg := GetLatestMsg(ctx, converId, channelType)
	if msgId == latestMsg.LatestMsgId || msgTime > latestMsg.LatestMsgTime {
		return true
	}
	return false
}

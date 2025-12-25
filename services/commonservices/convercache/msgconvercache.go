package convercache

import (
	"context"
	"fmt"
	"im-server/commons/bases"
	"im-server/commons/caches"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/commonservices/logs"
	"im-server/services/commonservices/msgdefines"
	"im-server/services/historymsg/storages"
	"time"
)

var msgConverCache *caches.LruCache
var msgConverLocks *tools.SegmentatedLocks

func init() {
	msgConverCache = caches.NewLruCacheWithReadTimeout("msgconver_cache", 100000, nil, 10*time.Minute)
	msgConverLocks = tools.NewSegmentatedLocks(512)
}

type MsgConverItem struct {
	cacheKey      string
	LatestMsgSeq  int64
	LatestMsgTime int64
}

func (item *MsgConverItem) GetMsgSeqWithIncr(msgFlag int32) int64 {
	if msgdefines.IsStoreMsg(msgFlag) {
		lock := msgConverLocks.GetLocks(item.cacheKey)
		lock.Lock()
		defer lock.Unlock()
		item.LatestMsgSeq = item.LatestMsgSeq + 1
		return item.LatestMsgSeq
	} else {
		return -1
	}
}

func (item *MsgConverItem) GetMsgSendTime(currentTime int64) int64 {
	lock := msgConverLocks.GetLocks(item.cacheKey)
	lock.Lock()
	defer lock.Unlock()
	if currentTime > item.LatestMsgTime {
		item.LatestMsgTime = currentTime
	} else {
		item.LatestMsgTime = item.LatestMsgTime + 1
	}
	return item.LatestMsgTime
}

func (item *MsgConverItem) GenerateMsgId(converId string, channelType pbobjs.ChannelType, currentTime int64, msgFlag int32) (string, int64, int64) {
	lock := msgConverLocks.GetLocks(item.cacheKey)
	lock.Lock()
	defer lock.Unlock()
	if currentTime > item.LatestMsgTime {
		item.LatestMsgTime = currentTime
	} else {
		item.LatestMsgTime = item.LatestMsgTime + 1
	}
	msgTime := item.LatestMsgTime
	var msgSeq int64 = -1
	if msgdefines.IsStoreMsg(msgFlag) {
		item.LatestMsgSeq = item.LatestMsgSeq + 1
		msgSeq = item.LatestMsgSeq
	}
	msgId := tools.GenerateMsgId(msgTime, int32(channelType), converId)
	return msgId, msgTime, msgSeq
}

func GetMsgConverCache(ctx context.Context, converId, subChannel string, channelType pbobjs.ChannelType) *MsgConverItem {
	appkey := bases.GetAppKeyFromCtx(ctx)
	cacheKey := getMsgConverCacheKey(appkey, converId, subChannel, channelType)
	lock := msgConverLocks.GetLocks(cacheKey)
	lock.Lock()
	defer lock.Unlock()
	if val, ok := msgConverCache.Get(cacheKey); ok {
		return val.(*MsgConverItem)
	} else {
		msgTime, msgSeq := GetLatestMsgTimeSeq(ctx, appkey, converId, subChannel, channelType)
		item := &MsgConverItem{
			cacheKey:      cacheKey,
			LatestMsgTime: msgTime,
			LatestMsgSeq:  msgSeq,
		}
		msgConverCache.Add(cacheKey, item)
		return item
	}
}

func getMsgConverCacheKey(appkey, converId, subChannel string, channelType pbobjs.ChannelType) string {
	return fmt.Sprintf("%s_%d_%s_%s", appkey, channelType, converId, subChannel)
}

func GetLatestMsgTimeSeq(ctx context.Context, appkey, converId, subChannel string, channelType pbobjs.ChannelType) (int64, int64) {
	var msgSeq int64 = 0
	var msgTime int64 = 0
	//从会话查询最新的index
	switch channelType {
	case pbobjs.ChannelType_Private:
		storage := storages.NewPrivateHisMsgStorage()
		latestMsg, err := storage.QryLatestMsg(appkey, converId, subChannel)
		if err != nil {
			logs.WithContext(ctx).Error(err.Error())
		}
		if latestMsg != nil {
			msgSeq = latestMsg.MsgSeqNo
			msgTime = latestMsg.SendTime
		}
	case pbobjs.ChannelType_Group:
		storage := storages.NewGroupHisMsgStorage()
		latestMsg, err := storage.QryLatestMsg(appkey, converId, subChannel)
		if err != nil {
			logs.WithContext(ctx).Error(err.Error())
		}
		if latestMsg != nil {
			msgSeq = latestMsg.MsgSeqNo
			msgTime = latestMsg.SendTime
		}
	case pbobjs.ChannelType_System:
		storage := storages.NewSystemHisMsgStorage()
		latestMsg, err := storage.QryLatestMsg(appkey, converId)
		if err != nil {
			logs.WithContext(ctx).Error(err.Error())
		}
		if latestMsg != nil {
			msgSeq = latestMsg.MsgSeqNo
			msgTime = latestMsg.SendTime
		}
	case pbobjs.ChannelType_GroupCast:
		storage := storages.NewGrpCastHisMsgStorage()
		latestMsg, err := storage.QryLatestMsg(appkey, converId)
		if err != nil {
			logs.WithContext(ctx).Error(err.Error())
		}
		if latestMsg != nil {
			msgSeq = latestMsg.MsgSeqNo
			msgTime = latestMsg.SendTime
		}
	case pbobjs.ChannelType_BroadCast:
		storage := storages.NewBrdCastHisMsgStorage()
		latestMsg, err := storage.QryLatestMsg(appkey, converId)
		if err != nil {
			logs.WithContext(ctx).Error(err.Error())
		}
		if latestMsg != nil {
			msgSeq = latestMsg.MsgSeqNo
			msgTime = latestMsg.SendTime
		}
	}
	/*
		code, resp, err := bases.SyncRpcCall(ctx, "qry_latest_hismsg", converId, &pbobjs.QryLatestMsgReq{
			ConverId:    converId,
			ChannelType: channelType,
			SubChannel:  subChannel,
		}, func() proto.Message {
			return &pbobjs.QryLatestMsgResp{}
		})
		if code == errs.IMErrorCode_SUCCESS && err == nil {
			if resp, ok := resp.(*pbobjs.QryLatestMsgResp); ok {
				msgSeq = resp.MsgSeqNo
				msgTime = resp.MsgTime
			}
		}
	*/
	return msgTime, msgSeq
}

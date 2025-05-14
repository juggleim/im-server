package services

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/caches"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/message/storages"
	"im-server/services/message/storages/models"
	"time"
)

var brdLocks *tools.SegmentatedLocks
var brdMsgCache *caches.LruCache

const (
	maxBrdMsg int = 10
)

func init() {
	brdLocks = tools.NewSegmentatedLocks(128)
	brdMsgCache = caches.NewLruCacheWithReadTimeout("brdmsg_cache", 10000, nil, 10*time.Minute)
}

func SyncBrdMsgs(ctx context.Context, appkey string, startTime int64, count int) []*pbobjs.DownMsg {
	lock := brdLocks.GetLocks(appkey)
	lock.Lock()
	defer lock.Unlock()

	var ringArr *tools.RingArray
	if obj, exist := brdMsgCache.Get(appkey); exist {
		ringArr = obj.(*tools.RingArray)
	} else {
		//load from db
		brdStorage := storages.NewBrdInboxMsgStorage()
		msgs, err := brdStorage.QryLatestMsg(appkey, maxBrdMsg)
		ringArr = tools.NewRingArray(maxBrdMsg)
		if err == nil {
			for _, msg := range msgs {
				var downMsg pbobjs.DownMsg
				err := tools.PbUnMarshal(msg.MsgBody, &downMsg)
				if err == nil {
					ringArr.Append(&downMsg)
				}
			}
		}
		brdMsgCache.Add(appkey, ringArr)
	}
	brdMsgs := []*pbobjs.DownMsg{}
	if ringArr != nil {
		ringArr.Foreach(func(item interface{}) bool {
			msg, ok := item.(*pbobjs.DownMsg)
			if ok {
				if msg.MsgTime > startTime {
					brdMsgs = append(brdMsgs, msg)
					if len(brdMsgs) > count {
						return false
					}
				}
			}
			return true
		})
	}
	return brdMsgs
}

func BrdAppendMsg(ctx context.Context, msg *pbobjs.DownMsg) errs.IMErrorCode {
	if msg != nil {
		appkey := bases.GetAppKeyFromCtx(ctx)
		lock := brdLocks.GetLocks(appkey)
		lock.Lock()
		defer lock.Unlock()
		if obj, exist := brdMsgCache.Get(appkey); exist {
			ringArr := obj.(*tools.RingArray)
			ringArr.Append(msg)
		}
	}
	return errs.IMErrorCode_SUCCESS
}

func SaveBroadcastMsg(ctx context.Context, msg *pbobjs.DownMsg) errs.IMErrorCode {
	data, _ := tools.PbMarshal(msg)
	brdStorage := storages.NewBrdInboxMsgStorage()
	brdStorage.SaveMsg(models.BrdInboxMsgMsg{
		SenderId:    msg.SenderId,
		SendTime:    msg.MsgTime,
		MsgId:       msg.MsgId,
		ChannelType: msg.ChannelType,
		MsgBody:     data,
		AppKey:      bases.GetAppKeyFromCtx(ctx),
	})
	return errs.IMErrorCode_SUCCESS
}

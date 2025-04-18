package commonservices

import (
	"context"
	"fmt"
	"im-server/commons/bases"
	"im-server/commons/caches"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/commonservices/msgdefines"
	"time"

	"google.golang.org/protobuf/proto"
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

func GetMsgConverCache(ctx context.Context, converId string, channelType pbobjs.ChannelType) *MsgConverItem {
	appkey := bases.GetAppKeyFromCtx(ctx)
	cacheKey := getMsgConverCacheKey(appkey, converId, channelType)
	lock := msgConverLocks.GetLocks(cacheKey)
	lock.Lock()
	defer lock.Unlock()
	if val, ok := msgConverCache.Get(cacheKey); ok {
		return val.(*MsgConverItem)
	} else {
		msgTime, msgSeq := GetLatestMsgTimeSeq(ctx, appkey, converId, channelType)
		item := &MsgConverItem{
			cacheKey:      cacheKey,
			LatestMsgTime: msgTime,
			LatestMsgSeq:  msgSeq,
		}
		msgConverCache.Add(cacheKey, item)
		return item
	}
}

func getMsgConverCacheKey(appkey, converId string, channelType pbobjs.ChannelType) string {
	return fmt.Sprintf("%s_%d_%s", appkey, channelType, converId)
}

func GetLatestMsgTimeSeq(ctx context.Context, appkey, converId string, channelType pbobjs.ChannelType) (int64, int64) {
	var msgSeq int64 = 0
	var msgTime int64 = 0
	//从会话查询最新的index
	code, resp, err := bases.SyncRpcCall(ctx, "qry_latest_hismsg", converId, &pbobjs.QryLatestMsgReq{
		ConverId:    converId,
		ChannelType: channelType,
	}, func() proto.Message {
		return &pbobjs.QryLatestMsgResp{}
	})
	if code == errs.IMErrorCode_SUCCESS && err == nil {
		if resp, ok := resp.(*pbobjs.QryLatestMsgResp); ok {
			msgSeq = resp.MsgSeqNo
			msgTime = resp.MsgTime
		}
	}

	return msgTime, msgSeq
}

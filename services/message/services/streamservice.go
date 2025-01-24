package services

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/caches"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/commonservices"
	"strings"
	"time"

	"github.com/bytedance/gopkg/collection/skipmap"
)

var streamMsgCache *caches.LruCache
var streamMsgLocks *tools.SegmentatedLocks

func init() {
	streamMsgCache = caches.NewLruCacheWithReadTimeout(10000, nil, 5*time.Minute)
	streamMsgLocks = tools.NewSegmentatedLocks(128)
}

type StreamMsg struct {
	Appkey   string
	SenderId string
	TargetId string
	MsgId    string
	MaxSeq   int
	MsgType  string

	streamMsgItems *skipmap.Int64Map
}

func AppendStreamMsgItem(ctx context.Context, req *pbobjs.StreamDownMsg) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	senderId := bases.GetRequesterIdFromCtx(ctx)
	targetId := bases.GetTargetIdFromCtx(ctx)
	key := strings.Join([]string{appkey, req.MsgId}, "_")
	l := streamMsgLocks.GetLocks(key)
	l.Lock()
	defer l.Unlock()
	isFinish := false
	if val, exist := streamMsgCache.Get(key); exist {
		sMsg := val.(*StreamMsg)
		for _, item := range req.MsgItems {
			sMsg.MaxSeq = sMsg.MaxSeq + 1
			if item.SubSeq <= 0 {
				item.SubSeq = int64(sMsg.MaxSeq)
			}
			sMsg.streamMsgItems.Store(item.SubSeq, item)
			if item.Event == pbobjs.StreamEvent_StreamComplete {
				isFinish = true
			}
		}
		if isFinish {
			streamMsgCache.Remove(key)
			//update
			updateStreamMsg(ctx, sMsg)
		}
	} else {
		sMsg := &StreamMsg{
			Appkey:         appkey,
			SenderId:       senderId,
			TargetId:       targetId,
			MsgId:          req.MsgId,
			MsgType:        req.MsgType,
			streamMsgItems: skipmap.NewInt64(),
		}
		for _, item := range req.MsgItems {
			sMsg.MaxSeq = sMsg.MaxSeq + 1
			if item.SubSeq <= 0 {
				item.SubSeq = int64(sMsg.MaxSeq)
			}
			sMsg.streamMsgItems.Store(item.SubSeq, item)
			if item.Event == pbobjs.StreamEvent_StreamComplete {
				isFinish = true
			}
		}
		if !isFinish {
			streamMsgCache.Add(key, sMsg)
		} else { //update
			updateStreamMsg(ctx, sMsg)
		}
	}
}

func updateStreamMsg(ctx context.Context, streamMsg *StreamMsg) {
	pbStreamMsg := &pbobjs.StreamDownMsg{
		TargetId:    streamMsg.TargetId,
		ChannelType: pbobjs.ChannelType_Private,
		MsgId:       streamMsg.MsgId,
		MsgItems:    []*pbobjs.StreamMsgItem{},
		MsgType:     streamMsg.MsgType,
	}
	streamMsg.streamMsgItems.Range(func(key int64, value interface{}) bool {
		if value != nil {
			if item, ok := value.(*pbobjs.StreamMsgItem); ok {
				pbStreamMsg.MsgItems = append(pbStreamMsg.MsgItems, item)
			}
		}
		return true
	})
	targetId := commonservices.GetConversationId(streamMsg.SenderId, streamMsg.TargetId, pbobjs.ChannelType_Private)
	bases.AsyncRpcCall(ctx, "upd_stream", targetId, pbStreamMsg)
}

func HandleStreamMsg(ctx context.Context, req *pbobjs.StreamDownMsg) errs.IMErrorCode {
	AppendStreamMsgItem(ctx, req)

	req.TargetId = bases.GetRequesterIdFromCtx(ctx)
	req.ChannelType = pbobjs.ChannelType_Private
	bases.AsyncRpcCall(ctx, "stream_msg", bases.GetTargetIdFromCtx(ctx), req)

	return errs.IMErrorCode_SUCCESS
}

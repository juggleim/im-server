package services

import (
	"context"
	"fmt"
	"im-server/commons/bases"
	"im-server/commons/caches"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/commonservices"
	"im-server/services/commonservices/logs"
	"im-server/services/commonservices/msgdefines"
	"im-server/services/historymsg/storages"
	"time"

	"github.com/bytedance/gopkg/collection/skipmap"
)

func SendStreamMsg(ctx context.Context, streamMsg *pbobjs.StreamMsg) errs.IMErrorCode {
	appkey := bases.GetAppKeyFromCtx(ctx)
	msgId := streamMsg.StreamMsgId

	userId := bases.GetRequesterIdFromCtx(ctx)
	targetId := streamMsg.TargetId
	streamMsgCacheItem := &StreamMsgCacheItem{
		Appkey:      appkey,
		SenderId:    userId,
		TargetId:    targetId,
		StreamMsgId: msgId,
		MaxSeq:      1,
	}

	if streamMsg.IsFinished && streamMsg.Seq == 1 {
		sendMsg(ctx, userId, targetId, msgdefines.NewStreamMsg(msgId, string(streamMsg.PartialContent), true, int(streamMsg.Seq)))
		return errs.IMErrorCode_SUCCESS
	}

	if streamMsg.Seq == 1 {
		succ := AddStreamMsg(ctx, streamMsgCacheItem)
		if !succ {
			return errs.IMErrorCode_DEFAULT
		}
		streamMsgCacheItem.Append(&PartialContent{
			Content:    string(streamMsg.PartialContent),
			Seq:        int(streamMsg.Seq),
			IsFinished: streamMsg.IsFinished,
		})
		sendMsg(ctx, userId, targetId, msgdefines.NewStreamMsg(msgId, string(streamMsg.PartialContent), false, int(streamMsg.Seq)))
	} else {
		cacheItem, exist := GetStreamMsg(ctx, msgId)
		if !exist {
			return errs.IMErrorCode_DEFAULT
		}
		streamMsgCacheItem = cacheItem
		streamMsgCacheItem.Append(&PartialContent{
			Content:    string(streamMsg.PartialContent),
			Seq:        int(streamMsg.Seq),
			IsFinished: streamMsg.IsFinished,
		})
		content := string(streamMsg.PartialContent)
		seq := streamMsg.Seq
		if streamMsg.IsFinished { //finished
			//update history msg
			finalContent, maxSeq := streamMsgCacheItem.FinalContent()
			content = finalContent
			seq = int64(maxSeq)
			storage := storages.NewPrivateHisMsgStorage()
			converId := commonservices.GetConversationId(userId, targetId, pbobjs.ChannelType_Private)
			dbMsg, err := storage.FindById(appkey, converId, "", msgId)
			if err == nil {
				newDownMsg := &pbobjs.DownMsg{}
				err = tools.PbUnMarshal(dbMsg.MsgBody, newDownMsg)
				if err == nil {
					newDownMsg.MsgContent = tools.ToJsonBs(msgdefines.StreamMsg{
						StreamId:   msgId,
						Content:    finalContent,
						IsFinished: true,
						Seq:        maxSeq,
					})
					newDownMsgBs, _ := tools.PbMarshal(newDownMsg)
					err = storage.UpdateMsgBody(appkey, converId, "", msgId, newDownMsg.MsgType, newDownMsgBs)
					if err != nil {
						logs.WithContext(ctx).Error(err.Error())
					}
				}
			}
			//remove from cache
			RemoveStreamMsg(ctx, msgId)
		}
		//send append msg
		sendMsg(ctx, userId, targetId, msgdefines.NewStreamAppendMsg(msgId, content, streamMsg.IsFinished, int(seq)))
	}
	return errs.IMErrorCode_SUCCESS
}

func sendMsg(ctx context.Context, senderId, targetId string, msg msgdefines.BaseStreamMsg) {
	msgType := msg.GetMsgType()
	msgId := msg.GetStreamId()

	var flag int32 = 0
	opts := []bases.BaseActorOption{}
	if msgType == msgdefines.InnerMsgType_StreamMsg {
		flag = msgdefines.SetCountMsg(flag)
		flag = msgdefines.SetStoreMsg(flag)
		opts = append(opts, &bases.WithMsgIdOption{
			MsgId: msgId,
		})
	}
	commonservices.AsyncPrivateMsgOverUpstream(ctx, senderId, targetId, &pbobjs.UpMsg{
		MsgType:    msgType,
		MsgContent: tools.ToJsonBs(msg),
		Flags:      flag,
	}, opts...)
}

var streamMsgCache *caches.LruCache
var streamMsgLocks *tools.SegmentatedLocks

func init() {
	streamMsgCache = caches.NewLruCacheWithReadTimeout("streammsg_cache", 10000, func(key, value interface{}) {
		if streamMsgCacheItem, ok := value.(*StreamMsgCacheItem); ok && streamMsgCacheItem != nil && !streamMsgCacheItem.isFinished {
			appkey := streamMsgCacheItem.Appkey
			userId := streamMsgCacheItem.SenderId
			targetId := streamMsgCacheItem.TargetId
			msgId := streamMsgCacheItem.StreamMsgId
			ctx := context.Background()
			ctx = context.WithValue(ctx, bases.CtxKey_AppKey, appkey)
			//update history msg
			finalContent, maxSeq := streamMsgCacheItem.FinalContent()
			maxSeq = maxSeq + 1
			storage := storages.NewPrivateHisMsgStorage()
			converId := commonservices.GetConversationId(userId, targetId, pbobjs.ChannelType_Private)
			dbMsg, err := storage.FindById(appkey, converId, "", msgId)
			if err == nil {
				newDownMsg := &pbobjs.DownMsg{}
				err = tools.PbUnMarshal(dbMsg.MsgBody, newDownMsg)
				if err == nil {
					newDownMsg.MsgContent = tools.ToJsonBs(msgdefines.StreamMsg{
						StreamId:   msgId,
						Content:    finalContent,
						IsFinished: true,
						Seq:        maxSeq,
					})
					newDownMsgBs, _ := tools.PbMarshal(newDownMsg)
					err = storage.UpdateMsgBody(appkey, converId, "", msgId, newDownMsg.MsgType, newDownMsgBs)
					if err != nil {
						logs.WithContext(ctx).Error(err.Error())
					}
				}
			}
			//send complete msg
			sendMsg(ctx, userId, targetId, msgdefines.NewStreamAppendMsg(msgId, finalContent, true, maxSeq))
		}
	}, 10*time.Minute)
	streamMsgLocks = tools.NewSegmentatedLocks(128)
}

type StreamMsgCacheItem struct {
	Appkey      string
	SenderId    string
	TargetId    string
	StreamMsgId string
	MaxSeq      int

	items *skipmap.Int64Map

	isFinished bool
}

type PartialContent struct {
	Content    string
	Seq        int
	IsFinished bool
}

func (item *StreamMsgCacheItem) Append(partial *PartialContent) {
	key := getStreamMsgKey(item.Appkey, item.StreamMsgId)
	l := streamMsgLocks.GetLocks(key)
	l.Lock()
	defer l.Unlock()
	if item.MaxSeq < partial.Seq {
		item.MaxSeq = partial.Seq
	}
	item.items.Store(int64(partial.Seq), partial)
}

func (item *StreamMsgCacheItem) FinalContent() (string, int) {
	key := getStreamMsgKey(item.Appkey, item.StreamMsgId)
	l := streamMsgLocks.GetLocks(key)
	l.Lock()
	defer l.Unlock()
	str := ""
	item.items.Range(func(key int64, value interface{}) bool {
		if partial, ok := value.(*PartialContent); ok {
			str = str + partial.Content
		}
		return true
	})
	return str, item.MaxSeq
}

func getStreamMsgKey(appkey, msgId string) string {
	return fmt.Sprintf("%s_%s", appkey, msgId)
}

func AddStreamMsg(ctx context.Context, item *StreamMsgCacheItem) bool {
	appkey := bases.GetAppKeyFromCtx(ctx)
	key := getStreamMsgKey(appkey, item.StreamMsgId)
	if item.items == nil {
		item.items = skipmap.NewInt64()
	}
	succ := streamMsgCache.AddIfAbsendNoGetOldVal(key, item)
	return succ
}

func GetStreamMsg(ctx context.Context, msgId string) (*StreamMsgCacheItem, bool) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	key := getStreamMsgKey(appkey, msgId)
	if !streamMsgCache.Contains(key) {
		return nil, false
	} else {
		if val, exist := streamMsgCache.Get(key); exist {
			return val.(*StreamMsgCacheItem), true
		} else {
			return nil, false
		}
	}
}

func RemoveStreamMsg(ctx context.Context, msgId string) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	key := getStreamMsgKey(appkey, msgId)
	if val, exist := streamMsgCache.Get(key); exist {
		item := val.(*StreamMsgCacheItem)
		item.isFinished = true
		streamMsgCache.Remove(key)
	}
}

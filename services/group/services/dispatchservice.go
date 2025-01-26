package services

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/gmicro/utils"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/commonservices"
	"im-server/services/commonservices/logs"
	"time"
)

var dispatchQueues []*Dispatcher

type GrpMsgDispatchItem struct {
	Ctx       context.Context
	GroupId   string
	MemberIds []string
	Msg       *pbobjs.DownMsg
}

type Dispatcher struct {
	Ratio       int
	msgQueue    chan *GrpMsgDispatchItem
	maxQueueLen int
	isContinued bool
	// limiter       *rate.Limiter
	latestUpdTime int64
}

func NewDispatcher(ratio, maxLen int) *Dispatcher {
	return &Dispatcher{
		Ratio:         ratio,
		maxQueueLen:   maxLen,
		latestUpdTime: time.Now().UnixMilli(),
	}
}

func (dis *Dispatcher) start() {
	utils.SafeGo(func() {
		for {
			if dis.isContinued {
				item := <-dis.msgQueue

				now := time.Now().UnixMilli()
				if now-item.Msg.MsgTime > 60*1000 {
					logs.WithContext(item.Ctx).Warnf("[discard_group_msg] group_id:%s\tmsg_id:%s\tmsg_time:%d", item.GroupId, item.Msg.MsgId, item.Msg.MsgTime)
					continue
				}

				memberCount := len(item.MemberIds)
				groupCastMsg(item.Ctx, item.GroupId, item.MemberIds, item.Msg)
				interval := dis.Ratio * 100
				time.Sleep(time.Duration(interval) * time.Millisecond)
				logs.WithContext(item.Ctx).Warnf("[group_dispatch_delay] group_id:%s\tmember_count:%d\tsleep:%d", item.GroupId, memberCount, interval)
			} else {
				break
			}
		}
		close(dis.msgQueue)
	})
}

func getGrpThresholdFromCtx(ctx context.Context) int {
	grpThreshold := 10000
	appinfo, exist := commonservices.GetAppInfo(bases.GetAppKeyFromCtx(ctx))
	if exist && appinfo != nil && appinfo.GrpMsgThreshold > 0 {
		grpThreshold = appinfo.GrpMsgThreshold
	}
	return grpThreshold
}

func (dis *Dispatcher) Stop() {
	dis.isContinued = false
}

func (dis *Dispatcher) Put(item *GrpMsgDispatchItem) {
	if !dis.isContinued {
		dis.isContinued = true
		dis.msgQueue = make(chan *GrpMsgDispatchItem, dis.maxQueueLen)
		dis.start()
	}
	dis.msgQueue <- item
}

func init() {
	dispatchQueues = make([]*Dispatcher, 10)
	for i := 0; i < 10; i++ {
		dispatchQueues[i] = NewDispatcher(i+1, 10000)
	}
}

func Dispatch2Message(ctx context.Context, groupId string, memberIds []string, msg *pbobjs.DownMsg) {
	memberCount := len(memberIds)

	ratio := memberCount / getGrpThresholdFromCtx(ctx)
	if ratio <= 0 {
		groupCastMsg(ctx, groupId, memberIds, msg)
	} else {
		item := &GrpMsgDispatchItem{
			Ctx:       ctx,
			GroupId:   groupId,
			MemberIds: memberIds,
			Msg:       msg,
		}
		if ratio >= 10 {
			dispatchQueues[9].Put(item)
		} else {
			dispatchQueues[ratio-1].Put(item)
		}
	}
}

func groupCastMsg(ctx context.Context, groupId string, memberIds []string, msg *pbobjs.DownMsg) {
	data, _ := tools.PbMarshal(msg)
	groups := bases.GroupTargets("msg_dispatch", memberIds)
	for _, ids := range groups {
		bases.UnicastRouteWithNoSender(&pbobjs.RpcMessageWraper{
			RpcMsgType:   pbobjs.RpcMsgType_UserPub,
			AppKey:       bases.GetAppKeyFromCtx(ctx),
			Session:      bases.GetSessionFromCtx(ctx),
			Method:       "msg_dispatch",
			RequesterId:  bases.GetRequesterIdFromCtx(ctx),
			ReqIndex:     bases.GetSeqIndexFromCtx(ctx),
			Qos:          bases.GetQosFromCtx(ctx),
			AppDataBytes: data,
			TargetId:     ids[0],
			GroupId:      groupId,
			TargetIds:    ids,
		})
	}
}

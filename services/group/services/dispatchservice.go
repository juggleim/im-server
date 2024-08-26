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

	"golang.org/x/time/rate"
)

var msgThreshold = 2000
var dispatchQueues []*Dispatcher

type GrpMsgDispatchItem struct {
	Ctx       context.Context
	GroupId   string
	MemberIds []string
	Msg       *pbobjs.DownMsg
}

type Dispatcher struct {
	msgQueue      chan *GrpMsgDispatchItem
	maxQueueLen   int
	isContinued   bool
	maxDisCount   int
	limiter       *rate.Limiter
	latestUpdTime int64
}

func NewDispatcher(maxLen int) *Dispatcher {
	return &Dispatcher{
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
					logs.WithContext(item.Ctx).Warnf("discard group msg:%s, msg_time:%d", item.Msg.MsgId, item.Msg.MsgTime)
					continue
				}

				memberCount := len(item.MemberIds)
				maxDisCount := dis.getMaxDisCount()
				ratio := memberCount / getGrpThresholdFromCtx(item.Ctx)
				if ratio <= 0 { //dispatch directly
					groupCastMsg(item.Ctx, item.GroupId, item.MemberIds, item.Msg)
				} else if memberCount > maxDisCount {
					groupCastMsg(item.Ctx, item.GroupId, item.MemberIds, item.Msg)
					sleepTime := memberCount / maxDisCount
					time.Sleep(time.Duration(sleepTime) * time.Second)
				} else {
					if dis.limiter == nil {
						dis.limiter = rate.NewLimiter(rate.Limit(maxDisCount), maxDisCount)
					}
					allow := dis.limiter.AllowN(time.Now(), memberCount)
					for !allow {
						time.Sleep(5 * time.Millisecond)
						allow = dis.limiter.AllowN(time.Now(), memberCount)
					}
					groupCastMsg(item.Ctx, item.GroupId, item.MemberIds, item.Msg)
				}
			} else {
				break
			}
		}
		close(dis.msgQueue)
	})
}

func getGrpThresholdFromCtx(ctx context.Context) int {
	grpThreshold := 100
	appinfo, exist := commonservices.GetAppInfo(bases.GetAppKeyFromCtx(ctx))
	if exist && appinfo != nil && appinfo.GrpMsgThreshold > 0 {
		grpThreshold = appinfo.GrpMsgThreshold
	}
	return grpThreshold
}

func (dis *Dispatcher) getMaxDisCount() int {
	now := time.Now().UnixMilli()
	if now-5000 > dis.latestUpdTime {
		msgNodeCount := bases.GetCluster().GetTargetNodeCount("g_msg_dispatch")
		if msgNodeCount <= 0 {
			msgNodeCount = 1
		}
		maxDisCount := msgNodeCount * msgThreshold
		if maxDisCount != dis.maxDisCount {
			dis.maxDisCount = maxDisCount
			dis.latestUpdTime = now
			dis.limiter = rate.NewLimiter(rate.Limit(dis.maxDisCount), dis.maxDisCount)
		}
	}
	return dis.maxDisCount
}

func (dis *Dispatcher) Stop() {
	dis.isContinued = false
}

func (dis *Dispatcher) Put(item *GrpMsgDispatchItem) {
	if !dis.isContinued {
		dis.isContinued = true
		dis.msgQueue = make(chan *GrpMsgDispatchItem, dis.maxQueueLen)
	}
	dis.msgQueue <- item
	dis.start()
}

func init() {
	dispatchQueues = make([]*Dispatcher, 9)
	for i := 0; i < 9; i++ {
		dispatchQueues[i] = NewDispatcher(10000)
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

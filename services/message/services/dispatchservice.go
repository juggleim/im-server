package services

import (
	"context"
	"fmt"
	"im-server/commons/bases"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/botmsg/botclient"
	"im-server/services/commonservices"
	"im-server/services/commonservices/logs"
	"im-server/services/commonservices/msgdefines"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"time"
)

var MsgSinglePools *tools.SinglePools

const (
	msgDispatchFanoutWorkers         = 64
	msgDispatchFanoutCapacity        = 512
	msgDispatchLocalBatchSize        = 500
	msgDispatchLocalBatchConcurrency = 4
)

var msgDispatchFanoutScheduler = newKeyedFanoutScheduler(msgDispatchFanoutWorkers, msgDispatchFanoutCapacity)

func init() {
	MsgSinglePools = tools.NewSinglePools(8192, 100, true)
}

func DispatchMsg(ctx context.Context, downMsg *pbobjs.DownMsg) {
	started := time.Now()
	targetCount := len(bases.GetTargetIdsFromCtx(ctx))
	if downMsg.ChannelType == pbobjs.ChannelType_Private || downMsg.ChannelType == pbobjs.ChannelType_System {
		targetCount = 1
	}
	dispatchTargets.observe(targetCount)
	if downMsg.ChannelType == pbobjs.ChannelType_Group {
		fanoutKey := msgDispatchFanoutKey(ctx, downMsg)
		fanoutSubmitDuration, blocked := msgDispatchFanoutScheduler.Submit(fanoutKey, func(queueWait time.Duration) {
			dispatchTimings.fanoutQueueWait.observe(queueWait)
			dispatchFanout(ctx, downMsg, targetCount)
		})
		dispatchTimings.fanoutSubmit.observe(fanoutSubmitDuration)
		if blocked {
			dispatchTimings.fanoutFullOnSubmit.Add(1)
		}
		dispatchTimings.actorSubmit.observe(time.Since(started))
		return
	}
	dispatchFanout(ctx, downMsg, targetCount)
	dispatchTimings.actorSubmit.observe(time.Since(started))
}

func msgDispatchFanoutKey(ctx context.Context, downMsg *pbobjs.DownMsg) string {
	targetId := bases.GetGroupIdFromCtx(ctx)
	if targetId == "" {
		targetId = downMsg.TargetId
	}
	return strings.Join([]string{bases.GetAppKeyFromCtx(ctx), strconv.Itoa(int(downMsg.ChannelType)), targetId}, "_")
}

func dispatchFanout(ctx context.Context, downMsg *pbobjs.DownMsg, targetCount int) {
	started := time.Now()
	preheatDuration := time.Duration(0)
	defer func() {
		totalDuration := time.Since(started)
		fanoutDuration := totalDuration - preheatDuration
		dispatchTimings.fanoutTask.observe(totalDuration)
		dispatchTimings.fanout.observe(fanoutDuration)
	}()
	appkey := bases.GetAppKeyFromCtx(ctx)
	if downMsg.ChannelType == pbobjs.ChannelType_Private || downMsg.ChannelType == pbobjs.ChannelType_System {
		receiverId := bases.GetTargetIdFromCtx(ctx)
		submitDispatch(ctx, appkey, receiverId, func() {
			doDispatch(ctx, receiverId, downMsg, false)
		})
	} else if downMsg.ChannelType == pbobjs.ChannelType_Group {
		memberIds := bases.GetTargetIdsFromCtx(ctx)
		ctx = bases.SetTargetIds2Ctx(ctx, []string{})
		threadhold := 1000
		forceCloseOfflineMsg := false
		appinfo, exist := commonservices.GetAppInfo(appkey)
		if exist && appinfo != nil {
			threadhold = appinfo.BigGrpThreshold
			forceCloseOfflineMsg = appinfo.ForceCloseOfflineMsg
		}
		if downMsg.MemberCount > int32(threadhold) {
			preheatStarted := time.Now()
			preheat(ctx, appkey, memberIds)
			preheatDuration = time.Since(preheatStarted)
			dispatchTimings.preheat.observe(preheatDuration)
			processDispatchMemberBatches(ctx, memberIds, func(receiverId string) {
				newDownMsg := commonservices.CopyDownMsg(downMsg, receiverId)
				closeOffline := forceCloseOfflineMsg
				if !closeOffline {
					userStatus := GetUserStatus(appkey, receiverId)
					closeOffline = (!userStatus.IsOnline())
				}
				submitDispatch(ctx, appkey, receiverId, func() {
					doDispatch(ctx, receiverId, newDownMsg, closeOffline)
				})
			})
		} else {
			processDispatchMemberBatches(ctx, memberIds, func(receiverId string) {
				newDownMsg := commonservices.CopyDownMsg(downMsg, receiverId)
				submitDispatch(ctx, appkey, receiverId, func() {
					doDispatch(ctx, receiverId, newDownMsg, false)
				})
			})
		}
	} else if downMsg.ChannelType == pbobjs.ChannelType_SubStatus {
		memberIds := bases.GetTargetIdsFromCtx(ctx)
		ctx = bases.SetTargetIds2Ctx(ctx, []string{})
		for _, receiverId := range memberIds {
			receId := receiverId
			if UserStatusCacheContains(appkey, receId) {
				userStatus := GetUserStatus(appkey, receId)
				if userStatus.IsOnline() {
					newDownMsg := commonservices.CopyDownMsg(downMsg, receiverId)
					MsgSinglePools.GetPool(strings.Join([]string{appkey, receId}, "_")).Submit(func() {
						dispatchStatusSubMsgs(ctx, receId, newDownMsg)
					})
				}
			}
		}
	}
}

func processDispatchMemberBatches(ctx context.Context, memberIds []string, process func(receiverId string)) {
	if len(memberIds) <= msgDispatchLocalBatchSize {
		for _, receiverId := range memberIds {
			process(receiverId)
		}
		return
	}
	semaphore := make(chan struct{}, msgDispatchLocalBatchConcurrency)
	var wg sync.WaitGroup
	for start := 0; start < len(memberIds); start += msgDispatchLocalBatchSize {
		end := start + msgDispatchLocalBatchSize
		if end > len(memberIds) {
			end = len(memberIds)
		}
		batch := memberIds[start:end]
		semaphore <- struct{}{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() {
				<-semaphore
				if err := recover(); err != nil {
					logs.WithContext(ctx).Errorf("msg dispatch local batch panic:%v\n%s", err, debug.Stack())
				}
			}()
			for _, receiverId := range batch {
				process(receiverId)
			}
		}()
	}
	wg.Wait()
}

func submitDispatch(ctx context.Context, appkey, receiverId string, task func()) {
	poolKey := strings.Join([]string{appkey, receiverId}, "_")
	pool := MsgSinglePools.GetPool(poolKey)
	queueLenOnSubmit := pool.QueueLen()
	if queueLenOnSubmit == pool.QueueCapacity() {
		dispatchTimings.queueFullOnSubmit.Add(1)
	}
	enqueued := time.Now()
	pool.Submit(func() {
		queueWait := time.Since(enqueued)
		dispatchTimings.queueWait.observe(queueWait)
		started := time.Now()
		defer func() { dispatchTimings.task.observe(time.Since(started)) }()
		task()
	})
	dispatchTimings.submit.observe(time.Since(enqueued))
	observeDispatchQueueLen(pool.QueueLen())
}

func preheat(ctx context.Context, appkey string, memberIds []string) {
	//preheat user status
	noStatusCacheUids := []string{}
	for _, receiverId := range memberIds {
		if !UserStatusCacheContains(appkey, receiverId) {
			noStatusCacheUids = append(noStatusCacheUids, receiverId)
		}
	}
	if len(noStatusCacheUids) > 0 {
		BatchInitUserStatus(ctx, appkey, noStatusCacheUids)
	}
}

func doDispatch(ctx context.Context, receiverId string, msg *pbobjs.DownMsg, closeOffline bool) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	//TODO save imediately when user online, other wise, use async queue.
	//TODO batch insert
	sendTime := RegenateSendTime(appkey, receiverId, msg.MsgTime)
	msg.MsgTime = sendTime
	//handle conversation check, such as undisturb, unread index
	phaseStarted := time.Now()
	HandleDownMsgByConver(ctx, receiverId, msg.TargetId, msg.SubChannel, msg.ChannelType, msg)
	dispatchTimings.conversation.observe(time.Since(phaseStarted))
	phaseStarted = time.Now()
	targetUserInfo := commonservices.GetTargetUserInfo(ctx, receiverId)
	dispatchTimings.userInfo.observe(time.Since(phaseStarted))
	if targetUserInfo.UserType == pbobjs.UserType_Bot {
		if msg.ChannelType != pbobjs.ChannelType_Private && msg.ChannelType != pbobjs.ChannelType_Group {
			return
		}
		if msgdefines.IsStateMsg(msg.Flags) || msgdefines.IsCmdMsg(msg.Flags) {
			return
		}
		if msg.ChannelType == pbobjs.ChannelType_Group {
			onlyMentioned := true
			if targetUserInfo.BotSettings != nil {
				onlyMentioned = targetUserInfo.BotSettings.OnlyMentioned
			}
			if onlyMentioned && !commonservices.IsDirectMentionedMe(targetUserInfo.UserId, msg) {
				return
			}
		}
		if targetUserInfo.HasBotConf {
			botclient.SendMsg2Bot(ctx, receiverId, msg)
		}
		if msgdefines.IsStateMsg(msg.Flags) {
			MsgDirect(ctx, receiverId, msg)
		} else {
			commonservices.SaveConversation(ctx, receiverId, msg)
			SaveMsg2Inbox(appkey, receiverId, msg)
			//send to client
			MsgOrNtf(ctx, receiverId, msg, true)
		}
	} else {
		if !msgdefines.IsStateMsg(msg.Flags) {
			if !closeOffline {
				phaseStarted = time.Now()
				commonservices.SaveConversation(ctx, receiverId, msg)
				dispatchTimings.downstream.observe(time.Since(phaseStarted))
				phaseStarted = time.Now()
				SaveMsg2Inbox(appkey, receiverId, msg)
				dispatchTimings.inbox.observe(time.Since(phaseStarted))
				// send to client
				phaseStarted = time.Now()
				MsgOrNtfWithPush(ctx, receiverId, msg)
				dispatchTimings.downstream.observe(time.Since(phaseStarted))
			} else {
				phaseStarted = time.Now()
				target := fmt.Sprintf("%s_%s_%d", appkey, msg.TargetId, msg.ChannelType)
				batchExecutorPool.GetBatchExecutor(target).Append(fmt.Sprintf("%s_%s", target, receiverId), &BatchConverItem{
					Appkey: appkey,
					UserId: receiverId,
					Msg:    msg,
				})
			}
		} else if !closeOffline {
			MsgDirect(ctx, receiverId, msg)
		}
	}
}

func copyDownMsg(msg *pbobjs.DownMsg, receiverId string) *pbobjs.DownMsg {
	referMsg := msg.ReferMsg
	if referMsg != nil {
		if referMsg.SenderId == receiverId {
			referMsg.IsSend = true
		} else {
			referMsg.IsSend = false
		}
	}
	return &pbobjs.DownMsg{
		TargetId:          msg.TargetId,
		ChannelType:       msg.ChannelType,
		MsgType:           msg.MsgType,
		SenderId:          msg.SenderId,
		MsgId:             msg.MsgId,
		MsgSeqNo:          msg.MsgSeqNo,
		MsgContent:        msg.MsgContent,
		MsgTime:           msg.MsgTime,
		Flags:             msg.Flags,
		IsSend:            msg.IsSend,
		Platform:          msg.Platform,
		ClientUid:         msg.ClientUid,
		PushData:          msg.PushData,
		MentionInfo:       msg.MentionInfo,
		IsRead:            msg.IsRead,
		ReferMsg:          referMsg,
		TargetUserInfo:    msg.TargetUserInfo,
		SenderInfo:        msg.SenderInfo,
		GroupInfo:         msg.GroupInfo,
		MergedMsgs:        msg.MergedMsgs,
		UndisturbType:     msg.UndisturbType,
		MemberCount:       msg.MemberCount,
		ReadCount:         msg.ReadCount,
		UnreadIndex:       msg.UnreadIndex,
		SearchText:        msg.SearchText,
		GrpMemberInfo:     msg.GrpMemberInfo,
		DestroyTime:       msg.DestroyTime,
		LifeTimeAfterRead: msg.LifeTimeAfterRead,
		SubChannel:        msg.SubChannel,
		ConverTags:        msg.ConverTags,
	}
}

func UserPush(ctx context.Context, msg *pbobjs.DownMsg) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	targetId := bases.GetTargetIdFromCtx(ctx)
	if !UserStatusCacheContains(appkey, targetId) {
		BatchInitUserStatus(ctx, appkey, []string{targetId})
	}
	userStatus := GetUserStatus(appkey, targetId)
	if userStatus != nil {
		if msg.PushData != nil && msg.PushData.IsVoip {
			SendPush(ctx, bases.GetRequesterIdFromCtx(ctx), targetId, msg)
		} else {
			if !userStatus.IsOnline() || userStatus.OpenPushSwitch() {
				SendPush(ctx, bases.GetRequesterIdFromCtx(ctx), targetId, msg)
			}
		}
	}
}

func dispatchStatusSubMsgs(ctx context.Context, receiverId string, msg *pbobjs.DownMsg) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	sendTime := RegenateSendTime(appkey, receiverId, msg.MsgTime)
	msg.MsgTime = sendTime
	MsgOrNtf(ctx, receiverId, msg, true)
}

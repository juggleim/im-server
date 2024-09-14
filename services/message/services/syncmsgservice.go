package services

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/commonservices"
	"im-server/services/message/storages"
	"im-server/services/message/storages/models"
	"sort"
	"time"
)

var msgSyncBatchCount int = 200

func SyncMessages(ctx context.Context, syncMsg *pbobjs.SyncMsgReq) (errs.IMErrorCode, *pbobjs.DownMsgSet) {
	userId := bases.GetTargetIdFromCtx(ctx)
	appKey := bases.GetAppKeyFromCtx(ctx)

	appinfo, exist := commonservices.GetAppInfo(appKey)
	restrictTime := time.Now().Add(-time.Minute * 1440).UnixMilli()
	if exist && appinfo != nil {
		restrictTime = time.Now().Add(-time.Minute * time.Duration(appinfo.OfflineMsgSaveTime)).UnixMilli()
	}
	syncTime := syncMsg.SyncTime
	sendboxSyncTime := syncMsg.SendBoxSyncTime
	if syncTime < restrictTime {
		syncTime = restrictTime
	}
	if sendboxSyncTime < restrictTime {
		sendboxSyncTime = restrictTime
	}

	//记录用户在离线状态
	RecordUserOnlineStatus(appKey, userId, true, int(bases.GetTerminalNumFromCtx(ctx)))
	//关闭直发
	userStatus := GetUserStatus(appKey, userId)
	userStatus.CheckNtfWithSwitch()
	//关闭强制推送
	userStatus.SetPushSwitch(0)
	//clear badge
	userStatus.SetBadge(0)

	ret := &pbobjs.DownMsgSet{
		Msgs: []*pbobjs.DownMsg{},
	}
	//拉取收件箱
	if userStatus.LatestMsgTime == nil || *userStatus.LatestMsgTime > syncTime {
		inboxMsgs := SyncInboxMessages(appKey, userId, syncTime, syncMsg.SyncTime, msgSyncBatchCount)
		for _, msg := range inboxMsgs {
			downMsg := &pbobjs.DownMsg{}
			err := tools.PbUnMarshal(msg.MsgBody, downMsg)
			if err == nil {
				ret.Msgs = append(ret.Msgs, downMsg)
			}
		}
	}

	//拉取发件箱
	if syncMsg.ContainsSendBox {
		sendboxMsgs := SyncSendboxMessages(appKey, userId, sendboxSyncTime, syncMsg.SendBoxSyncTime, msgSyncBatchCount)
		for _, msg := range sendboxMsgs {
			downMsg := &pbobjs.DownMsg{}
			err := tools.PbUnMarshal(msg.MsgBody, downMsg)
			if err == nil {
				ret.Msgs = append(ret.Msgs, downMsg)
			}
		}
	}

	//拉取广播消息
	brdMsgs := SyncBrdMsgs(ctx, appKey, syncMsg.SyncTime, msgSyncBatchCount)
	ret.Msgs = append(ret.Msgs, brdMsgs...)

	//re-sort
	sort.Slice(ret.Msgs, func(i, j int) bool {
		return ret.Msgs[i].MsgTime < ret.Msgs[j].MsgTime
	})
	if len(ret.Msgs) >= msgSyncBatchCount {
		ret.Msgs = ret.Msgs[:msgSyncBatchCount]
	} else {
		ret.IsFinished = true
		//变更通知拉取状态
		GetUserStatus(appKey, userId).SetNtfStatus(false)
		if userStatus.LatestMsgTime == nil {
			var maxMsgTime int64 = 0
			for _, msg := range ret.Msgs {
				if !msg.IsSend {
					if msg.MsgTime > maxMsgTime {
						maxMsgTime = msg.MsgTime
					}
				}
			}
			if maxMsgTime > 0 {
				userStatus.SetLatestMsgTime(maxMsgTime)
			}
		}
	}
	//statistic
	if len(ret.Msgs) > 0 {
		for _, msg := range ret.Msgs {
			commonservices.ReportDownMsg(appKey, msg.ChannelType, 1)
		}
	}

	return errs.IMErrorCode_SUCCESS, ret
}

func SyncInboxMessages(appkey, userid string, startTime, cmdStartTime int64, count int) []*models.Msg {
	retMsgs := []*models.Msg{}
	//cmd msg
	cmdStorage := storages.NewCmdInboxMsgStorage()
	msgs, err := cmdStorage.QryMsgsBaseTime(appkey, userid, cmdStartTime, count)
	if err == nil {
		retMsgs = append(retMsgs, msgs...)
	}

	//general msg
	storage := storages.NewInboxMsgStorage()
	msgs, err = storage.QryMsgsBaseTime(appkey, userid, startTime, count)
	if err == nil {
		retMsgs = append(retMsgs, msgs...)
	}
	return retMsgs
}

func SyncSendboxMessages(appkey, userid string, startTime, cmdStartTime int64, count int) []*models.Msg {
	retMsgs := []*models.Msg{}
	//cmd msg
	cmdStorage := storages.NewCmdSendboxMsgStorage()
	msgs, err := cmdStorage.QryMsgsBaseTime(appkey, userid, cmdStartTime, count)
	if err == nil {
		retMsgs = append(retMsgs, msgs...)
	}

	//general msg
	storage := storages.NewSendboxMsgStorage()
	msgs, err = storage.QryMsgsBaseTime(appkey, userid, startTime, count)
	if err == nil {
		retMsgs = append(retMsgs, msgs...)
	}
	return retMsgs
}

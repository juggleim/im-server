package services

import (
	"context"
	"fmt"
	"im-server/commons/bases"
	"im-server/commons/configures"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/commonservices"
	"im-server/services/commonservices/logs"
	"im-server/services/commonservices/msgdefines"
	"im-server/services/message/storages"
	"im-server/services/message/storages/models"
	"sync"
	"time"
)

const (
	msgPurgeWorkerCount   = 8
	msgPurgeQueueCapacity = 10000
	msgPurgeInterval      = time.Minute
)

var (
	msgPurgeOnce      sync.Once
	msgPurgeScheduler *purgeScheduler
)

// TODO save immediately when user online, other wise, use async queue.
func SaveMsg2Inbox(appkey, receiverId string, msg *pbobjs.DownMsg) error {
	if appkey == "" {
		err := fmt.Errorf("refuse to store inbox message with empty appkey, receiver_id:%s msg_id:%s", receiverId, msg.MsgId)
		logs.NewLogEntity().Error(err.Error())
		return err
	}
	var err error
	msgBs, _ := tools.PbMarshal(msg)
	message := models.Msg{
		UserId:      receiverId,
		SendTime:    msg.MsgTime,
		MsgId:       msg.MsgId,
		ChannelType: msg.ChannelType,
		MsgBody:     msgBs,
		AppKey:      appkey,
		TargetId:    msg.TargetId,
		MsgType:     msg.MsgType,
	}
	if msgdefines.IsCmdMsg(msg.Flags) {
		msgStorage := storages.NewCmdInboxMsgStorage()
		err = msgStorage.SaveMsg(message)
		cmdMsgExpired := getMsgExpired(appkey, true)
		purgeMsgs(appkey+":cmd_inbox", msg.MsgTime, msg.MsgTime-cmdMsgExpired, func(cutoff int64) error {
			purgeErr := msgStorage.DelMsgsBaseTime(appkey, cutoff)
			fmt.Println("clear offline cmd inbox msgs:", purgeErr, msg.MsgTime, cmdMsgExpired)
			return purgeErr
		})
	} else {
		msgStorage := storages.NewInboxMsgStorage()
		err = msgStorage.SaveMsg(message)
		msgExpired := getMsgExpired(appkey, false)
		purgeMsgs(appkey+":inbox", msg.MsgTime, msg.MsgTime-msgExpired, func(cutoff int64) error {
			purgeErr := msgStorage.DelMsgsBaseTime(appkey, cutoff)
			fmt.Println("clear offline inbox msgs:", purgeErr, msg.MsgTime, msgExpired)
			return purgeErr
		})
	}
	if err != nil {
		logs.NewLogEntity().Errorf("failed to store inbox. err:%v", err)
	}
	return err
}

func SaveMsg2Sendbox(ctx context.Context, appkey, senderId string, msg *pbobjs.DownMsg) error {
	if appkey == "" {
		err := fmt.Errorf("refuse to store sendbox message with empty appkey, sender_id:%s msg_id:%s", senderId, msg.MsgId)
		logs.NewLogEntity().Error(err.Error())
		return err
	}
	//save to sendbox
	msgBs, _ := tools.PbMarshal(msg)
	var err error
	message := models.Msg{
		UserId:      senderId,
		SendTime:    msg.MsgTime,
		MsgId:       msg.MsgId,
		ChannelType: msg.ChannelType,
		MsgBody:     msgBs,
		AppKey:      appkey,
		TargetId:    msg.TargetId,
		MsgType:     msg.MsgType,
	}
	if msgdefines.IsCmdMsg(msg.Flags) {
		storage := storages.NewCmdSendboxMsgStorage()
		if msg.MsgType == msgdefines.CmdMsgType_ClearUnread {
			rpcExts := bases.GetExtsFromCtx(ctx)
			if uniqTag, ok := rpcExts[commonservices.RpcExtKey_UniqTag]; ok && uniqTag != "" {
				message.UniqTag = uniqTag
				err = storage.UpsertMsg(message)
			} else {
				err = storage.SaveMsg(message)
			}
		} else {
			err = storage.SaveMsg(message)
		}
		cmdMsgExpired := getMsgExpired(appkey, true)
		purgeMsgs(appkey+":cmd_sendbox", msg.MsgTime, msg.MsgTime-cmdMsgExpired, func(cutoff int64) error {
			purgeErr := storage.DelMsgsBaseTime(appkey, cutoff)
			fmt.Println("clear offline cmd sendbox msgs:", purgeErr, msg.MsgTime, cmdMsgExpired)
			return purgeErr
		})
	} else {
		storage := storages.NewSendboxMsgStorage()
		err = storage.SaveMsg(message)
		msgExpired := getMsgExpired(appkey, false)
		purgeMsgs(appkey+":sendbox", msg.MsgTime, msg.MsgTime-msgExpired, func(cutoff int64) error {
			purgeErr := storage.DelMsgsBaseTime(appkey, cutoff)
			fmt.Println("clear offline sendbox msgs:", purgeErr, msg.MsgTime, msgExpired)
			return purgeErr
		})
	}
	if err != nil {
		logs.NewLogEntity().Errorf("failed to store sendbox. err:%v", err)
	}
	return err
}

func MsgDirect(ctx context.Context, targetId string, downMsg *pbobjs.DownMsg) {
	rpcMsg := bases.CreateServerPubWraper(ctx, bases.GetRequesterIdFromCtx(ctx), targetId, "msg", downMsg)
	if downMsg.IsSend {
		rpcMsg.PublishType = int32(commonservices.PublishType_AllSessionExceptSelf)
	}
	rpcMsg.Qos = 0
	bases.UnicastRouteWithNoSender(rpcMsg)
}

func purgeMsgs(key string, currentTime, cutoff int64, f purgeFunc) {
	if configures.Config.MsgStoreEngine == "" || configures.Config.MsgStoreEngine == configures.MsgStoreEngine_MySQL {
		msgPurgeOnce.Do(func() {
			msgPurgeScheduler = newPurgeScheduler(msgPurgeWorkerCount, msgPurgeQueueCapacity, msgPurgeInterval)
		})
		msgPurgeScheduler.Submit(key, currentTime, cutoff, f)
	} else {
		fmt.Println(configures.Config.MsgStoreEngine, ":", configures.MsgStoreEngine_MySQL)
	}
}

func getMsgExpired(appkey string, cmd bool) int64 {
	expired := configures.MsgExpired
	appinfo, exist := commonservices.GetAppInfo(appkey)
	if cmd {
		expired = configures.CmdMsgExpired
		if exist && appinfo != nil {
			expired = int64(appinfo.OfflineCmdMsgSaveTime) * 60 * 1000
		}
	} else if exist && appinfo != nil {
		expired = int64(appinfo.OfflineMsgSaveTime) * 60 * 1000
	}
	return expired
}

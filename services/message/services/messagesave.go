package services

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/caches"
	"im-server/commons/configures"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/commonservices"
	"im-server/services/commonservices/logs"
	"im-server/services/commonservices/msgdefines"
	"im-server/services/message/storages"
	"im-server/services/message/storages/models"
	"time"

	"github.com/Jeffail/tunny"
)

var pool *tunny.Pool
var taskCache *caches.LruCache

// TODO save immediately when user online, other wise, use async queue.
func SaveMsg2Inbox(appkey, receiverId string, msg *pbobjs.DownMsg) error {
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
		purgeMsgs(appkey+"1", msg.MsgTime, func() {
			cmdMsgExpired := configures.CmdMsgExpired
			appinfo, exist := commonservices.GetAppInfo(appkey)
			if exist && appinfo != nil {
				cmdMsgExpired = int64(appinfo.OfflineCmdMsgSaveTime) * 60 * 1000
			}
			msgStorage.DelMsgsBaseTime(appkey, msg.MsgTime-cmdMsgExpired)
		})
	} else {
		msgStorage := storages.NewInboxMsgStorage()
		err = msgStorage.SaveMsg(message)
		purgeMsgs(appkey+"2", msg.MsgTime, func() {
			msgExpired := configures.MsgExpired
			appinfo, exist := commonservices.GetAppInfo(appkey)
			if exist && appinfo != nil {
				msgExpired = int64(appinfo.OfflineMsgSaveTime) * 60 * 1000
			}
			msgStorage.DelMsgsBaseTime(appkey, msg.MsgTime-msgExpired)
		})
	}
	if err != nil {
		logs.NewLogEntity().Errorf("failed to store inbox. err:%v", err)
	}
	return err
}

func SaveMsg2Sendbox(ctx context.Context, appkey, senderId string, msg *pbobjs.DownMsg) error {
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
		purgeMsgs(appkey+"3", msg.MsgTime, func() {
			cmdMsgExpired := configures.CmdMsgExpired
			appinfo, exist := commonservices.GetAppInfo(appkey)
			if exist && appinfo != nil {
				cmdMsgExpired = int64(appinfo.OfflineCmdMsgSaveTime) * 60 * 1000
			}
			storage.DelMsgsBaseTime(appkey, msg.MsgTime-cmdMsgExpired)
		})
	} else {
		storage := storages.NewSendboxMsgStorage()
		err = storage.SaveMsg(message)
		purgeMsgs(appkey+"4", msg.MsgTime, func() {
			msgExpired := configures.MsgExpired
			appinfo, exist := commonservices.GetAppInfo(appkey)
			if exist && appinfo != nil {
				msgExpired = int64(appinfo.OfflineMsgSaveTime) * 60 * 1000
			}
			storage.DelMsgsBaseTime(appkey, msg.MsgTime-msgExpired)
		})
	}
	if err != nil {
		logs.NewLogEntity().Errorf("failed to store inbox. err:%v", err)
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

func purgeMsgs(key string, currentTime int64, f func()) {
	if configures.Config.MsgStoreEngine == "" || configures.Config.MsgStoreEngine == configures.MsgStoreEngine_MySQL {
		if pool == nil || taskCache == nil {
			pool = tunny.NewCallback(64)
			taskCache = caches.NewLruCacheWithReadTimeout("msgpure_cache", 10000, nil, 30*time.Minute)
		}
		if val, exist := taskCache.Get(key); exist {
			time := val.(int64)
			if currentTime-time > 60*1000 {
				go pool.Process(f)
				taskCache.Add(key, currentTime)
			}
		} else {
			taskCache.Add(key, currentTime)
		}
	}
}

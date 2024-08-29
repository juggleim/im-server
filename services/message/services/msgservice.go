package services

import (
	"context"
	"time"

	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/commonservices"
	"im-server/services/commonservices/logs"

	"google.golang.org/protobuf/proto"
)

var MsgSinglePools *tools.SinglePools

func init() {
	MsgSinglePools = tools.NewSinglePools(512)
}

func SendPrivateMsg(ctx context.Context, senderId, receiverId string, upMsg *pbobjs.UpMsg) (errs.IMErrorCode, string, int64, int64) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	converId := commonservices.GetConversationId(senderId, receiverId, pbobjs.ChannelType_Private)
	//statistic
	commonservices.ReportUpMsg(appkey, pbobjs.ChannelType_Private, 1)
	commonservices.ReportDispatchMsg(appkey, pbobjs.ChannelType_Private, 1)
	//check block user
	blockUsers := GetBlockUsers(appkey, receiverId)
	if blockUsers.CheckBlockUser(senderId) {
		sendTime := time.Now().UnixMilli()
		msgId := tools.GenerateMsgId(sendTime, int32(pbobjs.ChannelType_Private), receiverId)
		return errs.IMErrorCode_MSG_BLOCK, msgId, sendTime, 0
	}
	//check msg interceptor
	if code := commonservices.CheckMsgInterceptor(ctx, senderId, receiverId, pbobjs.ChannelType_Private, upMsg); code != errs.IMErrorCode_SUCCESS {
		sendTime := time.Now().UnixMilli()
		msgId := tools.GenerateMsgId(sendTime, int32(pbobjs.ChannelType_Private), receiverId)
		return code, msgId, sendTime, 0
	}
	msgConverCache := commonservices.GetMsgConverCache(ctx, converId, pbobjs.ChannelType_Private)
	msgId, sendTime, msgSeq := msgConverCache.GenerateMsgId(converId, pbobjs.ChannelType_Private, time.Now().UnixMilli(), upMsg.Flags)
	downMsg4Sendbox := &pbobjs.DownMsg{
		SenderId:       senderId,
		TargetId:       receiverId,
		ChannelType:    pbobjs.ChannelType_Private,
		MsgType:        upMsg.MsgType,
		MsgId:          msgId,
		MsgSeqNo:       msgSeq,
		MsgContent:     upMsg.MsgContent,
		MsgTime:        sendTime,
		Flags:          upMsg.Flags,
		ClientUid:      upMsg.ClientUid,
		IsSend:         true,
		MentionInfo:    upMsg.MentionInfo,
		ReferMsg:       commonservices.FillReferMsg(ctx, upMsg),
		TargetUserInfo: commonservices.GetTargetDisplayUserInfo(ctx, receiverId),
		MergedMsgs:     upMsg.MergedMsgs,
	}
	//send to sender's other device
	if !commonservices.IsStateMsg(upMsg.Flags) {
		//save msg to sendbox for sender
		//record conversation for sender
		commonservices.Save2Sendbox(ctx, downMsg4Sendbox)
	}
	if bases.GetOnlySendboxFromCtx(ctx) {
		return errs.IMErrorCode_SUCCESS, msgId, sendTime, msgSeq
	}
	commonservices.SubPrivateMsg(ctx, msgId, downMsg4Sendbox)

	downMsg := &pbobjs.DownMsg{
		SenderId:       senderId,
		TargetId:       senderId,
		ChannelType:    pbobjs.ChannelType_Private,
		MsgType:        upMsg.MsgType,
		MsgId:          msgId,
		MsgSeqNo:       msgSeq,
		MsgContent:     upMsg.MsgContent,
		MsgTime:        sendTime,
		Flags:          upMsg.Flags,
		ClientUid:      upMsg.ClientUid,
		MentionInfo:    upMsg.MentionInfo,
		ReferMsg:       commonservices.FillReferMsg(ctx, upMsg),
		TargetUserInfo: commonservices.GetSenderUserInfo(ctx),
		MergedMsgs:     upMsg.MergedMsgs,
	}

	//check merged msg
	if commonservices.IsMergedMsg(upMsg.Flags) && upMsg.MergedMsgs != nil && len(upMsg.MergedMsgs.Msgs) > 0 {
		bases.AsyncRpcCall(ctx, "merge_msgs", msgId, &pbobjs.MergeMsgReq{
			ParentMsgId: msgId,
			MergedMsgs:  upMsg.MergedMsgs,
		})
	}

	//save history msg
	if commonservices.IsStoreMsg(upMsg.Flags) {
		commonservices.SaveHistoryMsg(ctx, senderId, receiverId, pbobjs.ChannelType_Private, downMsg, 0)
	}

	//dispatch to receiver
	if senderId != receiverId {
		dispatchMsg(ctx, receiverId, downMsg)
	}

	return errs.IMErrorCode_SUCCESS, msgId, sendTime, msgSeq
}

func dispatchMsg(ctx context.Context, receiverId string, msg *pbobjs.DownMsg) {
	data, _ := tools.PbMarshal(msg)
	bases.UnicastRouteWithNoSender(&pbobjs.RpcMessageWraper{
		RpcMsgType:   pbobjs.RpcMsgType_UserPub,
		AppKey:       bases.GetAppKeyFromCtx(ctx),
		Session:      bases.GetSessionFromCtx(ctx),
		Method:       "msg_dispatch",
		RequesterId:  bases.GetRequesterIdFromCtx(ctx),
		ReqIndex:     bases.GetSeqIndexFromCtx(ctx),
		Qos:          bases.GetQosFromCtx(ctx),
		AppDataBytes: data,
		TargetId:     receiverId,
	})
}

func MsgOrNtf(ctx context.Context, targetId string, downMsg *pbobjs.DownMsg) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	userStatus := GetUserStatus(appkey, targetId)
	if userStatus.IsOnline() {
		isNtf := GetUserStatus(appkey, targetId).CheckNtfWithSwitch()
		if isNtf { //发送通知
			logs.WithContext(ctx).Infof("ntf target_id:%s", targetId)
			rpcNtf := bases.CreateServerPubWraper(ctx, bases.GetRequesterIdFromCtx(ctx), targetId, "ntf", &pbobjs.Notify{
				Type:     pbobjs.NotifyType_Msg,
				SyncTime: downMsg.MsgTime,
			})
			rpcNtf.Qos = 0
			bases.UnicastRouteWithNoSender(rpcNtf)
			bases.UnicastRouteWithCallback(rpcNtf, &SendMsgAckActor{
				appkey:      appkey,
				senderId:    bases.GetRequesterIdFromCtx(ctx),
				targetId:    targetId,
				channelType: downMsg.ChannelType,
				Msg:         downMsg,
				ctx:         ctx,
				HasPush:     userStatus.OpenPushSwitch(),
				IsNotify:    isNtf,
			}, 5*time.Second)
		} else {
			logs.WithContext(ctx).Infof("msg target_id:%s", targetId)
			rpcMsg := bases.CreateServerPubWraper(ctx, bases.GetRequesterIdFromCtx(ctx), targetId, "msg", downMsg)
			bases.UnicastRouteWithCallback(rpcMsg, &SendMsgAckActor{
				appkey:      appkey,
				senderId:    bases.GetRequesterIdFromCtx(ctx),
				targetId:    targetId,
				channelType: downMsg.ChannelType,
				Msg:         downMsg,
				ctx:         ctx,
				HasPush:     userStatus.OpenPushSwitch(),
				IsNotify:    isNtf,
			}, 5*time.Second)
		}
		if userStatus.OpenPushSwitch() {
			SendPush(ctx, bases.GetRequesterIdFromCtx(ctx), targetId, downMsg)
		}
	} else { //for push
		SendPush(ctx, bases.GetRequesterIdFromCtx(ctx), targetId, downMsg)
	}
}

func getPushLanguage(ctx context.Context, userId string) string {
	appkey := bases.GetAppKeyFromCtx(ctx)
	language := "en_US"
	appinfo, exist := commonservices.GetAppInfo(appkey)
	if exist {
		language = appinfo.PushLanguage
	}
	uSetting := commonservices.GetTargetUserSettings(ctx, userId)
	if uSetting != nil && uSetting.Language != "" {
		language = uSetting.Language
	}
	return language
}

type SendMsgAckActor struct {
	actorsystem.UntypedActor
	appkey      string
	senderId    string
	targetId    string
	channelType pbobjs.ChannelType
	// pushData    *pbobjs.PushData
	Msg      *pbobjs.DownMsg
	ctx      context.Context
	HasPush  bool
	IsNotify bool
}

func (actor *SendMsgAckActor) OnReceive(ctx context.Context, input proto.Message) {
	if rpcMsg, ok := input.(*pbobjs.RpcMessageWraper); ok {
		data := rpcMsg.AppDataBytes
		onlineStatus := &pbobjs.OnlineStatus{}
		err := tools.PbUnMarshal(data, onlineStatus)
		if err == nil {
			logs.WithContext(actor.ctx).Infof("target_id:%s\tonline_type:%d", actor.targetId, onlineStatus.Type)
			if onlineStatus.Type == pbobjs.OnlineType_Offline { //receiver is offline
				RecordUserOnlineStatus(actor.appkey, actor.targetId, false, 0)
				if !actor.HasPush {
					SendPush(actor.ctx, actor.senderId, actor.targetId, actor.Msg)
				}
			} else { //receiver is online, and response ack
				if !actor.IsNotify {
					us := GetUserStatus(actor.appkey, actor.targetId)
					if us != nil {
						us.SetNtfStatus(false)
					}
					//statistic
					commonservices.ReportDownMsg(actor.appkey, actor.channelType, 1)
				}
			}
		}
	}
}

func GetPushDataForDefaultMsg(msg *pbobjs.DownMsg, pushLanguage string) *pbobjs.PushData {
	if msg == nil {
		return nil
	}
	var retPushData *pbobjs.PushData
	if msg.PushData != nil && msg.PushData.PushText != "" {
		retPushData = msg.PushData
	} else {
		retPushData = &pbobjs.PushData{}
		var (
			title  string
			prefix string
		)
		nickName := msg.TargetUserInfo.GetNickname()
		if msg.ChannelType == pbobjs.ChannelType_Group {
			title = msg.GroupInfo.GroupName
			if nickName != "" {
				prefix = nickName + ": "
			}
		} else {
			title = nickName
		}
		retPushData.Title = title
		if msg.MsgType == "jg:text" {
			txtMsg := &commonservices.TextMsg{}
			err := tools.JsonUnMarshal(msg.MsgContent, txtMsg)

			pushText := txtMsg.Content
			if len(pushText) > 20 {
				pushText = pushText[:20] + "..."
			}
			if err == nil {
				retPushData.PushText = prefix + pushText
			} else {
				retPushData.PushText = prefix + "[Text]"
			}
		} else if msg.MsgType == "jg:img" {
			retPushData.PushText = prefix + "[Image]"
		} else if msg.MsgType == "jg:voice" {
			retPushData.PushText = prefix + "[Voice]"
		} else if msg.MsgType == "jg:file" {
			retPushData.PushText = prefix + "[File]"
		} else if msg.MsgType == "jg:video" {
			retPushData.PushText = prefix + "[Video]"
		} else if msg.MsgType == "jg:merge" {
			retPushData.PushText = prefix + "[Merge]"
		}
	}
	//add internal fields
	retPushData.MsgId = msg.MsgId
	retPushData.SenderId = msg.SenderId
	retPushData.ConverId = msg.TargetId
	retPushData.ChannelType = msg.ChannelType
	return retPushData
}

func (actor *SendMsgAckActor) CreateInputObj() proto.Message {
	return &pbobjs.RpcMessageWraper{}
}
func (actor *SendMsgAckActor) OnTimeout() {

}

func SendPush(ctx context.Context, senderId, receiverId string, msg *pbobjs.DownMsg) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	appInfo, exist := commonservices.GetAppInfo(appkey)
	if exist && appInfo != nil && appInfo.IsOpenPush {
		if !commonservices.IsUndisturbMsg(msg.Flags) {
			pushData := GetPushDataForDefaultMsg(msg, getPushLanguage(ctx, receiverId))
			if pushData != nil {
				pushRpc := bases.CreateServerPubWraper(ctx, senderId, receiverId, "push", pushData)
				bases.UnicastRouteWithNoSender(pushRpc)
			}
		}
	}
}

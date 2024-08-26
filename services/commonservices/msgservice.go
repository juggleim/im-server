package commonservices

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
)

/*
	Flags:
	1: is cmd msg
	2: is count msg
	4: is state msg
	8: is store msg
	16: is modified msg
	32: is merged msg
	64: is undisturb msg
	128: is no affect sender's conversation msg
	256: is ext msg
	512: is reaction msg
	1024: is stream msg
*/

func IsCmdMsg(flag int32) bool {
	return (flag & 0x1) == 1
}

func SetCmdMsg(flag int32) int32 {
	return flag | 0x1
}

func IsCountMsg(flag int32) bool {
	return (flag & 0x2) == 2
}

func SetCountMsg(flag int32) int32 {
	return flag | 0x2
}

func IsStateMsg(flag int32) bool {
	return (flag & 0x4) == 4
}

func SetStateMsg(flag int32) int32 {
	return flag | 0x4
}

func IsStoreMsg(flag int32) bool {
	return (flag & 0x8) == 8
}

func SetStoreMsg(flag int32) int32 {
	return flag | 0x8
}

func IsModifiedMsg(flag int32) bool {
	return (flag & 0x10) == 16
}

func SetModifiedMsg(flag int32) int32 {
	return flag | 0x10
}

func IsMergedMsg(flag int32) bool {
	return (flag & 0x20) == 32
}

func SetMergedMsg(flag int32) int32 {
	return flag | 0x20
}

func IsUndisturbMsg(flag int32) bool {
	return (flag & 0x40) == 64
}

func SetUndisturbMsg(flag int32) int32 {
	return flag | 0x40
}

func IsNoAffectSenderConver(flag int32) bool {
	return (flag & 0x80) == 128
}

func SetNoAffectSenderConver(flag int32) int32 {
	return flag | 0x80
}

func IsExtMsg(flag int32) bool {
	return (flag & 0x100) == 256
}

func SetExtMsg(flag int32) int32 {
	return flag | 0x100
}

func IsReactionMsg(flag int32) bool {
	return (flag & 0x200) == 512
}

func SetReactionMsg(flag int32) int32 {
	return flag | 0x200
}

func IsStreamMsg(flag int32) bool {
	return (flag & 0x400) == 1024
}

func SetStreamMsg(flag int32) int32 {
	return flag | 0x400
}

type PublishType int

var (
	PublishType_AllSession           = 0
	PublishType_OnlineSelfSession    = 1
	PublishType_AllSessionExceptSelf = 2
)

func FillReferMsg(ctx context.Context, upMsg *pbobjs.UpMsg) *pbobjs.DownMsg {
	if upMsg.ReferMsg != nil {
		if upMsg.ReferMsg.SenderId != "" {
			upMsg.ReferMsg.TargetUserInfo = GetUserInfoFromRpc(ctx, upMsg.ReferMsg.SenderId)
		}
		return upMsg.ReferMsg
	}
	return nil
}

func Save2Sendbox(ctx context.Context, downMsg *pbobjs.DownMsg) {
	noSendbox := bases.GetNoSendboxFromCtx(ctx)
	if !noSendbox {
		data, _ := tools.PbMarshal(downMsg)
		bases.UnicastRouteWithNoSender(&pbobjs.RpcMessageWraper{
			RpcMsgType:   pbobjs.RpcMsgType_UserPub,
			AppKey:       bases.GetAppKeyFromCtx(ctx),
			Session:      bases.GetSessionFromCtx(ctx),
			Method:       "sendbox",
			RequesterId:  bases.GetRequesterIdFromCtx(ctx),
			ReqIndex:     bases.GetSeqIndexFromCtx(ctx),
			Qos:          bases.GetQosFromCtx(ctx),
			AppDataBytes: data,
			TargetId:     bases.GetRequesterIdFromCtx(ctx),
		})
	}
}

func StreamMsgDirect(ctx context.Context, targetId string, streamMsg *pbobjs.StreamDownMsg) {
	rpcMsg := bases.CreateServerPubWraper(ctx, bases.GetRequesterIdFromCtx(ctx), targetId, "stream_msg", streamMsg)
	rpcMsg.Qos = 0
	bases.UnicastRouteWithNoSender(rpcMsg)
}

func SyncPrivateMsg(ctx context.Context, targetId string, upMsg *pbobjs.UpMsg, opts ...bases.BaseActorOption) (errs.IMErrorCode, string, int64, int64) {
	tmpOpts := []bases.BaseActorOption{}
	tmpOpts = append(tmpOpts, &bases.ExtsOption{Exts: map[string]string{
		RpcExtKey_RealTargetId: targetId,
	}})
	tmpOpts = append(tmpOpts, opts...)
	targetId = GetConversationId(bases.GetRequesterIdFromCtx(ctx), targetId, pbobjs.ChannelType_Private)
	result, err := bases.SyncOriginalRpcCall(ctx, "p_msg", targetId, upMsg, tmpOpts...)
	if err != nil || result == nil {
		return errs.IMErrorCode_DEFAULT, "", 0, 0
	}
	return errs.IMErrorCode(result.ResultCode), result.MsgId, result.MsgSendTime, result.MsgSeqNo
}

func AsyncPrivateMsg(ctx context.Context, userId, targetId string, upMsg *pbobjs.UpMsg, opts ...bases.BaseActorOption) {
	ctx = bases.SetRequesterId2Ctx(ctx, userId)
	tmpOpts := []bases.BaseActorOption{}
	tmpOpts = append(tmpOpts, &bases.ExtsOption{Exts: map[string]string{
		RpcExtKey_RealTargetId: targetId,
	}})
	tmpOpts = append(tmpOpts, opts...)
	targetId = GetConversationId(userId, targetId, pbobjs.ChannelType_Private)
	bases.AsyncRpcCall(ctx, "p_msg", targetId, upMsg, tmpOpts...)
}

func SyncGroupMsg(ctx context.Context, targetId string, upMsg *pbobjs.UpMsg, opts ...bases.BaseActorOption) (errs.IMErrorCode, string, int64, int64) {
	resp, err := bases.SyncOriginalRpcCall(ctx, "g_msg", targetId, upMsg, opts...)
	if err != nil || resp == nil {
		return errs.IMErrorCode_DEFAULT, "", 0, 0
	}
	return errs.IMErrorCode(resp.ResultCode), resp.MsgId, resp.MsgSendTime, resp.MsgSeqNo
}

func AsyncGroupMsg(ctx context.Context, userId, groupId string, upMsg *pbobjs.UpMsg, opts ...bases.BaseActorOption) {
	ctx = bases.SetRequesterId2Ctx(ctx, userId)
	bases.AsyncRpcCall(ctx, "g_msg", groupId, upMsg, opts...)
}

func GroupMsgFromApi(ctx context.Context, userId, groupId string, upMsg *pbobjs.UpMsg, noSendbox bool) {
	data, _ := tools.PbMarshal(upMsg)
	bases.UnicastRouteWithNoSender(&pbobjs.RpcMessageWraper{
		RpcMsgType:   pbobjs.RpcMsgType_UserPub,
		AppKey:       bases.GetAppKeyFromCtx(ctx),
		Session:      bases.GetSessionFromCtx(ctx),
		Method:       "g_msg",
		RequesterId:  userId,
		ReqIndex:     bases.GetSeqIndexFromCtx(ctx),
		Qos:          bases.GetQosFromCtx(ctx),
		AppDataBytes: data,
		TargetId:     groupId,
		IsFromApi:    true,
		NoSendbox:    noSendbox,
	})
}

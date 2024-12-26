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

type PublishType int32

var (
	PublishType_AllSession           PublishType = 0
	PublishType_OnlineSelfSession    PublishType = 1
	PublishType_AllSessionExceptSelf PublishType = 2
)

func FillReferMsg(ctx context.Context, upMsg *pbobjs.UpMsg) *pbobjs.DownMsg {
	if upMsg.ReferMsg != nil {
		if upMsg.ReferMsg.SenderId != "" {
			upMsg.ReferMsg.TargetUserInfo = GetTargetDisplayUserInfo(ctx, upMsg.ReferMsg.SenderId)
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
			IsFromApi:    bases.GetIsFromApiFromCtx(ctx),
			ExtParams:    bases.GetExtsFromCtx(ctx),
		})
	}
}

func SyncMsg(ctx context.Context, method, userId, targetId string, upMsg *pbobjs.UpMsg, opts ...bases.BaseActorOption) (errs.IMErrorCode, string, int64, int64) {
	tmpOpts := []bases.BaseActorOption{}
	if method == "p_msg" || method == "imp_pri_msg" {
		tmpOpts = append(tmpOpts, &bases.ExtsOption{Exts: map[string]string{
			RpcExtKey_RealTargetId: targetId,
		}})
		tmpOpts = append(tmpOpts, opts...)
		targetId = GetConversationId(userId, targetId, pbobjs.ChannelType_Private)
	} else if method == "s_msg" {
		tmpOpts = append(tmpOpts, &bases.ExtsOption{Exts: map[string]string{
			RpcExtKey_RealTargetId: targetId,
		}})
		tmpOpts = append(tmpOpts, opts...)
		targetId = GetConversationId(userId, targetId, pbobjs.ChannelType_System)
	}
	ctx = bases.SetRequesterId2Ctx(ctx, userId)
	result, err := bases.SyncOriginalRpcCall(ctx, method, targetId, upMsg, tmpOpts...)
	if err != nil || result == nil {
		return errs.IMErrorCode_DEFAULT, "", 0, 0
	}
	return errs.IMErrorCode(result.ResultCode), result.MsgId, result.MsgSendTime, result.MsgSeqNo
}

func SyncMsgOverUpstream(ctx context.Context, method, userId, targetId string, upMsg *pbobjs.UpMsg, opts ...bases.BaseActorOption) (errs.IMErrorCode, string, int64, int64) {
	tmpOpts := []bases.BaseActorOption{}
	tmpOpts = append(tmpOpts, &bases.ExtsOption{
		Exts: map[string]string{
			RpcExtKey_RealMethod: method,
		},
	})
	requestId := userId
	tmpOpts = append(tmpOpts, opts...)
	ctx = bases.SetRequesterId2Ctx(ctx, targetId)
	result, err := bases.SyncOriginalRpcCall(ctx, "upstream", requestId, upMsg, tmpOpts...)
	if err != nil || result == nil {
		return errs.IMErrorCode_DEFAULT, "", 0, 0
	}
	return errs.IMErrorCode(result.ResultCode), result.MsgId, result.MsgSendTime, result.MsgSeqNo
}

func AsyncMsg(ctx context.Context, method, userId, targetId string, upMsg *pbobjs.UpMsg, opts ...bases.BaseActorOption) {
	ctx = bases.SetRequesterId2Ctx(ctx, userId)
	tmpOpts := []bases.BaseActorOption{}
	if method == "p_msg" || method == "imp_pri_msg" {
		tmpOpts = append(tmpOpts, &bases.ExtsOption{Exts: map[string]string{
			RpcExtKey_RealTargetId: targetId,
		}})
		tmpOpts = append(tmpOpts, opts...)
		targetId = GetConversationId(userId, targetId, pbobjs.ChannelType_Private)
	} else if method == "s_msg" {
		tmpOpts = append(tmpOpts, &bases.ExtsOption{Exts: map[string]string{
			RpcExtKey_RealTargetId: targetId,
		}})
		tmpOpts = append(tmpOpts, opts...)
		targetId = GetConversationId(userId, targetId, pbobjs.ChannelType_System)
	}
	bases.AsyncRpcCall(ctx, method, targetId, upMsg, tmpOpts...)
}

func AsyncMsgOverUpstream(ctx context.Context, method, userId, targetId string, upMsg *pbobjs.UpMsg, opts ...bases.BaseActorOption) {
	tmpOpts := []bases.BaseActorOption{}
	tmpOpts = append(tmpOpts, &bases.ExtsOption{
		Exts: map[string]string{
			RpcExtKey_RealMethod: method,
		},
	})
	tmpOpts = append(tmpOpts, opts...)
	ctx = bases.SetRequesterId2Ctx(ctx, targetId)
	bases.AsyncRpcCall(ctx, "upstream", userId, upMsg, tmpOpts...)
}

func SyncPrivateMsg(ctx context.Context, userId, targetId string, upMsg *pbobjs.UpMsg, opts ...bases.BaseActorOption) (errs.IMErrorCode, string, int64, int64) {
	return SyncMsg(ctx, "p_msg", userId, targetId, upMsg, opts...)
}

func SyncPrivateMsgOverUpstream(ctx context.Context, userId, targetId string, upMsg *pbobjs.UpMsg, opts ...bases.BaseActorOption) (errs.IMErrorCode, string, int64, int64) {
	return SyncMsgOverUpstream(ctx, "p_msg", userId, targetId, upMsg, opts...)
}

func AsyncPrivateMsg(ctx context.Context, userId, targetId string, upMsg *pbobjs.UpMsg, opts ...bases.BaseActorOption) {
	AsyncMsg(ctx, "p_msg", userId, targetId, upMsg, opts...)
}

func AsyncPrivateMsgOverUpstream(ctx context.Context, userId, targetId string, upMsg *pbobjs.UpMsg, opts ...bases.BaseActorOption) {
	AsyncMsgOverUpstream(ctx, "p_msg", userId, targetId, upMsg, opts...)
}

func AsyncImportPrivateMsgOverUpstream(ctx context.Context, userId, targetId string, upMsg *pbobjs.UpMsg, opts ...bases.BaseActorOption) {
	AsyncMsgOverUpstream(ctx, "imp_pri_msg", userId, targetId, upMsg, opts...)
}

func SyncSystemMsg(ctx context.Context, systemId, targetId string, upMsg *pbobjs.UpMsg, opts ...bases.BaseActorOption) (errs.IMErrorCode, string, int64, int64) {
	return SyncMsg(ctx, "s_msg", systemId, targetId, upMsg, opts...)
}

func SyncSystemMsgOverUpstream(ctx context.Context, systemId, targetId string, upMsg *pbobjs.UpMsg, opts ...bases.BaseActorOption) (errs.IMErrorCode, string, int64, int64) {
	return SyncMsgOverUpstream(ctx, "s_msg", systemId, targetId, upMsg, opts...)
}

func AsyncSystemMsg(ctx context.Context, systemId, targetId string, upMsg *pbobjs.UpMsg, opts ...bases.BaseActorOption) {
	AsyncMsg(ctx, "s_msg", systemId, targetId, upMsg, opts...)
}

func AsyncSystemMsgOverUpstream(ctx context.Context, systemId, targetId string, upMsg *pbobjs.UpMsg, opts ...bases.BaseActorOption) {
	AsyncMsgOverUpstream(ctx, "s_msg", systemId, targetId, upMsg, opts...)
}

func SyncGroupMsg(ctx context.Context, userId, targetId string, upMsg *pbobjs.UpMsg, opts ...bases.BaseActorOption) (errs.IMErrorCode, string, int64, int64) {
	return SyncMsg(ctx, "g_msg", userId, targetId, upMsg, opts...)
}

func SyncGroupMsgOverUpstream(ctx context.Context, userId, targetId string, upMsg *pbobjs.UpMsg, opts ...bases.BaseActorOption) (errs.IMErrorCode, string, int64, int64) {
	return SyncMsgOverUpstream(ctx, "g_msg", userId, targetId, upMsg, opts...)
}

func AsyncGroupMsg(ctx context.Context, userId, groupId string, upMsg *pbobjs.UpMsg, opts ...bases.BaseActorOption) {
	AsyncMsg(ctx, "g_msg", userId, groupId, upMsg, opts...)
}

func AsyncGroupMsgOverUpstream(ctx context.Context, userId, groupId string, upMsg *pbobjs.UpMsg, opts ...bases.BaseActorOption) {
	AsyncMsgOverUpstream(ctx, "g_msg", userId, groupId, upMsg, opts...)
}

func AsyncImportGroupMsgOverUpstream(ctx context.Context, userId, groupId string, upMsg *pbobjs.UpMsg, opts ...bases.BaseActorOption) {
	AsyncMsgOverUpstream(ctx, "imp_grp_msg", userId, groupId, upMsg, opts...)
}

func AsyncChatMsg(ctx context.Context, userId, chatId string, upMsg *pbobjs.UpMsg, opts ...bases.BaseActorOption) {
	AsyncMsg(ctx, "c_msg", userId, chatId, upMsg, opts...)
}

func AsyncChatMsgOverUpstream(ctx context.Context, userId, chatId string, upMsg *pbobjs.UpMsg, opts ...bases.BaseActorOption) {
	AsyncMsgOverUpstream(ctx, "c_msg", userId, chatId, upMsg, opts...)
}

func IsMentionedMe(userId string, downMsg *pbobjs.DownMsg) bool {
	if downMsg != nil && downMsg.MentionInfo != nil {
		if downMsg.MentionInfo.MentionType == pbobjs.MentionType_All || downMsg.MentionInfo.MentionType == pbobjs.MentionType_AllAndSomeone {
			//mention all
			return true
		} else if downMsg.MentionInfo.MentionType == pbobjs.MentionType_Someone {
			for _, mentionedUser := range downMsg.MentionInfo.TargetUsers {
				if userId == mentionedUser.UserId {
					return true
				}
			}
		}
	}
	return false
}

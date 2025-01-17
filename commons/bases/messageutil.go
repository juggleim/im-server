package bases

import (
	"context"
	"im-server/commons/errs"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"

	"google.golang.org/protobuf/proto"
)

func CreateQueryWraper(ctx context.Context, method string, msg proto.Message) *pbobjs.RpcMessageWraper {
	wraper := &pbobjs.RpcMessageWraper{
		RpcMsgType: pbobjs.RpcMsgType_QueryMsg,
	}
	bs, _ := tools.PbMarshal(msg)
	wraper.AppDataBytes = bs
	handleBaseContext(ctx, wraper)
	wraper.Method = method
	return wraper
}

func CreateQueryAckWraper(ctx context.Context, code errs.IMErrorCode, respMsg proto.Message) *pbobjs.RpcMessageWraper {
	queryAck := &pbobjs.RpcMessageWraper{
		RpcMsgType: pbobjs.RpcMsgType_QueryAck,
		ResultCode: int32(code),
	}
	if respMsg != nil {
		bs, _ := tools.PbMarshal(respMsg)
		queryAck.AppDataBytes = bs
	}
	handleBaseContext(ctx, queryAck)

	return queryAck
}

func CreateQueryAckWraperWithTime(ctx context.Context, code errs.IMErrorCode, time int64, respMsg proto.Message) *pbobjs.RpcMessageWraper {
	queryAck := &pbobjs.RpcMessageWraper{
		RpcMsgType:  pbobjs.RpcMsgType_QueryAck,
		ResultCode:  int32(code),
		MsgSendTime: time,
	}
	if respMsg != nil {
		bs, _ := tools.PbMarshal(respMsg)
		queryAck.AppDataBytes = bs
	}
	handleBaseContext(ctx, queryAck)

	return queryAck
}

func CreateServerPubWraper(ctx context.Context, requesterId, targetId, method string, msg proto.Message) *pbobjs.RpcMessageWraper {
	serverPub := &pbobjs.RpcMessageWraper{
		RpcMsgType: pbobjs.RpcMsgType_ServerPub,
	}
	bs, _ := tools.PbMarshal(msg)
	serverPub.AppDataBytes = bs
	handleBaseContext(ctx, serverPub)
	serverPub.RequesterId = requesterId
	serverPub.TargetId = targetId
	serverPub.Method = method
	serverPub.Qos = 1
	return serverPub
}

func CreateUserPubAckWraper(ctx context.Context, code errs.IMErrorCode, msgId string, msgSendTime, msgSeqNo int64, clientMsgId string, modifiedMsg *pbobjs.DownMsg) *pbobjs.RpcMessageWraper {
	userPubAck := &pbobjs.RpcMessageWraper{
		RpcMsgType:  pbobjs.RpcMsgType_UserPubAck,
		ResultCode:  int32(code),
		MsgId:       msgId,
		MsgSendTime: msgSendTime,
		MsgSeqNo:    msgSeqNo,
		ClientMsgId: clientMsgId,
		ModifiedMsg: modifiedMsg,
	}
	handleBaseContext(ctx, userPubAck)
	return userPubAck
}

func CreateGrpPubAckWraper(ctx context.Context, code errs.IMErrorCode, msgId string, msgSendTime, msgSeqNo int64, clientMsgId string, memberCount int32, modifiedMsg *pbobjs.DownMsg) *pbobjs.RpcMessageWraper {
	userPubAck := &pbobjs.RpcMessageWraper{
		RpcMsgType:  pbobjs.RpcMsgType_UserPubAck,
		ResultCode:  int32(code),
		MsgId:       msgId,
		MsgSendTime: msgSendTime,
		MsgSeqNo:    msgSeqNo,
		MemberCount: memberCount,
		ClientMsgId: clientMsgId,
		ModifiedMsg: modifiedMsg,
	}
	handleBaseContext(ctx, userPubAck)
	return userPubAck
}

func handleBaseContext(ctx context.Context, rpcMsg *pbobjs.RpcMessageWraper) {
	rpcMsg.ReqIndex = GetSeqIndexFromCtx(ctx)
	rpcMsg.AppKey = GetAppKeyFromCtx(ctx)
	rpcMsg.Qos = GetQosFromCtx(ctx)
	rpcMsg.Session = GetSessionFromCtx(ctx)
	rpcMsg.Method = GetMethodFromCtx(ctx)
	rpcMsg.SourceMethod = GetSourceMethodFromCtx(ctx)
	rpcMsg.RequesterId = GetRequesterIdFromCtx(ctx)
	rpcMsg.TargetId = GetTargetIdFromCtx(ctx)
	rpcMsg.PublishType = GetPublishTypeFromCtx(ctx)
	rpcMsg.IsFromApi = GetIsFromApiFromCtx(ctx)
}

func Redirect(ctx context.Context, method, targetId string, data proto.Message, sender actorsystem.ActorRef) {
	msg := &pbobjs.RpcMessageWraper{}
	handleBaseContext(ctx, msg)
	msg.Method = method
	msg.TargetId = targetId
	appdata, _ := tools.PbMarshal(data)
	msg.AppDataBytes = appdata
	UnicastRouteWithSenderActor(msg, sender)
}

func CreateRpcCtx(appkey, requesterId string) context.Context {
	ctx := context.Background()
	ctx = setCtxValue(ctx, CtxKey_AppKey, appkey)
	ctx = setCtxValue(ctx, CtxKey_RequesterId, requesterId)
	ctx = setCtxValue(ctx, CtxKey_Session, tools.GenerateUUIDShort11())
	return ctx
}

package services

import (
	"time"

	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/apigateway/models"
	"im-server/services/commonservices"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/proto"
)

const (
	CtxKey_AppKey  string = "CtxKey_AppKey"
	CtxKey_Session string = "CtxKey_Session"
)

func AsyncSendMsg(ctx *gin.Context, method, requestId, targetId string, req proto.Message, isNotifySender bool) {
	dataBytes, _ := tools.PbMarshal(req)
	exts := map[string]string{}
	exts[commonservices.RpcExtKey_RealMethod] = method
	method = "upstream"

	bases.UnicastRouteWithNoSender(&pbobjs.RpcMessageWraper{
		RpcMsgType:   pbobjs.RpcMsgType_ServerPub,
		AppKey:       GetCtxString(ctx, CtxKey_AppKey),
		Session:      GetCtxString(ctx, CtxKey_Session),
		Method:       method,
		RequesterId:  targetId,
		Qos:          1,
		TargetId:     requestId,
		AppDataBytes: dataBytes,
		IsFromApi:    true,
		NoSendbox:    !isNotifySender,
		ExtParams:    exts,
	})
}

func SyncSendMsg(ctx *gin.Context, method, requestId, targetId string, req proto.Message, isNotifySender bool) (errs.IMErrorCode, *models.SendMsgResp, error) {
	dataBytes, _ := tools.PbMarshal(req)
	exts := map[string]string{}
	// exts[commonservices.RpcExtKey_RealMethod] = method
	// method = "upstream"
	if method == "p_msg" {
		exts[commonservices.RpcExtKey_RealTargetId] = targetId
		targetId = commonservices.GetConversationId(requestId, targetId, pbobjs.ChannelType_Private)
	} else if method == "s_msg" {
		exts[commonservices.RpcExtKey_RealTargetId] = targetId
		targetId = commonservices.GetConversationId(requestId, targetId, pbobjs.ChannelType_System)
	}
	result, err := bases.SyncUnicastRoute(&pbobjs.RpcMessageWraper{
		RpcMsgType:   pbobjs.RpcMsgType_ServerPub,
		AppKey:       GetCtxString(ctx, CtxKey_AppKey),
		Session:      GetCtxString(ctx, CtxKey_Session),
		Method:       method,
		RequesterId:  requestId,
		Qos:          1,
		TargetId:     targetId,
		AppDataBytes: dataBytes,
		IsFromApi:    true,
		NoSendbox:    !isNotifySender,
		ExtParams:    exts,
	}, 5*time.Second)
	if err != nil {
		return errs.IMErrorCode_API_INTERNAL_RESP_FAIL, nil, err
	}
	return errs.IMErrorCode(result.ResultCode), &models.SendMsgResp{
		MsgId:   result.MsgId,
		MsgTime: result.MsgSendTime,
		MsgSeq:  result.MsgSeqNo,
	}, nil
}

func AsyncApiCall(ctx *gin.Context, method, requestId, targetId string, req proto.Message) {
	dataBytes, _ := tools.PbMarshal(req)
	bases.UnicastRouteWithNoSender(&pbobjs.RpcMessageWraper{
		RpcMsgType:   pbobjs.RpcMsgType_ServerPub,
		AppKey:       GetCtxString(ctx, CtxKey_AppKey),
		Session:      GetCtxString(ctx, CtxKey_Session),
		Method:       method,
		RequesterId:  requestId,
		Qos:          1,
		TargetId:     targetId,
		AppDataBytes: dataBytes,
		IsFromApi:    true,
	})
}

func SyncApiCall(ctx *gin.Context, method, requestId, targetId string, req proto.Message, respFactory func() proto.Message) (int32, proto.Message, error) {
	dataBytes, _ := tools.PbMarshal(req)
	result, err := bases.SyncUnicastRoute(&pbobjs.RpcMessageWraper{
		RpcMsgType:   pbobjs.RpcMsgType_QueryMsg,
		AppKey:       GetCtxString(ctx, CtxKey_AppKey),
		Session:      GetCtxString(ctx, CtxKey_Session),
		Method:       method,
		RequesterId:  requestId,
		Qos:          1,
		AppDataBytes: dataBytes,
		TargetId:     targetId,
		IsFromApi:    true,
	}, 5*time.Second)
	if err != nil {
		return int32(errs.IMErrorCode_API_INTERNAL_RESP_FAIL), nil, err
	}
	if respFactory != nil {
		respObj := respFactory()
		err = tools.PbUnMarshal(result.AppDataBytes, respObj)
		if err != nil {
			return int32(errs.IMErrorCode_API_INTERNAL_RESP_FAIL), nil, err
		}
		return result.ResultCode, respObj, nil
	} else {
		return result.ResultCode, nil, nil
	}
}
func GetCtxString(ctx *gin.Context, key string) string {
	val, exist := ctx.Get(key)
	if exist {
		return val.(string)
	} else {
		return ""
	}
}

func GetAppkeyFromCtx(ctx *gin.Context) string {
	val, exist := ctx.Get(CtxKey_AppKey)
	if exist {
		return val.(string)
	} else {
		return ""
	}
}

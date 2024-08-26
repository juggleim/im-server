package services

import (
	"im-server/commons/bases"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/proto"
)

const (
	CtxKey_AppKey  string = "CtxKey_AppKey"
	CtxKey_Account string = "CtxKey_Account"
	CtxKey_Session string = "CtxKey_Session"
)

func SyncApiCall(ctx *gin.Context, method, requestId, targetId string, req proto.Message, respFactory func() proto.Message) (AdminErrorCode, interface{}, error) {
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
		return AdminErrorCode_ServerErr, nil, err
	}
	if respFactory != nil {
		respObj := respFactory()
		err = tools.PbUnMarshal(result.AppDataBytes, respObj)
		if err != nil {
			return AdminErrorCode_ServerErr, nil, err
		}
		return AdminErrorCode(result.ResultCode), respObj, nil
	} else {
		return AdminErrorCode(result.ResultCode), nil, nil
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

func SetCtxString(ctx *gin.Context, key string, val interface{}) {
	ctx.Set(key, val)
}

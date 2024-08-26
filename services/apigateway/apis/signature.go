package apis

import (
	"fmt"

	"im-server/commons/errs"
	"im-server/commons/tools"
	"im-server/services/apigateway/services"
	"im-server/services/commonservices"

	"github.com/gin-gonic/gin"
)

const (
	Header_AppKey    string = "appkey"
	Header_Nonce     string = "nonce"
	Header_Timestamp string = "timestamp"
	Header_Signature string = "signature"

	Header_RequestId string = "request-id"
)

func Signature(ctx *gin.Context) {
	session := fmt.Sprintf("api_%s", tools.GenerateUUIDShort11())
	ctx.Header(Header_RequestId, session)
	ctx.Set(services.CtxKey_Session, session)

	appKey := ctx.Request.Header.Get(Header_AppKey)
	nonce := ctx.Request.Header.Get(Header_Nonce)
	tsStr := ctx.Request.Header.Get(Header_Timestamp)
	signature := ctx.Request.Header.Get(Header_Signature)
	if appKey == "" {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_APPKEY_REQUIRED)
		ctx.Abort()
		return
	}
	if nonce == "" {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_NONCE_REQUIRED)
		ctx.Abort()
		return
	}
	if tsStr == "" {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_TIMESTAMP_REQUIRED)
		ctx.Abort()
		return
	}
	if signature == "" {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_SIGNATURE_REQUIRED)
		ctx.Abort()
		return
	}
	appInfo, exist := commonservices.GetAppInfo(appKey)
	if exist && appInfo != nil {
		str := fmt.Sprintf("%s%s%s", appInfo.AppSecret, nonce, tsStr)
		sig := tools.SHA1(str)
		if sig == signature {
			ctx.Set(services.CtxKey_AppKey, appKey)
		} else {
			tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_SIGNATURE_FAIL)
			ctx.Abort()
			return
		}
	} else {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_APP_NOT_EXISTED)
		ctx.Abort()
		return
	}
}

func TestSha1() {
	str := fmt.Sprintf("%s%s%s", "appsecret", "nonce", "1672568121910")
	sig := tools.SHA1(str)
	fmt.Println(sig)
}

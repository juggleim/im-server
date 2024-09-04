package apis

import (
	"im-server/commons/errs"
	"im-server/commons/tools"
	"im-server/services/commonservices"
	"im-server/services/commonservices/tokens"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	CtxKey_AppKey string = "CtxKey_AppKey"
	CtxKey_UserId string = "CtxKey_UserId"
)

func CheckToken(ctx *gin.Context) {
	appKey := ctx.GetHeader("x-appkey")
	tokenStr := ctx.GetHeader("x-token")

	userId, code := parseToken(appKey, tokenStr)
	if code != errs.IMErrorCode_SUCCESS {
		tools.ErrorHttpResp(ctx, code)
		ctx.Abort()
		return
	}
	ctx.Set(CtxKey_AppKey, appKey)
	ctx.Set(CtxKey_UserId, userId)
	ctx.Next()
}

func parseToken(appkey, tokenStr string) (string, errs.IMErrorCode) {
	tokenWrap, err := tokens.ParseTokenString(tokenStr)
	if err != nil {
		return "", errs.IMErrorCode_CONNECT_TOKEN_ILLEGAL
	}
	if tokenWrap.AppKey != appkey {
		return "", errs.IMErrorCode_CONNECT_TOKEN_AUTHFAIL
	}
	appInfo, exist := commonservices.GetAppInfo(appkey)
	if !exist || appInfo == nil {
		return "", errs.IMErrorCode_CONNECT_APP_NOT_EXISTED
	}
	token, err := tokens.ParseToken(tokenWrap, []byte(appInfo.AppSecureKey))
	if err != nil {
		return "", errs.IMErrorCode_CONNECT_TOKEN_AUTHFAIL
	}
	if appInfo.TokenEffectiveMinute > 0 && (token.TokenTime+int64(appInfo.TokenEffectiveMinute)*60*1000) < time.Now().UnixMilli() {
		return "", errs.IMErrorCode_CONNECT_TOKEN_EXPIRED
	}
	return token.UserId, errs.IMErrorCode_SUCCESS
}

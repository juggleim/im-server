package apis

import (
	"compress/gzip"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices"
	"im-server/services/commonservices/tokens"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func CorsHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Header("Access-Control-Allow-Origin", "*")
		context.Header("Access-Control-Allow-Credentials", "true")
		context.Header("Access-Control-Allow-Headers", "*")
		context.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		context.Header("Access-Control-Expose-Headers", "Content-Length")
		context.Header("Access-Control-Max-Age", "86400")
		method := context.Request.Method

		if method == "OPTIONS" {
			context.AbortWithStatus(http.StatusNoContent)
			return
		}
		context.Next()
	}
}

func GzipDecompress() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("Content-Encoding") == "gzip" {
			reader, err := gzip.NewReader(c.Request.Body)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
				c.Abort()
				return
			}
			c.Request.Body = reader
		}
		c.Next()
	}
}

const (
	CtxKey_AppKey string = "CtxKey_AppKey"
)

func CheckToken(ctx *gin.Context) {
	appKey := ctx.GetHeader("x-appkey")
	tokenStr := ctx.GetHeader("x-token")

	_, _, code := parseToken(appKey, tokenStr)
	if code != errs.IMErrorCode_SUCCESS {
		msg := errs.GetApiErrorByCode(code)
		FailHttpResp(ctx, ApiErrorCode_AuthFail, msg.Msg)
		ctx.Abort()
	}
	ctx.Set(CtxKey_AppKey, appKey)
	ctx.Next()
}

func parseToken(appKey string, tokenStr string) (appInfo *commonservices.AppInfo, token tokens.ImToken, code errs.IMErrorCode) {
	var (
		tokenWrap *pbobjs.TokenWrap
		err       error
		exist     bool
	)
	tokenWrap, err = tokens.ParseTokenString(tokenStr)
	if err != nil {
		code = errs.IMErrorCode_CONNECT_TOKEN_AUTHFAIL
		return
	}
	appInfo, exist = commonservices.GetAppInfo(appKey)
	if !exist || appInfo == nil {
		code = errs.IMErrorCode_CONNECT_APP_NOT_EXISTED
		return
	}
	token, err = tokens.ParseToken(tokenWrap, []byte(appInfo.AppSecureKey))
	if err != nil {
		code = errs.IMErrorCode_CONNECT_TOKEN_AUTHFAIL
		return
	}

	if appInfo.TokenEffectiveMinute > 0 && (token.TokenTime+int64(appInfo.TokenEffectiveMinute)*60*1000) < time.Now().UnixMilli() {
		code = errs.IMErrorCode_CONNECT_TOKEN_EXPIRED
		return
	}

	code = errs.IMErrorCode_SUCCESS
	return
}

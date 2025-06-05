package routers

import (
	"compress/gzip"
	"im-server/commons/errs"
	"im-server/commons/tools"
	"im-server/services/navigator/apis"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Route(eng *gin.Engine, prefix string) {
	group := eng.Group("/" + prefix)
	group.Use(apis.CheckToken)
	group.Use(corsHandler())
	group.GET("/general", apis.NaviGet)

	group.POST("/upload-log", apis.UploadClientLog)
	group.POST("/upload-log-plain", apis.UploadClientLogPlain, gzipDecompress())
	group.POST("/log-status", apis.UploadLogStatus)
}

func corsHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		method := context.Request.Method
		context.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		context.Writer.Header().Add("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, x-token, x-appkey, x-platform")
		context.Writer.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PATCH, PUT")
		context.Writer.Header().Add("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type, X-Token, X-Appid")
		context.Writer.Header().Add("Access-Control-Allow-Credentials", "true")

		if method == "OPTIONS" {
			context.AbortWithStatus(http.StatusNoContent)
		}
		context.Next()
	}
}

func gzipDecompress() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.GetHeader("Content-Encoding") == "gzip" {
			reader, err := gzip.NewReader(ctx.Request.Body)
			if err != nil {
				tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_PARAM_ILLEGAL)
				ctx.Abort()
				return
			}
			ctx.Request.Body = reader
		}
		ctx.Next()
	}
}

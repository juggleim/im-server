package routers

import (
	"im-server/services/admingateway/apis"
	"im-server/services/admingateway/ctxs"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Route(eng *gin.Engine, prefix string) *gin.RouterGroup {
	eng.Use(CorsHandler(), InjectCtx())
	group := eng.Group("/" + prefix)
	group.Use(apis.Validate)

	group.POST("/imapiagent", apis.ApiAgent)
	group.GET("/common/address", apis.GetAccessAddress)

	group.POST("/login", apis.Login)
	group.POST("/accounts/updpass", apis.UpdPassword)
	group.POST("/accounts/add", apis.AddAccount)
	group.POST("/accounts/delete", apis.DeleteAccounts)
	group.POST("/accounts/disable", apis.DisableAccounts)
	group.POST("/accounts/bindapps", apis.BindApps)
	group.POST("/accounts/unbindapps", apis.UnBindApps)
	group.GET("/accounts/list", apis.QryAccounts)

	group.POST("/apps/create", apis.CreateApp)
	group.GET("/apps/list", apis.QryApps)
	group.GET("/apps/info", apis.QryAppInfo)
	group.POST("/apps/configs/set", apis.UpdateAppConfigs)
	group.POST("/apps/configs/get", apis.QryAppConfigs)
	group.POST("/apps/eventsubconfig/set", apis.SetEventSubConfig)
	group.GET("/apps/eventsubconfig/get", apis.GetEventSubConfig)
	//translate
	group.POST("/apps/translate/set", apis.SetTranslateConf)
	group.GET("/apps/translate/get", apis.GetTranslateConf)
	//sms
	group.POST("/apps/sms/set", apis.SetSmsConf)
	group.GET("/apps/sms/get", apis.GetSmsConf)
	//rtc
	group.POST("/apps/rtcconf/set", apis.SetRtcConf)
	group.GET("/apps/rtcconf/get", apis.GetRtcConf)

	group.POST("/apps/iospushcer/set", apis.SetIosPushConf)
	group.POST("/apps/iospushcer/upload", apis.UploadIosCer)
	group.GET("/apps/iospushcer/get", apis.GetIosCer)
	group.POST("/apps/fcmpushconf/upload", apis.UploadFcmPushConf)
	group.GET("/apps/fcmpushconf/get", apis.GetFcmPushConf)
	group.POST("/apps/androidpushconf/set", apis.SetAndroidPushConf)
	group.GET("/apps/androidpushconf/get", apis.GetAndroidPushConf)

	group.POST("/apps/fileconf/set", apis.SetFileConf)
	group.GET("/apps/fileconf/get", apis.GetFileConf)
	group.GET("/apps/fileconf/switch/get", apis.GetFileConfs)
	group.POST("/apps/fileconf/switch/set", apis.SetFileConfSwitch)
	//logs
	group.POST("/apps/clientlogs/notify", apis.ClientLogNtf)
	group.GET("/apps/clientlogs/list", apis.ClientLogList)
	group.GET("/apps/clientlogs/download", apis.ClientLogDownload)
	group.GET("/apps/serverlogs/userconnect", apis.QryUserConnectLogs)
	group.GET("/apps/serverlogs/connect", apis.QryConnectLogs)
	group.GET("/apps/serverlogs/business", apis.QryBusinessLogs)

	//statistic
	group.GET("/apps/statistic/msg", apis.QryMsgStatistic)
	group.GET("/apps/statistic/useractivity", apis.QryUserActivities)
	group.GET("/apps/statistic/connectcount", apis.QryConnectCount)
	group.GET("/apps/statistic/userreg", apis.QryUserRegiste)
	group.GET("/apps/statistic/maxconnectcount", apis.QryMaxConnectCount)

	group.GET("/apps/sensitivewords/list", apis.SensitiveWords)
	group.POST("/apps/sensitivewords/import", apis.ImportSensitiveWords)
	group.POST("/apps/sensitivewords/add", apis.AddSensitiveWord)
	group.POST("/apps/sensitivewords/delete", apis.DeleteSensitiveWord)

	return group
}

func CorsHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		method := context.Request.Method
		context.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		context.Writer.Header().Add("Access-Control-Allow-Headers", "*")
		context.Writer.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PATCH, PUT")
		context.Writer.Header().Add("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		context.Writer.Header().Add("Access-Control-Allow-Credentials", "true")

		if method == "OPTIONS" {
			context.AbortWithStatus(http.StatusNoContent)
		}
		context.Next()
	}
}

func InjectCtx() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		appKey := ctx.Request.Header.Get("appkey")
		ctx.Set(string(ctxs.CtxKey_AppKey), appKey)
		ctx.Next()
	}
}

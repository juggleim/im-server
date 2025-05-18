package admingateway

import (
	"embed"
	"fmt"
	"im-server/commons/configures"
	"im-server/commons/gmicro"
	"im-server/services/admingateway/apis"
	"im-server/services/apigateway/services"
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

type AdminGateway struct {
	httpServer *gin.Engine
}

func (ser *AdminGateway) RegisterActors(register gmicro.IActorRegister) {

}

func (ser *AdminGateway) Startup(args map[string]interface{}) {
	ser.httpServer = gin.Default()
	ser.httpServer.HEAD("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, nil)
	})
	loadAdminWeb(ser.httpServer)
	ser.httpServer.Use(CorsHandler(), InjectCtx())
	group := ser.httpServer.Group("/admingateway")
	group.Use(apis.Validate)

	group.POST("/imapiagent", apis.ApiAgent)
	group.GET("/common/address", apis.GetAccessAddress)

	group.POST("/login", apis.Login)
	group.POST("/accounts/updpass", apis.UpdPassword)
	group.POST("/accounts/add", apis.AddAccount)
	group.POST("/accounts/delete", apis.DeleteAccounts)
	group.POST("/accounts/disable", apis.DisableAccounts)
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

	httpPort := configures.Config.AdminGateway.HttpPort
	go ser.httpServer.Run(fmt.Sprintf(":%d", httpPort))
	fmt.Println("Start admingateway with port:", httpPort)
}

func (ser *AdminGateway) Shutdown(force bool) {

}

//go:embed admin/*
var adminFiles embed.FS

func TestFile() {
	files, err := adminFiles.ReadDir("admin/assets")
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, f := range files {
		if !f.IsDir() {
			fmt.Println(f.Name())
		}
	}
}

func loadAdminWeb(httpServer *gin.Engine) {
	files, err := adminFiles.ReadDir("admin/assets")
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, f := range files {
		if !f.IsDir() {
			httpServer.GET("/assets/"+f.Name(), assetsFile)
		}
	}
	httpServer.GET("/dashboard", dashboardPage)
	httpServer.GET("/", dashboardPage)
	httpServer.GET("/:param1", dashboardPage)
	httpServer.GET("/:param1/:param2", dashboardPage)
	httpServer.GET("/:param1/:param2/:param3", dashboardPage)
}

func dashboardPage(ctx *gin.Context) {
	ctx.Writer.Header().Add("Content-Type", "text/html; charset=utf-8")

	var body string
	cacheBody, ok := htmlCache.Load("index.html")
	if ok {
		body = cacheBody.(string)
	} else {
		body = ReadFromFile("admin/index.html")
		htmlCache.Store("index.html", body)
	}
	ctx.String(200, body)
}

var htmlCache sync.Map

func assetsFile(ctx *gin.Context) {
	filePath := ctx.Request.URL.Path
	if strings.HasSuffix(filePath, ".js") {
		ctx.Writer.Header().Add("Content-Type", "application/javascript")
	} else if strings.HasSuffix(filePath, ".css") {
		ctx.Writer.Header().Add("Content-Type", "text/css")
	} else if strings.HasSuffix(filePath, ".png") {
		ctx.Writer.Header().Add("Content-Type", "image/png")
	} else if strings.HasSuffix(filePath, ".ico") {
		ctx.Writer.Header().Add("Content-Type", "image/x-icon")
	}
	var body string
	if cacheBody, ok := htmlCache.Load(filePath); ok {
		body = cacheBody.(string)
	} else {
		body = ReadFromFile("admin" + filePath)
		htmlCache.Store(filePath, body)
	}
	ctx.String(200, body)
}

func ReadFromFile(path string) string {
	// bs, err := os.ReadFile(path)
	bs, err := adminFiles.ReadFile(path)
	if err != nil {
		fmt.Println("read file failed:", err)
		return ""
	}
	return string(bs)
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
		ctx.Set(services.CtxKey_AppKey, appKey)
	}
}

package admingateway

import (
	"embed"
	"fmt"
	"im-server/commons/configures"
	"im-server/commons/gmicro"
	"im-server/services/admingateway/routers"
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

	routers.Route(ser.httpServer, "admingateway")

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

	httpServer.GET("/", dashboardPage)
	httpServer.GET("/dashboard", dashboardPage)
	httpServer.GET("/login", dashboardPage)
	// httpServer.GET("/:param1", dashboardPage)
	// httpServer.GET("/:param1/:param2", dashboardPage)
	// httpServer.GET("/:param1/:param2/:param3", dashboardPage)
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

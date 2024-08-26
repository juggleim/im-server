package navigator

import (
	"fmt"
	"net/http"

	"im-server/commons/configures"
	"im-server/commons/gmicro"
	"im-server/services/navigator/apis"

	"github.com/gin-gonic/gin"
)

type Navigator struct {
	httpServer *gin.Engine
}

func (ser *Navigator) RegisterActors(register gmicro.IActorRegister) {

}
func (ser *Navigator) Startup(args map[string]interface{}) {
	ser.httpServer = gin.Default()
	ser.httpServer.Use(CorsHandler())
	ser.httpServer.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "ok")
	})
	group := ser.httpServer.Group("/navigator")
	group.GET("/general", apis.NaviGet)

	httpPort := configures.Config.NavGateway.HttpPort
	go ser.httpServer.Run(fmt.Sprintf(":%d", httpPort))
	fmt.Println("start Navitor with port :", httpPort)
}

func CorsHandler() gin.HandlerFunc {
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

func (ser *Navigator) Shutdown() {

}

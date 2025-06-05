package apigateway

import (
	"fmt"
	"net/http"

	"im-server/commons/configures"
	"im-server/commons/gmicro"
	"im-server/services/apigateway/routers"

	"github.com/gin-gonic/gin"
)

type ApiGateway struct {
	httpServer *gin.Engine
}

func (ser *ApiGateway) RegisterActors(register gmicro.IActorRegister) {

}

func (ser *ApiGateway) Startup(args map[string]interface{}) {
	if configures.Config.ApiGateway.HttpPort > 0 {
		ser.httpServer = gin.Default()
		ser.httpServer.GET("/", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, "ok-ok")
		})
		ser.httpServer.HEAD("/", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, nil)
		})
		routers.Route(ser.httpServer, "apigateway")

		httpPort := configures.Config.ApiGateway.HttpPort
		go ser.httpServer.Run(fmt.Sprintf(":%d", httpPort))
		fmt.Println("Startup apigateway with port:", httpPort)
	}
}

func (ser *ApiGateway) Shutdown(force bool) {

}

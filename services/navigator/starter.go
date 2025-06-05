package navigator

import (
	"fmt"
	"net/http"

	"im-server/commons/configures"
	"im-server/commons/gmicro"
	"im-server/services/navigator/routers"

	"github.com/gin-gonic/gin"
)

type Navigator struct {
	httpServer *gin.Engine
}

func (ser *Navigator) RegisterActors(register gmicro.IActorRegister) {

}
func (ser *Navigator) Startup(args map[string]interface{}) {
	if configures.Config.NavGateway.HttpPort > 0 {
		ser.httpServer = gin.Default()
		ser.httpServer.GET("/", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, "ok")
		})
		ser.httpServer.HEAD("/", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, nil)
		})
		routers.Route(ser.httpServer, "navigator")

		httpPort := configures.Config.NavGateway.HttpPort
		go ser.httpServer.Run(fmt.Sprintf(":%d", httpPort))
		fmt.Println("Startup navitor with port:", httpPort)
	}
}

func (ser *Navigator) Shutdown(force bool) {

}

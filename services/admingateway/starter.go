package admingateway

import (
	"fmt"
	"im-server/commons/configures"
	"im-server/commons/gmicro"
	"im-server/commons/logs"

	"github.com/gin-gonic/gin"
	consoleConfigures "github.com/juggleim/imserver-console/commons/configures"
	consoleDb "github.com/juggleim/imserver-console/commons/dbcommons"
	consoleLogger "github.com/juggleim/imserver-console/commons/logs"
	consoleRouters "github.com/juggleim/imserver-console/routers"
	jimAdminRouters "github.com/juggleim/jugglechat-server/admins/routers"
)

type AdminGateway struct {
	httpServer *gin.Engine
}

func (ser *AdminGateway) RegisterActors(register gmicro.IActorRegister) {

}

func (ser *AdminGateway) Startup(args map[string]interface{}) {
	if err := consoleConfigures.InitConfigures(); err != nil {
		fmt.Println("Init Console Configures failed", err)
		return
	}
	consoleConfigures.Config.ImApiDomain = fmt.Sprintf("http://127.0.0.1:%d", configures.Config.GetApiPort())
	consoleConfigures.Config.ImAdminDomain = fmt.Sprintf("http://127.0.0.1:%d", configures.Config.AdminGateway.HttpPort)
	consoleConfigures.Config.AdminSecret = configures.Config.AdminSecret
	consoleLogger.SetLogger(logs.GetInfoLogger(), logs.GetErrorLogger())
	if err := consoleDb.InitMysql(); err != nil {
		fmt.Println("Init Console Mysql failed", err)
		return
	}
	consoleDb.Upgrade()

	ser.httpServer = gin.Default()
	group := consoleRouters.Route(ser.httpServer, "admingateway")
	consoleRouters.LoadJuggleChatAdminWeb(ser.httpServer)
	jimAdminRouters.Route(group)

	httpPort := configures.Config.AdminGateway.HttpPort
	go ser.httpServer.Run(fmt.Sprintf(":%d", httpPort))
	fmt.Println("Start admingateway with port:", httpPort)
}

func (ser *AdminGateway) Shutdown(force bool) {

}

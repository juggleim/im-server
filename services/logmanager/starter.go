package logmanager

import (
	"fmt"
	"im-server/commons/bases"
	"im-server/commons/configures"
	"im-server/commons/gmicro"
	"im-server/commons/gmicro/actorsystem"
	"im-server/services/logmanager/actors"
	"im-server/services/logmanager/apis"
	"im-server/services/logmanager/services"

	"github.com/gin-gonic/gin"
)

var serviceName string = "logmanager"

type LogManager struct {
	ginEngine *gin.Engine
}

func (manager *LogManager) RegisterActors(register gmicro.IActorRegister) {
	register.RegisterActor("vlog", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.LogServiceActor{}, serviceName)
	}, 64)
}

func (manager *LogManager) Startup(args map[string]interface{}) {
	if configures.Config.Log.Visual {
		fmt.Println("Startup logmanager.")
		err := services.InitLogDB(fmt.Sprintf("%s/visual_logs", configures.Config.Log.LogPath))
		if err != nil {
			fmt.Printf("Init log db failed. %+v\n", err)
		}
		if configures.Config.Log.VLogHttpPort > 0 {
			manager.startHttp(configures.Config.Log.VLogHttpPort)
		}
	}
}

func (manager *LogManager) Shutdown() {
	fmt.Println("Shutdown logmanager.")
	services.CloseLogDB()
}

func (manager *LogManager) startHttp(httpPort int) {
	engine := gin.Default()
	engine.Use(apis.CorsHandler(), apis.GzipDecompress(), apis.CheckToken)

	group := engine.Group("/api")
	{
		group.POST("/upload-log", apis.UploadClientLog)
		group.POST("/upload-log-plain", apis.UploadClientLogPlain)
	}
	go engine.Run(fmt.Sprintf(":%d", httpPort))

	manager.ginEngine = engine
}

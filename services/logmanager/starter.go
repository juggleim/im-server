package logmanager

import (
	"fmt"
	"im-server/commons/bases"
	"im-server/commons/configures"
	"im-server/commons/gmicro"
	"im-server/commons/gmicro/actorsystem"
	"im-server/services/logmanager/actors"
	"im-server/services/logmanager/services"
)

var serviceName string = "logmanager"

type LogManager struct {
}

func (manager *LogManager) RegisterActors(register gmicro.IActorRegister) {
	register.RegisterActor("vlog", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.LogServiceActor{}, serviceName)
	}, 64)
	register.RegisterActor("qry_vlog", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.QryLogsActor{}, serviceName)
	}, 8)
}

func (manager *LogManager) Startup(args map[string]interface{}) {
	if configures.Config.Log.Visual {
		fmt.Println("Startup logmanager.")
		err := services.InitLogDB(fmt.Sprintf("%s/visual_logs", configures.Config.Log.LogPath))
		if err != nil {
			fmt.Printf("Init log db failed. %+v\n", err)
		}
	}
}

func (manager *LogManager) Shutdown() {
	fmt.Println("Shutdown logmanager.")
	services.CloseLogDB()
}

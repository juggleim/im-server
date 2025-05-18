package logmanager

import (
	"fmt"
	"im-server/commons/bases"
	"im-server/commons/gmicro"
	"im-server/commons/gmicro/actorsystem"
	"im-server/services/logmanager/actors"
)

var serviceName string = "logmanager"

type LogManager struct {
}

func (manager *LogManager) RegisterActors(register gmicro.IActorRegister) {
	register.RegisterActor("qry_vlog", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.QryLogsActor{}, serviceName)
	})
}

func (manager *LogManager) Startup(args map[string]interface{}) {
	fmt.Println("Startup logmanager.")
}

func (manager *LogManager) Shutdown(force bool) {
	fmt.Println("Shutdown logmanager.")
}

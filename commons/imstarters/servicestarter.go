package imstarters

import (
	"im-server/commons/bases"
	"im-server/commons/gmicro"
)

type IServiceStarter interface{}

type IRegisterActorsHandler interface {
	RegisterActors(register gmicro.IActorRegister)
}

type IStartupHandler interface {
	Startup(args map[string]interface{})
}

type IShutdownHandler interface {
	Shutdown(force bool)
}

var serverList []IServiceStarter

func Loaded(server IServiceStarter) {
	if server != nil {
		//register actors
		if regActorsHandler, ok := server.(IRegisterActorsHandler); ok {
			regActorsHandler.RegisterActors(bases.GetCluster())
		}
		serverList = append(serverList, server)
	}
}

func Startup() {
	for _, server := range serverList {
		//execute startup
		if startHandler, ok := server.(IStartupHandler); ok {
			startHandler.Startup(map[string]interface{}{})
		}
	}
	bases.Startup()
}

func Shutdown(force bool) {
	//remove self from zk TODO
	for _, server := range serverList {
		//execute startup
		if shutdownHandler, ok := server.(IShutdownHandler); ok {
			shutdownHandler.Shutdown(force)
		}
	}
}

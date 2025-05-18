package userstatussub

import (
	"fmt"
	"im-server/commons/bases"
	"im-server/commons/gmicro"
	"im-server/commons/gmicro/actorsystem"
	"im-server/services/userstatussub/actors"
)

var serviceName string = "userstatussub"

type UserStatusSubManager struct{}

func (manager *UserStatusSubManager) RegisterActors(register gmicro.IActorRegister) {
	register.RegisterActor("pub_status", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.PublishStatusActor{}, serviceName)
	})
	register.RegisterActor("sub_users", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.SubscribeActor{}, serviceName)
	})
	register.RegisterActor("inner_sub_users", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.InnerSubscribeActor{}, serviceName)
	})
	register.RegisterActor("unsub_users", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.UnSubscribeActor{}, serviceName)
	})
	register.RegisterActor("inner_unsub_users", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.InnerUnSubscribeActor{}, serviceName)
	})
}

func (manager *UserStatusSubManager) Startup(args map[string]interface{}) {
	fmt.Println("Startup userstatussub.")
}

func (manager *UserStatusSubManager) Shutdown(force bool) {
	fmt.Println("Shutdown userstatussub.")
}

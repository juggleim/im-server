package statussubscriptions

import (
	"fmt"
	"im-server/commons/bases"
	"im-server/commons/gmicro"
	"im-server/commons/gmicro/actorsystem"
	"im-server/services/statussubscriptions/actors"
)

var serviceName string = "statussubscriptions"

type StatusSubscriptionsManager struct{}

func (manager *StatusSubscriptionsManager) RegisterActors(register gmicro.IActorRegister) {
	register.RegisterActor("sub_users", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.SubUsersActor{}, serviceName)
	})
	register.RegisterActor("qry_user_status", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.QryUserStatusActor{}, serviceName)
	})
	register.RegisterActor("unsub_users", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.UnSubUsersActor{}, serviceName)
	})
	register.RegisterActor("pub_user_status", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.PubUserStatusActor{}, serviceName)
	})
	register.RegisterActor("sync_sub_change", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.SyncSubChangeActor{}, serviceName)
	})
}

func (manager *StatusSubscriptionsManager) Startup(args map[string]interface{}) {
	fmt.Println("Startup", serviceName)
}

func (manager *StatusSubscriptionsManager) Shutdown(force bool) {
	fmt.Println("Shutdown", serviceName)
}

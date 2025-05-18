package subscriptions

import (
	"fmt"
	"im-server/commons/bases"
	"im-server/commons/gmicro"
	"im-server/commons/gmicro/actorsystem"
	"im-server/services/subscriptions/actors"
)

var serviceName string = "subscriptions"

type SubscriptionManager struct{}

func (manager *SubscriptionManager) RegisterActors(register gmicro.IActorRegister) {
	register.RegisterActor("msg_sub", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.MsgSubActor{}, serviceName)
	})
	register.RegisterActor("online_offline_sub", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.OnlineOfflineSubActor{}, serviceName)
	})
}

func (manager *SubscriptionManager) Startup(args map[string]interface{}) {
	fmt.Println("Startup subscriptions.")
}
func (manager *SubscriptionManager) Shutdown(force bool) {
	fmt.Println("Shutdown subscriptions.")
}

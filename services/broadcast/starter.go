package broadcast

import (
	"fmt"
	"im-server/commons/bases"
	"im-server/commons/gmicro"
	"im-server/commons/gmicro/actorsystem"
	"im-server/services/broadcast/actors"
)

var serviceName string = "broadcast"

type BroadcastManager struct{}

func (manager *BroadcastManager) RegisterActors(register gmicro.IActorRegister) {
	register.RegisterActor("bc_msg", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.BroadcastMsgActor{}, serviceName)
	}, 8)
	register.RegisterActor("gc_msg", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.GroupCastMsgActor{}, serviceName)
	}, 16)
}

func (manager *BroadcastManager) Startup(args map[string]interface{}) {
	fmt.Println("Startup broadcast.")
}
func (manager *BroadcastManager) Shutdown() {
	fmt.Println("Shutdown broadcast.")
}

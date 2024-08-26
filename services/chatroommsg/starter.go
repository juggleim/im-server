package chatroommsg

import (
	"fmt"
	"im-server/commons/bases"
	"im-server/commons/gmicro"
	"im-server/commons/gmicro/actorsystem"
	"im-server/services/chatroommsg/actors"
)

type ChatroomMsgManager struct{}

var serviceName string = "chatroommsg"

func (manager *ChatroomMsgManager) RegisterActors(register gmicro.IActorRegister) {
	register.RegisterActor("c_members_dispatch", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.DispatchMembersActor{}, serviceName)
	}, 32)
	register.RegisterActor("c_msgs_dispatch", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.DispatchMsgsActor{}, serviceName)
	}, 64)
	register.RegisterActor("c_atts_dispatch", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.DispatchAttsActor{}, serviceName)
	}, 32)
	register.RegisterActor("c_sync_msgs", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.SyncChatroomMsgsActor{}, serviceName)
	}, 64)
	register.RegisterActor("c_sync_atts", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.SyncChatroomAttsActor{}, serviceName)
	}, 64)
	register.RegisterActor("c_chrm_dispatch", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.DispatchChrmActor{}, serviceName)
	}, 8)
}

func (manager *ChatroomMsgManager) Startup(args map[string]interface{}) {
	fmt.Println("Startup chatroommsg.")
}
func (manager *ChatroomMsgManager) Shutdown() {
	fmt.Println("Shutdown chatroommsg.")
}

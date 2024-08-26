package chatroom

import (
	"fmt"

	"im-server/commons/bases"
	"im-server/commons/gmicro"
	"im-server/commons/gmicro/actorsystem"
	"im-server/services/chatroom/actors"
)

type ChatroomManager struct{}

var serviceName string = "chatroom"

func (manager *ChatroomManager) RegisterActors(register gmicro.IActorRegister) {
	register.RegisterActor("c_msg", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.ChatMsgActor{}, serviceName)
	}, 64)
	register.RegisterActor("c_add_att", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.AddAttActor{}, serviceName)
	}, 32)
	register.RegisterActor("c_batch_add_att", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.BatchAddAttActor{}, serviceName)
	}, 32)
	register.RegisterActor("c_del_att", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.DelAttActor{}, serviceName)
	}, 32)
	register.RegisterActor("c_batch_del_att", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.BatchDelAttActor{}, serviceName)
	}, 32)
	register.RegisterActor("c_qry_atts", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.QryAttsActor{}, serviceName)
	}, 32)
	register.RegisterActor("c_join", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.JoinChatroomActor{}, serviceName)
	}, 32)
	register.RegisterActor("c_quit", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.QuitChatroomActor{}, serviceName)
	}, 32)
	register.RegisterActor("c_sync_quit", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.SyncQuitActor{}, serviceName)
	}, 32)
	register.RegisterActor("c_create", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.CreateChrmActor{}, serviceName)
	}, 16)
	register.RegisterActor("c_destroy", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.DestroyChrmActor{}, serviceName)
	}, 16)
	register.RegisterActor("c_sync_partial", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.SyncPartialInfoActor{}, serviceName)
	}, 16)
	register.RegisterActor("c_qry_chrm", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.QryChatroomActor{}, serviceName)
	}, 16)
	register.RegisterActor("c_ban_user", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.BanUserActor{}, serviceName)
	}, 16)
	register.RegisterActor("c_qry_ban_user", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.QryBanUsersActor{}, serviceName)
	}, 8)
	register.RegisterActor("chrm_mute", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.ChrmMuteActor{}, serviceName)
	}, 8)
}

func (manager *ChatroomManager) Startup(args map[string]interface{}) {
	fmt.Println("Startup chatroom.")
}
func (manager *ChatroomManager) Shutdown() {
	fmt.Println("Shutdown chatroom.")
}

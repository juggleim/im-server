package message

import (
	"fmt"

	"im-server/commons/bases"
	"im-server/commons/gmicro"
	"im-server/commons/gmicro/actorsystem"
	"im-server/services/message/actors"
)

var serviceName string = "message"

type MessageManager struct{}

func (manager *MessageManager) RegisterActors(register gmicro.IActorRegister) {
	register.RegisterActor("p_msg", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.PrivateMsgActor{}, serviceName)
	}, 6144)
	register.RegisterActor("s_msg", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.SystemMsgActor{}, serviceName)
	}, 64)
	register.RegisterActor("msg_dispatch", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.MsgDispatchActor{}, serviceName)
	}, 6144)
	register.RegisterActor("sendbox", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.SendBoxActor{}, serviceName)
	}, 32)
	register.RegisterActor("pri_stream", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.PrivateStreamActor{}, serviceName)
	}, 32)
	register.RegisterActor("sync_msgs", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.SyncMsgActor{}, serviceName)
	}, 64)
	register.RegisterActor("push_switch", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.PushSwitchActor{}, serviceName)
	}, 16)
	register.RegisterActor("brd_inbox", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.BrdcastInboxActor{}, serviceName)
	}, 8)
	register.RegisterActor("brd_append", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.BrdAppendActor{}, serviceName)
	}, 8)
	register.RegisterActor("block_users", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.BlockUsersActor{}, serviceName)
	}, 8)
	register.RegisterActor("qry_block_users", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.QryBlockUsersActor{}, serviceName)
	}, 8)
	register.RegisterActor("del_conver_cache", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.DelConverCacheActor{}, serviceName)
	}, 8)
}

func (manager *MessageManager) Startup(args map[string]interface{}) {
	fmt.Println("Startup message.")
}
func (manager *MessageManager) Shutdown() {
	fmt.Println("Shutdown message.")
}

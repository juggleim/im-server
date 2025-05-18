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
	register.RegisterStandaloneActor("p_msg", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.PrivateMsgActor{}, serviceName)
	}, 3072)
	register.RegisterActor("s_msg", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.SystemMsgActor{}, serviceName)
	})
	register.RegisterStandaloneActor("msg_dispatch", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.MsgDispatchActor{}, serviceName)
	}, 6144)
	register.RegisterActor("sendbox", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.SendBoxActor{}, serviceName)
	})
	register.RegisterActor("sync_msgs", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.SyncMsgActor{}, serviceName)
	})
	register.RegisterActor("push_switch", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.PushSwitchActor{}, serviceName)
	})
	register.RegisterActor("brd_inbox", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.BrdcastInboxActor{}, serviceName)
	})
	register.RegisterActor("brd_append", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.BrdAppendActor{}, serviceName)
	})
	register.RegisterActor("block_users", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.BlockUsersActor{}, serviceName)
	})
	register.RegisterActor("qry_block_users", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.QryBlockUsersActor{}, serviceName)
	})
	register.RegisterActor("check_block_user", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.CheckBlockUserActor{}, serviceName)
	})
	register.RegisterActor("del_conver_cache", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.DelConverCacheActor{}, serviceName)
	})
	register.RegisterActor("msg_ack", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.MsgAckActor{}, serviceName)
	})
	register.RegisterActor("imp_pri_msg", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.ImportPrivateHisMsgActor{}, serviceName)
	})
	register.RegisterActor("upd_push_status", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.UpdPushStatusActor{}, serviceName)
	})
	register.RegisterActor("user_push", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.UserPushActor{}, serviceName)
	})
}

func (manager *MessageManager) Startup(args map[string]interface{}) {
	fmt.Println("Startup message.")
}
func (manager *MessageManager) Shutdown(force bool) {
	fmt.Println("Shutdown message.")
}

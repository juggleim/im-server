package conversation

import (
	"fmt"
	"im-server/commons/bases"
	"im-server/commons/gmicro"
	"im-server/commons/gmicro/actorsystem"
	"im-server/services/conversation/actors"
)

var serviceName string = "conversation"

type ConversationManager struct{}

func (manager *ConversationManager) RegisterActors(register gmicro.IActorRegister) {
	register.RegisterActor("sync_convers", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.SyncConversationsActor{}, serviceName)
	})
	register.RegisterActor("qry_convers", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.QryConversationsActor{}, serviceName)
	})
	register.RegisterActor("qry_conver", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.QryConversationActor{}, serviceName)
	})
	register.RegisterActor("qry_global_convers", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.QryGlobalConversActor{}, serviceName)
	})
	register.RegisterActor("qry_total_unread_count", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.QryTotalUnreadCountActor{}, serviceName)
	})
	register.RegisterActor("clear_unread", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.ClearUnReadActor{}, serviceName)
	})
	register.RegisterActor("mark_unread", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.MarkUnreadActor{}, serviceName)
	})
	register.RegisterActor("clear_total_unread", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.ClearTotalUnreadActor{}, serviceName)
	})
	register.RegisterActor("del_convers", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.DelConversationsActor{}, serviceName)
	})
	register.RegisterActor("top_convers", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.TopConversActor{}, serviceName)
	})
	register.RegisterActor("qry_top_convers", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.QryTopConversActor{}, serviceName)
	})
	register.RegisterActor("add_conver", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.AddConversationActor{}, serviceName)
	})
	// register.RegisterActor("qry_undisturb", func() actorsystem.IUntypedActor {
	// 	return bases.BaseProcessActor(&actors.QryUndisturbActor{})
	// }, 64)
	register.RegisterActor("undisturb_convers", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.UndisturbConversActor{}, serviceName)
	})
	register.RegisterActor("upd_latest_msg", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.UpdLatestMsgActor{}, serviceName)
	})

	register.RegisterActor("qry_mention_msgs", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.QryMentionMsgsActor{}, serviceName)
	})
}

func (manager *ConversationManager) Startup(args map[string]interface{}) {
	fmt.Println("Startup conversation.")
}
func (manager *ConversationManager) Shutdown() {
	fmt.Println("Shutdown conversation.")
}

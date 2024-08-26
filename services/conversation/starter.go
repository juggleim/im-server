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
	}, 64)
	register.RegisterActor("qry_convers", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.QryConversationsActor{}, serviceName)
	}, 64)
	register.RegisterActor("qry_conver", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.QryConversationActor{}, serviceName)
	}, 64)
	register.RegisterActor("qry_global_convers", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.QryGlobalConversActor{}, serviceName)
	}, 8)
	register.RegisterActor("qry_total_unread_count", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.QryTotalUnreadCountActor{}, serviceName)
	}, 64)
	register.RegisterActor("clear_unread", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.ClearUnReadActor{}, serviceName)
	}, 32)
	register.RegisterActor("mark_unread", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.MarkUnreadActor{}, serviceName)
	}, 16)
	register.RegisterActor("clear_total_unread", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.ClearTotalUnreadActor{}, serviceName)
	}, 16)
	register.RegisterActor("del_convers", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.DelConversationsActor{}, serviceName)
	}, 16)
	register.RegisterActor("top_convers", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.TopConversActor{}, serviceName)
	}, 16)
	register.RegisterActor("qry_top_convers", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.QryTopConversActor{}, serviceName)
	}, 64)
	register.RegisterActor("add_conver", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.AddConversationActor{}, serviceName)
	}, 64)
	// register.RegisterActor("qry_undisturb", func() actorsystem.IUntypedActor {
	// 	return bases.BaseProcessActor(&actors.QryUndisturbActor{})
	// }, 64)
	register.RegisterActor("undisturb_convers", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.UndisturbConversActor{}, serviceName)
	}, 8)
	register.RegisterActor("upd_latest_msg", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.UpdLatestMsgActor{}, serviceName)
	}, 8)

	register.RegisterActor("qry_mention_msgs", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.QryMentionMsgsActor{}, serviceName)
	}, 32)
}

func (manager *ConversationManager) Startup(args map[string]interface{}) {
	fmt.Println("Startup conversation.")
}
func (manager *ConversationManager) Shutdown() {
	fmt.Println("Shutdown conversation.")
}

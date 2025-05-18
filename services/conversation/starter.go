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
	register.RegisterActor("add_conver", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.AddConversationActor{}, serviceName)
	})
	register.RegisterActor("batch_add_conver", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.BatchAddConversationActor{}, serviceName)
	})
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
	register.RegisterActor("tag_add_convers", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.TagAddConversActor{}, serviceName)
	})
	register.RegisterActor("tag_del_convers", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.TagDelConversActor{}, serviceName)
	})
	register.RegisterActor("qry_user_conver_tags", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.QryUserConverTagsActor{}, serviceName)
	})
	register.RegisterActor("del_user_conver_tags", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.DelUserConverTagsActor{}, serviceName)
	})
}

func (manager *ConversationManager) Startup(args map[string]interface{}) {
	fmt.Println("Startup conversation.")
}
func (manager *ConversationManager) Shutdown(force bool) {
	fmt.Println("Shutdown conversation.")
}

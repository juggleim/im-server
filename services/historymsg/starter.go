package historymsg

import (
	"fmt"
	"im-server/commons/bases"
	"im-server/commons/gmicro"
	"im-server/commons/gmicro/actorsystem"
	"im-server/services/historymsg/actors"
)

var serviceName string = "historymsg"

type HistoryMsgManager struct{}

func (manager *HistoryMsgManager) RegisterActors(register gmicro.IActorRegister) {
	register.RegisterActor("add_hismsg", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.AddHisMsgActor{}, serviceName)
	})
	register.RegisterMultiMethodActor([]string{"del_msg", "del_hismsg"}, func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.DelMsgActor{}, serviceName)
	})
	register.RegisterActor("upd_stream", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.UpdStreamActor{}, serviceName)
	})
	register.RegisterActor("qry_latest_hismsg", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.QryLatestMsgActor{}, serviceName)
	})
	register.RegisterActor("qry_hismsgs", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.QryHistoryMsgsActor{}, serviceName)
	})
	register.RegisterActor("qry_first_unread_msg", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.QryFirstUnreadMsgActor{}, serviceName)
	})
	register.RegisterActor("qry_hismsg_by_ids", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.QryHistoryMsgByIdsActor{}, serviceName)
	})
	register.RegisterActor("clean_hismsg", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.CleanHisMsgActor{}, serviceName)
	})
	register.RegisterActor("recall_msg", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.RecallMsgActor{}, serviceName)
	})
	register.RegisterActor("modify_msg", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.ModifyMsgActor{}, serviceName)
	})
	register.RegisterActor("mark_read", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.MarkReadActor{}, serviceName)
	})
	register.RegisterActor("mark_grp_msg_read", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.MarkGrpMsgReadActor{}, serviceName)
	})
	register.RegisterActor("qry_read_infos", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.QryReadInfosActor{}, serviceName)
	})
	register.RegisterActor("qry_read_detail", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.QryReadDetailActor{}, serviceName)
	})
	register.RegisterActor("merge_msgs", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.MergeMsgActor{}, serviceName)
	})
	register.RegisterActor("qry_merged_msgs", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.QryMergedMsgsActor{}, serviceName)
	})
	register.RegisterActor("msg_ext", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.SetMsgExtActor{}, serviceName)
	})
	register.RegisterActor("del_msg_ext", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.DelMsgExtActor{}, serviceName)
	})
	register.RegisterActor("msg_exset", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.AddMsgExSetActor{}, serviceName)
	})
	register.RegisterActor("del_msg_exset", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.DelMsgExSetActor{}, serviceName)
	})
	register.RegisterActor("batch_trans", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.MultiTransActor{}, serviceName)
	})
}

func (manager *HistoryMsgManager) Startup(args map[string]interface{}) {
	fmt.Println("Startup historymsg.")
}
func (manager *HistoryMsgManager) Shutdown() {
	fmt.Println("Shutdown historymsg.")
}

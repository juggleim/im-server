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
	}, 64)
	register.RegisterActor("del_msg", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.DelMsgActor{}, serviceName)
	}, 32)
	register.RegisterActor("qry_latest_hismsg", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.QryLatestMsgActor{}, serviceName)
	}, 32)
	register.RegisterActor("qry_hismsgs", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.QryHistoryMsgsActor{}, serviceName)
	}, 64)
	register.RegisterActor("qry_first_unread_msg", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.QryFirstUnreadMsgActor{}, serviceName)
	}, 32)
	register.RegisterActor("qry_hismsg_by_ids", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.QryHistoryMsgByIdsActor{}, serviceName)
	}, 32)
	register.RegisterActor("clean_hismsg", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.CleanHisMsgActor{}, serviceName)
	}, 32)
	register.RegisterActor("recall_msg", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.RecallMsgActor{}, serviceName)
	}, 32)
	register.RegisterActor("modify_msg", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.ModifyMsgActor{}, serviceName)
	}, 32)
	register.RegisterActor("mark_read", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.MarkReadActor{}, serviceName)
	}, 64)
	register.RegisterActor("mark_grp_msg_read", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.MarkGrpMsgReadActor{}, serviceName)
	}, 64)
	register.RegisterActor("qry_read_infos", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.QryReadInfosActor{}, serviceName)
	}, 32)
	register.RegisterActor("qry_read_detail", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.QryReadDetailActor{}, serviceName)
	}, 16)
	register.RegisterActor("merge_msgs", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.MergeMsgActor{}, serviceName)
	}, 16)
	register.RegisterActor("qry_merged_msgs", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.QryMergedMsgsActor{}, serviceName)
	}, 32)
	register.RegisterActor("msg_ext", func() actorsystem.IUntypedActor { //TODO
		return bases.BaseProcessActor(&actors.SetMsgExtActor{}, serviceName)
	}, 32)
}

func (manager *HistoryMsgManager) Startup(args map[string]interface{}) {
	fmt.Println("Startup historymsg.")
}
func (manager *HistoryMsgManager) Shutdown() {
	fmt.Println("Shutdown historymsg.")
}

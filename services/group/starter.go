package group

import (
	"fmt"
	"im-server/commons/bases"
	"im-server/commons/gmicro"
	"im-server/commons/gmicro/actorsystem"
	"im-server/services/group/actors"
)

var serviceName string = "group"

type GroupManager struct{}

func (manager *GroupManager) RegisterActors(register gmicro.IActorRegister) {
	register.RegisterActor("g_msg", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.GroupMsgActor{}, serviceName)
	})
	register.RegisterActor("g_add_members", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.AddMemberActor{}, serviceName)
	})
	register.RegisterActor("g_del_members", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.DelMemberActor{}, serviceName)
	})
	register.RegisterActor("g_dissolve", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.DissolveGroupActor{}, serviceName)
	})
	register.RegisterActor("g_qry_members", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.QryGroupMembersActor{}, serviceName)
	})
	register.RegisterActor("g_check_members", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.CheckGroupMemberActor{}, serviceName)
	})
	register.RegisterActor("group_mute", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.GroupMuteActor{}, serviceName)
	})
	register.RegisterActor("qry_group_info", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.QryGroupInfoActor{}, serviceName)
	})
	register.RegisterActor("qry_group_info_by_ids", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.QryGroupInfoByIdsActor{}, serviceName)
	})
	register.RegisterActor("upd_group_info", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.UpdGroupInfoActor{}, serviceName)
	})
	register.RegisterActor("group_member_mute", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.GroupMemberMuteActor{}, serviceName)
	})
	register.RegisterActor("group_member_allow", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.GroupMemberAllowActor{}, serviceName)
	})
	register.RegisterActor("qry_group_members_by_ids", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.QryGroupMembersByIdsActor{}, serviceName)
	})
	register.RegisterActor("qry_group_snapshot", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.QryGrpSnapshotActor{}, serviceName)
	})
	register.RegisterActor("upd_grp_conver", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.UpdGrpConverActor{}, serviceName)
	})
	//settings
	register.RegisterActor("qry_grp_member_settings", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.QryMemberSettingsActor{}, serviceName)
	})
	register.RegisterActor("set_grp_member_setting", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.SetGrpMemberSettingActor{}, serviceName)
	})
	register.RegisterActor("imp_grp_msg", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.ImportGroupHisMsgActor{}, serviceName)
	})
}

func (manager *GroupManager) Startup(args map[string]interface{}) {
	fmt.Println("Startup group.")
}
func (manager *GroupManager) Shutdown(force bool) {
	fmt.Println("Shutdown group.")
}

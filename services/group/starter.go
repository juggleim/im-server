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
	}, 64)
	register.RegisterActor("g_add_members", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.AddMemberActor{}, serviceName)
	}, 16)
	register.RegisterActor("g_del_members", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.DelMemberActor{}, serviceName)
	}, 8)
	register.RegisterActor("g_dissolve", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.DissolveGroupActor{}, serviceName)
	}, 8)
	register.RegisterActor("g_qry_members", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.QryGroupMembersActor{}, serviceName)
	}, 8)
	register.RegisterActor("g_check_members", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.CheckGroupMemberActor{}, serviceName)
	}, 16)
	register.RegisterActor("group_mute", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.GroupMuteActor{}, serviceName)
	}, 8)
	register.RegisterActor("qry_group_info", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.QryGroupInfoActor{}, serviceName)
	}, 8)
	register.RegisterActor("qry_group_info_by_ids", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.QryGroupInfoByIdsActor{}, serviceName)
	}, 8)
	register.RegisterActor("upd_group_info", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.UpdGroupInfoActor{}, serviceName)
	}, 8)
	register.RegisterActor("group_member_mute", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.GroupMemberMuteActor{}, serviceName)
	}, 8)
	register.RegisterActor("group_member_allow", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.GroupMemberAllowActor{}, serviceName)
	}, 8)
	register.RegisterActor("qry_group_members_by_ids", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.QryGroupMembersByIdsActor{}, serviceName)
	}, 8)
	register.RegisterActor("qry_group_snapshot", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.QryGrpSnapshotActor{}, serviceName)
	}, 8)
	register.RegisterActor("upd_grp_conver", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.UpdGrpConverActor{}, serviceName)
	}, 8)
	//settings
	register.RegisterActor("set_grp_setting", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.SetGrpSettingActor{}, serviceName)
	}, 8)
	register.RegisterActor("qry_grp_member_settings", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.QryMemberSettingsActor{}, serviceName)
	}, 16)
}

func (manager *GroupManager) Startup(args map[string]interface{}) {
	fmt.Println("Startup group.")
}
func (manager *GroupManager) Shutdown() {
	fmt.Println("Shutdown group.")
}

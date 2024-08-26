package usermanager

import (
	"fmt"

	"im-server/commons/bases"
	"im-server/commons/gmicro"
	"im-server/commons/gmicro/actorsystem"
	"im-server/services/usermanager/actors"
)

var serviceName string = "usermanager"

type UserManager struct{}

func (manager *UserManager) RegisterActors(register gmicro.IActorRegister) {
	register.RegisterActor("upstream", func() actorsystem.IUntypedActor {
		return &actors.UpstreamActor{}
	}, 64)
	register.RegisterActor("reg_user", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.UserRegistActor{}, serviceName)
	}, 64)
	register.RegisterActor("add_bot", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.AddBotActor{}, serviceName)
	}, 16)
	register.RegisterActor("qry_user_info_by_ids", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.QryUserInfoByIdsActor{}, serviceName)
	}, 32)
	register.RegisterActor("qry_user_info", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.QryUserInfoActor{}, serviceName)
	}, 32)
	register.RegisterActor("upd_user_info", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.UpdUserInfoActor{}, serviceName)
	}, 8)
	register.RegisterActor("set_user_settings", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.SetUserSettingActor{}, serviceName)
	}, 16)
	register.RegisterActor("set_user_undisturb", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.SetUserUndisturbActor{}, serviceName)
	}, 8)
	register.RegisterActor("get_user_undisturb", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.GetUserUndisturbActor{}, serviceName)
	}, 8)
	register.RegisterActor("set_user_status", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.SetUserStatusActor{}, serviceName)
	}, 16)
	register.RegisterActor("qry_user_status", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.QryUserStatusActor{}, serviceName)
	}, 32)
	register.RegisterActor("inner_qry_user_status", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.InnerQryUserStatusActor{}, serviceName)
	}, 32)
}

func (manager *UserManager) Startup(args map[string]interface{}) {
	fmt.Println("Startup usermanager.")
}
func (manager *UserManager) Shutdown() {
	fmt.Println("Shutdown usermanager.")
}

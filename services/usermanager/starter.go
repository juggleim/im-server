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
	register.RegisterStandaloneActor("upstream", func() actorsystem.IUntypedActor {
		return &actors.UpstreamActor{}
	}, 6144)
	register.RegisterActor("send_stream", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.SendStreamActor{}, serviceName)
	})
	register.RegisterActor("reg_user", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.UserRegistActor{}, serviceName)
	})
	register.RegisterActor("add_bot", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.AddBotActor{}, serviceName)
	})
	register.RegisterActor("qry_user_info_by_ids", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.QryUserInfoByIdsActor{}, serviceName)
	})
	register.RegisterActor("qry_user_info", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.QryUserInfoActor{}, serviceName)
	})
	register.RegisterActor("upd_user_info", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.UpdUserInfoActor{}, serviceName)
	})
	register.RegisterActor("set_user_settings", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.SetUserSettingActor{}, serviceName)
	})
	register.RegisterActor("get_user_settings", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.GetUserSettingActor{}, serviceName)
	})
	register.RegisterActor("set_user_undisturb", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.SetUserUndisturbActor{}, serviceName)
	})
	register.RegisterActor("get_user_undisturb", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.GetUserUndisturbActor{}, serviceName)
	})
	register.RegisterActor("set_user_status", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.SetUserStatusActor{}, serviceName)
	})
	register.RegisterActor("qry_user_status", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.QryUserStatusActor{}, serviceName)
	})
	register.RegisterActor("inner_qry_user_status", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.InnerQryUserStatusActor{}, serviceName)
	})
	register.RegisterActor("pri_global_mute", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.PriGlobalMuteActor{}, serviceName)
	})
	register.RegisterActor("qry_pri_global_mute", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.QryPriGlobalMuteActor{}, serviceName)
	})
}

func (manager *UserManager) Startup(args map[string]interface{}) {
	fmt.Println("Startup usermanager.")
}
func (manager *UserManager) Shutdown(force bool) {
	fmt.Println("Shutdown usermanager.")
}

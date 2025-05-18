package push

import (
	"fmt"
	"im-server/commons/bases"
	"im-server/commons/gmicro"
	"im-server/commons/gmicro/actorsystem"
	"im-server/services/pushmanager/actors"
)

var serviceName string = "pushmanager"

type PushManager struct{}

func (manager *PushManager) RegisterActors(register gmicro.IActorRegister) {
	register.RegisterActor("push", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.PushActor{}, serviceName)
	})
	register.RegisterActor("reg_push_token", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.RegPushTokenActor{}, serviceName)
	})
	register.RegisterActor("remove_push_token", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.RemovePushTokenActor{}, serviceName)
	})
	register.RegisterActor("qry_user_tags", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.QryUserTagsActor{}, serviceName)
	})
	register.RegisterActor("add_user_tags", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.AddUserTagsActor{}, serviceName)
	})
	register.RegisterActor("del_user_tags", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.DelUserTagsActor{}, serviceName)
	})
	register.RegisterActor("clear_user_tags", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.ClearUserTagsActor{}, serviceName)
	})
	register.RegisterActor("push_with_tags", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.PushWithTagsActor{}, serviceName)
	})
}

func (manager *PushManager) Startup(args map[string]interface{}) {
	fmt.Println("Startup pushmanager.")
}
func (manager *PushManager) Shutdown(force bool) {
	fmt.Println("Shutdown pushmanager.")
}

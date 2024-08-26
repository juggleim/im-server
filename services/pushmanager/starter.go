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
	}, 64)
	register.RegisterActor("reg_push_token", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.RegPushTokenActor{}, serviceName)
	}, 64)
	register.RegisterActor("qry_user_tags", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.QryUserTagsActor{}, serviceName)
	}, 64)
	register.RegisterActor("add_user_tags", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.AddUserTagsActor{}, serviceName)
	}, 64)
	register.RegisterActor("del_user_tags", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.DelUserTagsActor{}, serviceName)
	}, 64)
	register.RegisterActor("clear_user_tags", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.ClearUserTagsActor{}, serviceName)
	}, 64)
	register.RegisterActor("push_with_tags", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.PushWithTagsActor{}, serviceName)
	}, 64)
}

func (manager *PushManager) Startup(args map[string]interface{}) {
	fmt.Println("Startup pushmanager.")
}
func (manager *PushManager) Shutdown() {
	fmt.Println("Shutdown pushmanager.")
}

package appbusiness

import (
	"fmt"
	"im-server/commons/gmicro"
)

type AppBusiness struct{}

// var serviceName string = "appbusiness"

func (bus *AppBusiness) RegisterActors(register gmicro.IActorRegister) {
	// register.RegisterActor("app_upd_user", func() actorsystem.IUntypedActor {
	// 	return bases.BaseProcessActor(&users.UserUpdateActor{}, serviceName)
	// })
}

func (bus *AppBusiness) Startup(args map[string]interface{}) {
	fmt.Println("Startup appbusiness.")
}

func (bus *AppBusiness) Shutdown() {
	fmt.Println("Shutdown appbusiness.")
}

package appbusiness

import (
	"fmt"
	"im-server/commons/gmicro"

	"github.com/juggleim/jugglechat-server/configures"
	"github.com/juggleim/jugglechat-server/log"
	"github.com/juggleim/jugglechat-server/storages/dbs/dbcommons"
)

type AppBusiness struct{}

// var serviceName string = "appbusiness"

func (bus *AppBusiness) RegisterActors(register gmicro.IActorRegister) {
	// register.RegisterActor("app_upd_user", func() actorsystem.IUntypedActor {
	// 	return bases.BaseProcessActor(&users.UserUpdateActor{}, serviceName)
	// })
}

func (bus *AppBusiness) Startup(args map[string]interface{}) {
	//init configure
	if err := configures.InitConfigures(); err != nil {
		fmt.Println("Init Configures failed", err)
		return
	}
	//init log
	log.InitLogs()
	//init mysql
	if err := dbcommons.InitMysql(); err != nil {
		log.Error("Init Mysql failed.", err)
		return
	}

	fmt.Println("Startup appbusiness.")
}

func (bus *AppBusiness) Shutdown(force bool) {
	fmt.Println("Shutdown appbusiness.")
}

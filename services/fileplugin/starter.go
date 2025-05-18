package fileplugin

import (
	"fmt"
	"im-server/commons/bases"
	"im-server/commons/gmicro"
	"im-server/commons/gmicro/actorsystem"
	"im-server/services/fileplugin/actors"
)

var serviceName string = "fileplugin"

type FilePlugin struct{}

func (manager *FilePlugin) RegisterActors(register gmicro.IActorRegister) {
	register.RegisterActor("file_cred", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.QryUploadTokenActor{}, serviceName)
	})
	register.RegisterActor("upd_log_state", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.ReportClientLogStateActor{}, serviceName)
	})
}

func (manager *FilePlugin) Startup(args map[string]interface{}) {
	fmt.Println("Startup fileplugin.")
}
func (manager *FilePlugin) Shutdown(force bool) {
	fmt.Println("Shutdown fileplugin.")
}

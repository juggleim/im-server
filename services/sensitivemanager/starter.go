package sensitivemanager

import (
	"fmt"
	"im-server/commons/bases"
	"im-server/commons/gmicro"
	"im-server/commons/gmicro/actorsystem"
	"im-server/services/sensitivemanager/actors"
)

var serviceName string = "sensitivemanager"

type SensitiveManager struct {
}

func (manager *SensitiveManager) RegisterActors(register gmicro.IActorRegister) {
	register.RegisterActor("add_sensitive_words", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.AddWordsActor{}, serviceName)
	})
	register.RegisterActor("del_sensitive_words", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.DelWordsActor{}, serviceName)
	})
	register.RegisterActor("sensitive_filter_text", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.FilterTextActor{}, serviceName)
	})
	register.RegisterActor("qry_sensitive_words", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.QrySensitiveWordsActor{}, serviceName)
	})
}

func (manager *SensitiveManager) Startup(args map[string]interface{}) {
	fmt.Println("Startup SensitiveManager.")
}
func (manager *SensitiveManager) Shutdown(force bool) {
	fmt.Println("Shutdown SensitiveManager.")
}

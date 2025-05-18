package botmsg

import (
	"fmt"
	"im-server/commons/bases"
	"im-server/commons/gmicro"
	"im-server/commons/gmicro/actorsystem"
	"im-server/services/botmsg/actors"
)

type BotMsgManager struct{}

var serviceName string = "botmsgmanager"

func (manager *BotMsgManager) RegisterActors(register gmicro.IActorRegister) {
	register.RegisterActor("bot_msg", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.BotMsgActor{}, serviceName)
	})
}

func (manager *BotMsgManager) Startup(args map[string]interface{}) {
	fmt.Println("Startup botmsgmanager.")
}
func (manager *BotMsgManager) Shutdown(force bool) {
	fmt.Println("Shutdown botmsgmanager.")
}

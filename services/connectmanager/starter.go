package connectmanager

import (
	"fmt"

	"im-server/commons/bases"
	"im-server/commons/configures"
	"im-server/commons/gmicro"
	"im-server/commons/gmicro/actorsystem"
	"im-server/services/connectmanager/actors"
	"im-server/services/connectmanager/server"
)

var serviceName string = "connectmanager"

type ConnectManager struct {
	wsServer *server.ImWebsocketServer
}

func (ser *ConnectManager) RegisterActors(register gmicro.IActorRegister) {
	register.RegisterMultiMethodActor([]string{"connect", "msg", "ntf", "ustatus", "stream_msg"}, func() actorsystem.IUntypedActor {
		return &actors.ConnectActor{}
	}, 64)
	register.RegisterActor("qry_online_status", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.UserOnlineStatusActor{}, serviceName)
	}, 32)
	register.RegisterActor("ban_users", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.BanUsersActor{}, serviceName)
	}, 8)
	register.RegisterActor("qry_ban_users", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.QryBanUsersActor{}, serviceName)
	}, 8)
	register.RegisterActor("kick_user", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.KickUserActor{}, serviceName)
	}, 8)
}
func (ser *ConnectManager) Startup(args map[string]interface{}) {
	wsPort := configures.Config.ConnectManager.WsPort
	ser.wsServer = &server.ImWebsocketServer{
		MessageListener: &server.ImListenerImpl{},
	}
	go ser.wsServer.SyncStart(wsPort)
	fmt.Println("start with gorilla ws port:", wsPort)
}

func (ser *ConnectManager) Shutdown() {
	if ser.wsServer != nil {
		ser.wsServer.Stop()
	}
}

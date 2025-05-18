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
	register.RegisterMultiMethodActor([]string{"connect", "msg", "ntf", "ustatus", "stream_msg", "rtc_room_event", "rtc_invite_event"}, func() actorsystem.IUntypedActor {
		return &actors.ConnectActor{}
	})
	register.RegisterActor("qry_online_status", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.UserOnlineStatusActor{}, serviceName)
	})
	register.RegisterActor("ban_users", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.BanUsersActor{}, serviceName)
	})
	register.RegisterActor("qry_ban_users", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.QryBanUsersActor{}, serviceName)
	})
	register.RegisterActor("kick_user", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.KickUserActor{}, serviceName)
	})
}
func (ser *ConnectManager) Startup(args map[string]interface{}) {
	wsPort := configures.Config.ConnectManager.WsPort
	ser.wsServer = &server.ImWebsocketServer{
		MessageListener: &server.ImListenerImpl{},
	}
	go ser.wsServer.SyncStart(wsPort)
	fmt.Println("Start connectmanager with port:", wsPort)
}

func (ser *ConnectManager) Shutdown(force bool) {
	if ser.wsServer != nil {
		ser.wsServer.Stop()
	}
}

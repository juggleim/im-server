package rtcroom

import (
	"fmt"
	"im-server/commons/bases"
	"im-server/commons/gmicro"
	"im-server/commons/gmicro/actorsystem"
	"im-server/services/rtcroom/actors"
)

type RtcRoomManager struct{}

var serviceName string = "rtcroom"

func (manager *RtcRoomManager) RegisterActors(register gmicro.IActorRegister) {
	register.RegisterActor("rtc_create", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.CreateRoomActor{}, serviceName)
	})
	register.RegisterActor("rtc_destroy", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.DestroyRoomActor{}, serviceName)
	})
	register.RegisterActor("rtc_join", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.JoinRoomActor{}, serviceName)
	})
	register.RegisterActor("rtc_quit", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.QuitRoomActor{}, serviceName)
	})
	register.RegisterActor("rtc_qry", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.QryRtcRoomActor{}, serviceName)
	})
	register.RegisterActor("rtc_ping", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.PingActor{}, serviceName)
	})
	register.RegisterActor("rtc_upd_state", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.UpdateActor{}, serviceName)
	})

	register.RegisterActor("rtc_invite", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.InviteActor{}, serviceName)
	})
	register.RegisterActor("rtc_accept", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.AcceptActor{}, serviceName)
	})
	register.RegisterActor("rtc_hangup", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.HangupActor{}, serviceName)
	})
	register.RegisterActor("rtc_member_rooms", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.QryRtcMemberRoomsActor{}, serviceName)
	})
	register.RegisterActor("rtc_grab_member", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.GrabMemberActor{}, serviceName)
	})
	register.RegisterActor("rtc_sync_member", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.SyncMemberActor{}, serviceName)
	})
	register.RegisterActor("rtc_sync_user_quit", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.SyncQuitActor{}, serviceName)
	})
}

func (manager *RtcRoomManager) Startup(args map[string]interface{}) {
	fmt.Println("Startup rtcroom.")
}

func (manager *RtcRoomManager) Shutdown(force bool) {
	fmt.Println("Shutdown rtcroom.")
}

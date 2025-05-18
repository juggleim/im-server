package friends

import (
	"fmt"
	"im-server/commons/bases"
	"im-server/commons/gmicro"
	"im-server/commons/gmicro/actorsystem"
	"im-server/services/friends/actors"
)

var serviceName string = "friends"

type FriendManager struct{}

func (manager *FriendManager) RegisterActors(register gmicro.IActorRegister) {
	register.RegisterActor("add_friends", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.AddFriendActor{}, serviceName)
	})
	register.RegisterActor("del_friends", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.DelFriendActor{}, serviceName)
	})
	register.RegisterActor("qry_friends", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.QryFriendsActor{}, serviceName)
	})
	register.RegisterActor("qry_friends_with_page", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.QryFriendsWithPageActor{}, serviceName)
	})
	register.RegisterActor("check_friends", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.CheckFriendActor{}, serviceName)
	})
}

func (manager *FriendManager) Startup(args map[string]interface{}) {
	fmt.Println("Startup friendmanager.")
}

func (manager *FriendManager) Shutdown(force bool) {
	fmt.Println("Shutdown friendmanager.")
}

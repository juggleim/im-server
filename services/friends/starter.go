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
	register.RegisterActor("add_friend", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.AddFriendActor{}, serviceName)
	})
	register.RegisterActor("del_friend", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.DelFriendActor{}, serviceName)
	})
	register.RegisterActor("qry_friends", func() actorsystem.IUntypedActor {
		return bases.BaseProcessActor(&actors.QryFriendsActor{}, serviceName)
	})
}

func (manager *FriendManager) Startup(args map[string]interface{}) {
	fmt.Println("Startup friendmanager.")
}

func (manager *FriendManager) Shutdown() {
	fmt.Println("Shutdown friendmanager.")
}

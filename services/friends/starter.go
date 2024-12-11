package friends

import (
	"fmt"
	"im-server/commons/gmicro"
)

type FriendManager struct{}

func (manager *FriendManager) RegisterActors(register gmicro.IActorRegister) {

}

func (manager *FriendManager) Startup(args map[string]interface{}) {
	fmt.Println("Startup friendmanager.")
}

func (manager *FriendManager) Shutdown() {
	fmt.Println("Shutdown friendmanager.")
}

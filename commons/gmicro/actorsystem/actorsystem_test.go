package actorsystem

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"im-server/commons/gmicro/utils"

	"google.golang.org/protobuf/proto"
)

func TestNewActorSystemNoRpc(t *testing.T) {
	// actorSystem := NewActorSystemNoRpc("MyActorSystem")
	// actorSystem.RegisterActorProcessor("m1", func() proto.Message {
	// 	return &utils.Student{}
	// }, ActorSystemDemoProcessor, 10)

	// stu := &utils.Student{
	// 	Name: "name2",
	// 	Age:  1,
	// }
	// actor := actorSystem.LocalActorOf("m1")
	// actor.Tell(stu, NoSender)

	time.Sleep(5 * time.Second)
}

func ActorSystemDemoProcessor(sender ActorRef, input proto.Message) {
	fmt.Println("process has been executed.")
	fmt.Println("type:", reflect.TypeOf(input))
	stu := input.(*utils.Student)
	fmt.Println(stu.Name)

	sender.Tell(stu, NoSender)
}

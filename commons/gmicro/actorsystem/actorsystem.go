package actorsystem

import (
	"time"

	"im-server/commons/gmicro/utils"
)

const (
	NoRpcHost string = "-"
	NoRpcPort int    = 0
)

type ActorSystem struct {
	Name        string
	Host        string
	Prot        int
	sender      *MsgSender
	receiver    *MsgReceiver
	RecvDecoder func([]byte, interface{})
	dispatcher  *ActorDispatcher
}

func NewActorSystemNoRpc(name string) *ActorSystem {
	return NewActorSystem(name, NoRpcHost, NoRpcPort)
}

func NewActorSystem(name, host string, port int) *ActorSystem {
	sender := NewMsgSender()
	dispatcher := NewActorDispatcher(sender)
	receiver := NewMsgReceiver(host, port, dispatcher)

	sender.SetMsgReceiver(receiver)
	system := &ActorSystem{
		Name:       name,
		Host:       host,
		Prot:       port,
		sender:     sender,
		receiver:   receiver,
		dispatcher: dispatcher,
	}
	return system
}

func (system *ActorSystem) LocalActorOf(method string) ActorRef {
	return system.ActerOf(system.Host, system.Prot, method)
}

func (system *ActorSystem) ActerOf(host string, port int, method string) ActorRef {
	uid := utils.GenerateUUIDBytes()
	ref := NewActorRef(host, port, method, uid, system.sender)
	return ref
}

func (system *ActorSystem) CallbackActerOf(ttl time.Duration, actor ICallbackUntypedActor) ActorRef {
	uid := utils.GenerateUUIDBytes()
	ref := NewCallbackActorRef(ttl, system.Host, system.Prot, uid, actor, system.sender, system.dispatcher)
	return ref
}

func (system *ActorSystem) RegisterActor(method string, actorCreateFun func() IUntypedActor) {
	system.dispatcher.RegisterActor(method, actorCreateFun)
}

func (system *ActorSystem) RegisterMultiMethodActor(methods []string, actorCreateFun func() IUntypedActor) {
	system.dispatcher.RegisterMultiMethodActor(methods, actorCreateFun)
}

func (system *ActorSystem) RegisterStandaloneActor(method string, actorCreateFun func() IUntypedActor, concurrentCount int) {
	system.dispatcher.RegisterStandaloneActor(method, actorCreateFun, concurrentCount)
}

func (system *ActorSystem) RegisterStandaloneMultiMethodActor(methods []string, actorCreateFun func() IUntypedActor, concurrentCount int) {
	system.dispatcher.RegisterStandaloneMultiMethodActor(methods, actorCreateFun, concurrentCount)
}

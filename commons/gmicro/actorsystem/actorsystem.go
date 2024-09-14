package actorsystem

import (
	"time"

	"im-server/commons/gmicro/utils"
)

type ActorSystem struct {
	Name        string
	sender      *MsgSender
	receiver    *MsgReceiver
	RecvDecoder func([]byte, interface{})
	dispatcher  *ActorDispatcher
}

type MessageRequest struct {
	Type       int32
	Session    []byte
	TarMethod  string
	SrcMethod  string
	Data       []byte
	Extra      []byte
	TraceId    string
	IsNoSender bool
}

func NewActorSystem(name string) *ActorSystem {
	sender := NewMsgSender()
	dispatcher := NewActorDispatcher(sender)
	receiver := NewMsgReceiver(dispatcher)

	sender.SetMsgReceiver(receiver)
	system := &ActorSystem{
		Name:       name,
		sender:     sender,
		receiver:   receiver,
		dispatcher: dispatcher,
	}
	return system
}

func (system *ActorSystem) LocalActorOf(method string) ActorRef {
	return system.ActerOf(method)
}

func (system *ActorSystem) ActerOf(method string) ActorRef {
	uid := utils.GenerateUUIDBytes()
	ref := NewActorRef(method, uid, system.sender)
	return ref
}

func (system *ActorSystem) CallbackActerOf(ttl time.Duration, actor ICallbackUntypedActor) ActorRef {
	uid := utils.GenerateUUIDBytes()
	ref := NewCallbackActorRef(ttl, uid, actor, system.sender, system.dispatcher)
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

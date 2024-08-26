package actorsystem

import (
	"context"

	"google.golang.org/protobuf/proto"
)

type IUntypedActor interface{}

type IReceiveHandler interface {
	OnReceive(ctx context.Context, msg proto.Message)
}
type ISenderHandler interface {
	SetSender(sender ActorRef)
	GetSender() ActorRef
}
type ISelfHandler interface {
	SetSelf(self ActorRef)
	GetSelf() ActorRef
}
type ITimeoutHandler interface {
	OnTimeout()
}
type ICreateInputHandler interface {
	CreateInputObj() proto.Message
}

type IDefaultUntypedActor interface {
	IReceiveHandler
	ISenderHandler
	ISelfHandler
	ICreateInputHandler
}

type ICallbackUntypedActor interface {
	IDefaultUntypedActor
	ITimeoutHandler
}

type UntypedActor struct {
	Sender ActorRef
	Self   ActorRef
}

func (act *UntypedActor) SetSender(sender ActorRef) {
	act.Sender = sender
}

func (act *UntypedActor) SetSelf(self ActorRef) {
	act.Self = self
}

func (act *UntypedActor) GetSender() ActorRef {
	return act.Sender
}

func (act *UntypedActor) GetSelf() ActorRef {
	return act.Self
}

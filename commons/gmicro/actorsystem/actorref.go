package actorsystem

import (
	"time"

	"google.golang.org/protobuf/proto"
)

type ActorRef interface {
	Tell(proto.Message, ActorRef)
	TellAndNoSender(proto.Message)
	GetMethod() string
	isNoSender() bool
	isCallback() bool
	startCallbackActor(session []byte)
	getSession() []byte
}

type DefaultActorRef struct {
	Method        string
	Session       []byte
	Sender        *MsgSender
	is_Callback   bool
	callbackActor ICallbackUntypedActor
	ttl           time.Duration
	dispatcher    *ActorDispatcher
}

type DeadLetterActorRef struct {
	DefaultActorRef
}

func (ref *DeadLetterActorRef) TellAndNoSender(message proto.Message) {
	//do nothing
}

func (ref *DeadLetterActorRef) Tell(message proto.Message, sender ActorRef) {
	//do nothing
}
func (ref *DeadLetterActorRef) GetMethod() string {
	return ref.Method
}

func (ref *DeadLetterActorRef) getSession() []byte {
	return nil
}

func (ref *DeadLetterActorRef) isNoSender() bool {
	return true
}

var NoSender *DeadLetterActorRef

func init() {
	NoSender = &DeadLetterActorRef{}
}

func IsNoSender(req *MessageRequest) bool {
	if req != nil {
		return req.IsNoSender
	} else {
		return true
	}
}

func NewActorRef(method string, session []byte, sender *MsgSender) ActorRef {
	ref := &DefaultActorRef{
		Method:  method,
		Session: session,
		Sender:  sender,
	}
	return ref
}

func NewCallbackActorRef(ttl time.Duration, session []byte, callbackActor ICallbackUntypedActor, sender *MsgSender, dispatcher *ActorDispatcher) ActorRef {
	ref := &DefaultActorRef{
		Session:       session,
		Sender:        sender,
		is_Callback:   true,
		callbackActor: callbackActor,
		ttl:           ttl,
		dispatcher:    dispatcher,
	}
	return ref
}

func (ref *DefaultActorRef) TellAndNoSender(message proto.Message) {
	ref.Tell(message, NoSender)
}

func (ref *DefaultActorRef) Tell(message proto.Message, sender ActorRef) {
	if message != nil {
		bytes, _ := proto.Marshal(message)
		session := sender.getSession()
		if len(session) <= 0 {
			session = ref.Session
		}
		rpcReq := &MessageRequest{
			Session:   session,
			TarMethod: ref.Method,

			SrcMethod:  sender.GetMethod(),
			IsNoSender: sender.isNoSender(),

			Data: bytes,
		}
		if ref.is_Callback {
			ref.startCallbackActor(session)
		}
		if sender.isCallback() {
			sender.startCallbackActor(session)
		}
		ref.Sender.Send(rpcReq)
	}
}

func (ref *DefaultActorRef) startCallbackActor(session []byte) {
	if ref.callbackActor != nil && ref.dispatcher != nil {
		//start callback actor
		ref.dispatcher.AddCallbackActor(session, ref.callbackActor, ref.ttl)
	}
}
func (ref *DefaultActorRef) GetMethod() string {
	return ref.Method
}
func (ref *DefaultActorRef) isNoSender() bool {
	return false
}
func (ref *DefaultActorRef) isCallback() bool {
	return ref.is_Callback
}
func (ref *DefaultActorRef) getSession() []byte {
	return ref.Session
}

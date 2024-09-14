package actorsystem

import (
	"context"
	"sync"
	"time"

	"im-server/commons/gmicro/logs"
	"im-server/commons/gmicro/utils"

	"github.com/Jeffail/tunny"
	timewheel "github.com/rfyiamcool/go-timewheel"
	"google.golang.org/protobuf/proto"
)

const buffersize int = 8192

type ActorDispatcher struct {
	dispatchMap        sync.Map
	callbackMap        sync.Map
	msgSender          *MsgSender
	timer              *timewheel.TimeWheel
	callbackPool       *tunny.Pool
	callbackWraperChan chan wraper

	executorCommonPool *tunny.Pool
}

func NewActorDispatcher(sender *MsgSender) *ActorDispatcher {
	timer, err := timewheel.NewTimeWheel(1*time.Second, 360)
	if err != nil {
		logs.Error("error when start timewheel of dispatcher")
	}
	dispatcher := &ActorDispatcher{
		msgSender:          sender,
		timer:              timer,
		callbackPool:       tunny.NewCallback(buffersize),
		callbackWraperChan: make(chan wraper, buffersize),
		executorCommonPool: tunny.NewCallback(2 * buffersize),
	}
	timer.Start()
	go callbackActorExecute(dispatcher.callbackPool, dispatcher.callbackWraperChan)
	return dispatcher
}

func (dispatcher *ActorDispatcher) Dispatch(req *MessageRequest) {
	targetMethod := req.TarMethod
	var executor IExecutor

	if targetMethod == "" { //callback actor
		key := utils.Bytes2ShortString(req.Session)
		obj, ok := dispatcher.callbackMap.LoadAndDelete(key)
		if ok {
			callbackExecutor := obj.(*CallbackActorExecutor)
			//remove from timer task
			task := callbackExecutor.Task
			if task != nil {
				dispatcher.timer.Remove(task)
			}
			executor = callbackExecutor
		}
	} else {
		obj, ok := dispatcher.dispatchMap.Load(targetMethod)
		if ok {
			executor = obj.(IExecutor)
		}
	}
	if executor != nil {
		executor.Execute(req, dispatcher.msgSender)
	}
}

func (dispatcher *ActorDispatcher) Destroy() {
	if dispatcher.timer != nil {
		dispatcher.timer.Stop()
	}
}

func (dispatcher *ActorDispatcher) RegisterActor(method string, actorCreateFun func() IUntypedActor) {
	executor := NewActorExecutorWithDefaultPool(dispatcher.executorCommonPool, actorCreateFun)
	dispatcher.dispatchMap.Store(method, executor)
}

func (dispatcher *ActorDispatcher) RegisterStandaloneActor(method string, actorCreateFun func() IUntypedActor, concurrentCount int) {
	var executor *ActorExecutor
	if concurrentCount > 0 {
		executor = NewActorExecutor(concurrentCount, actorCreateFun)
	} else {
		executor = NewActorExecutorWithDefaultPool(dispatcher.executorCommonPool, actorCreateFun)
	}
	dispatcher.dispatchMap.Store(method, executor)
}

func (dispatcher *ActorDispatcher) RegisterMultiMethodActor(methods []string, actorCreateFun func() IUntypedActor) {
	executor := NewActorExecutorWithDefaultPool(dispatcher.executorCommonPool, actorCreateFun)
	for _, method := range methods {
		dispatcher.dispatchMap.Store(method, executor)
	}
}

func (dispatcher *ActorDispatcher) RegisterStandaloneMultiMethodActor(methods []string, actorCreateFun func() IUntypedActor, concurrentCount int) {
	var executor *ActorExecutor
	if concurrentCount > 0 {
		executor = NewActorExecutor(concurrentCount, actorCreateFun)
	} else {
		executor = NewActorExecutorWithDefaultPool(dispatcher.executorCommonPool, actorCreateFun)
	}
	for _, method := range methods {
		dispatcher.dispatchMap.Store(method, executor)
	}
}

func (dispatcher *ActorDispatcher) AddCallbackActor(session []byte, actor ICallbackUntypedActor, ttl time.Duration) {
	executor := NewCallbackActorExecutor(dispatcher.callbackPool, dispatcher.callbackWraperChan, actor)
	key := utils.Bytes2ShortString(session)
	dispatcher.callbackMap.Store(key, executor)
	task := dispatcher.timer.Add(ttl, func() {
		obj, ok := dispatcher.callbackMap.LoadAndDelete(key)
		if ok {
			executor := obj.(*CallbackActorExecutor)
			executor.doTimeout()
		}
	})
	executor.Task = task
}

func commonExecute(req *MessageRequest, msgSender *MsgSender, actor IUntypedActor) wraper {
	var sender ActorRef

	// srcHost := req.SrcHost
	// srcPort := req.SrcPort
	srcMethod := req.SrcMethod
	srcSession := req.Session

	if IsNoSender(req) {
		sender = NoSender
	} else {
		sender = &DefaultActorRef{
			// Host:    srcHost,
			// Port:    int(srcPort),
			Method:  srcMethod,
			Session: srcSession,
			Sender:  msgSender,
		}
	}

	bytes := req.Data

	createInputHandler, ok := actor.(ICreateInputHandler)
	var input proto.Message
	if ok {
		input = createInputHandler.CreateInputObj()
		proto.Unmarshal(bytes, input)
	}
	return wraper{
		sender: sender,
		msg:    input,
		actor:  actor,
	}
}

type wraper struct {
	sender ActorRef
	msg    proto.Message
	actor  IUntypedActor
}

func callbackActorExecute(pool *tunny.Pool, callbackWraperChan chan wraper) {
	for {
		wrapper := <-callbackWraperChan
		go pool.Process(func() {
			actor := wrapper.actor

			senderHandler, ok := actor.(ISenderHandler)
			if ok {
				senderHandler.SetSender(wrapper.sender)
			}
			receiveHandler, ok := actor.(IReceiveHandler)
			if ok {
				receiveHandler.OnReceive(context.Background(), wrapper.msg)
			}
		})
	}
}

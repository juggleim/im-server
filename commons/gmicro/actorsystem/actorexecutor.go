package actorsystem

import (
	"context"
	"sync"

	"im-server/commons/gmicro/utils"

	"github.com/Jeffail/tunny"
)

type IExecutor interface {
	Execute(req *MessageRequest, msgSender *MsgSender)
}

type ActorExecutor struct {
	wraperChan   chan wraper
	executePool  *tunny.Pool
	processSlots chan struct{}
	actorPool    sync.Pool
}

func NewActorExecutorWithDefaultPool(pool *tunny.Pool, actorCreateFun func() IUntypedActor) *ActorExecutor {
	return newActorExecutor(pool, make(chan struct{}, pool.GetSize()), actorCreateFun)
}

func newActorExecutor(pool *tunny.Pool, processSlots chan struct{}, actorCreateFun func() IUntypedActor) *ActorExecutor {
	executor := &ActorExecutor{
		wraperChan:   make(chan wraper, buffersize),
		executePool:  pool,
		processSlots: processSlots,
		actorPool: sync.Pool{
			New: func() any {
				return actorCreateFun()
			},
		},
	}
	go actorExecute(executor)
	return executor
}

func NewActorExecutor(concurrentCount int, actorCreateFun func() IUntypedActor) *ActorExecutor {
	executor := &ActorExecutor{
		wraperChan:   make(chan wraper, buffersize),
		executePool:  tunny.NewCallback(concurrentCount),
		processSlots: make(chan struct{}, concurrentCount),
		actorPool: sync.Pool{
			New: func() interface{} {
				return actorCreateFun()
			},
		},
	}
	go actorExecute(executor)
	return executor
}

func (executor *ActorExecutor) Execute(req *MessageRequest, msgSender *MsgSender) {
	actorObj := executor.actorPool.Get()
	executor.wraperChan <- commonExecute(req, msgSender, actorObj)
	executor.actorPool.Put(actorObj)
}

func actorExecute(executor *ActorExecutor) {
	for {
		wrapper := <-executor.wraperChan
		executor.processSlots <- struct{}{}
		go func(wrapper wraper) {
			defer func() { <-executor.processSlots }()
			executor.executePool.Process(func() {
				defer utils.Recovery()

				actorObj := executor.actorPool.Get()

				senderHandler, ok := actorObj.(ISenderHandler)
				if ok {
					senderHandler.SetSender(wrapper.sender)
				}

				receiveHandler, ok := actorObj.(IReceiveHandler)
				if ok {
					receiveHandler.OnReceive(context.Background(), wrapper.msg)
				}
				executor.actorPool.Put(actorObj)
			})
		}(wrapper)
	}
}

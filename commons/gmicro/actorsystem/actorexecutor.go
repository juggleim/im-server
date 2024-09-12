package actorsystem

import (
	"context"
	"sync"

	"im-server/commons/gmicro/actorsystem/rpc"
	"im-server/commons/gmicro/utils"

	"github.com/Jeffail/tunny"
)

type IExecutor interface {
	Execute(req *rpc.RpcMessageRequest, msgSender *MsgSender)
}

type ActorExecutor struct {
	wraperChan  chan wraper
	executePool *tunny.Pool
	actorPool   sync.Pool
	// actorCreateFun func() IUntypedActor
}

func NewActorExecutorWithDefaultPool(pool *tunny.Pool, actorCreateFun func() IUntypedActor) *ActorExecutor {
	executor := &ActorExecutor{
		wraperChan:  make(chan wraper, buffersize),
		executePool: pool,
		actorPool: sync.Pool{
			New: func() any {
				return actorCreateFun()
			},
		},
		// actorCreateFun: actorCreateFun,
	}
	go actorExecute(executor)
	return executor
}

func NewActorExecutor(concurrentCount int, actorCreateFun func() IUntypedActor) *ActorExecutor {
	executor := &ActorExecutor{
		wraperChan:  make(chan wraper, buffersize),
		executePool: tunny.NewCallback(concurrentCount),
		actorPool: sync.Pool{
			New: func() interface{} {
				return actorCreateFun()
			},
		},
		// actorCreateFun: actorCreateFun,
	}
	go actorExecute(executor)
	return executor
}

func (executor *ActorExecutor) Execute(req *rpc.RpcMessageRequest, msgSender *MsgSender) {
	actorObj := executor.actorPool.Get()
	executor.wraperChan <- commonExecute(req, msgSender, actorObj)
	executor.actorPool.Put(actorObj)
}

func actorExecute(executor *ActorExecutor) {
	for {
		wraper := <-executor.wraperChan
		go executor.executePool.Process(func() {
			defer utils.Recovery()

			actorObj := executor.actorPool.Get()

			senderHandler, ok := actorObj.(ISenderHandler)
			if ok {
				senderHandler.SetSender(wraper.sender)
			}

			receiveHandler, ok := actorObj.(IReceiveHandler)
			if ok {
				receiveHandler.OnReceive(context.Background(), wraper.msg)
			}
			executor.actorPool.Put(actorObj)
		})
	}
}

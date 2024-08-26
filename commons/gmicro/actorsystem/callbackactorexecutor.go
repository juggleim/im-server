package actorsystem

import (
	"im-server/commons/gmicro/actorsystem/rpc"

	"github.com/Jeffail/tunny"
	"github.com/rfyiamcool/go-timewheel"
)

type CallbackActorExecutor struct {
	Task         *timewheel.Task
	wraperChan   chan wraper
	callbackPool *tunny.Pool
	actor        ICallbackUntypedActor
}

func NewCallbackActorExecutor(callbackPool *tunny.Pool, wraperChan chan wraper, actor ICallbackUntypedActor) *CallbackActorExecutor {
	executor := &CallbackActorExecutor{
		wraperChan:   wraperChan,
		callbackPool: callbackPool,
		actor:        actor,
	}
	return executor
}

func (executor *CallbackActorExecutor) Execute(req *rpc.RpcMessageRequest, msgSender *MsgSender) {
	executor.wraperChan <- commonExecute(req, msgSender, executor.actor)
}

func (executor *CallbackActorExecutor) doTimeout() {
	if executor.actor != nil {
		timeoutHandler, ok := executor.actor.(ITimeoutHandler)
		if ok {
			timeoutHandler.OnTimeout()
		}
	}
}

package actorsystem

type MsgReceiver struct {
	recQueue   chan *MessageRequest
	dispatcher *ActorDispatcher
}

func NewMsgReceiver(dispatcher *ActorDispatcher) *MsgReceiver {
	rec := &MsgReceiver{
		recQueue:   make(chan *MessageRequest, 10000),
		dispatcher: dispatcher,
	}
	//start receiver queue
	go rec.start()
	return rec
}

func (rec *MsgReceiver) Receive(req *MessageRequest) {
	if req != nil {
		rec.recQueue <- req
	}
}

func (rec *MsgReceiver) start() {
	for {
		req := <-rec.recQueue
		rec.dispatcher.Dispatch(req)
	}
}

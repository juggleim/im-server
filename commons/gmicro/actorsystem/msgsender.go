package actorsystem

type MsgSender struct {
	msgReceiver *MsgReceiver
}

func NewMsgSender() *MsgSender {
	send := &MsgSender{}
	return send
}

func (sender *MsgSender) SetMsgReceiver(receiver *MsgReceiver) {
	sender.msgReceiver = receiver
}

func (sender *MsgSender) Send(req *MessageRequest) {
	if req != nil {
		sender.msgReceiver.Receive(req)
	}
}

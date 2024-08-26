package actorsystem

import (
	"strconv"
	"sync"

	"im-server/commons/gmicro/actorsystem/rpc"
)

type MsgSender struct {
	sendQueue   chan *rpc.RpcMessageRequest
	clientMap   sync.Map
	msgReceiver *MsgReceiver
}

func NewMsgSender() *MsgSender {
	send := &MsgSender{
		sendQueue: make(chan *rpc.RpcMessageRequest, 10000),
	}
	go send.start()
	return send
}

func (sender *MsgSender) SetMsgReceiver(receiver *MsgReceiver) {
	sender.msgReceiver = receiver
}

func (sender *MsgSender) Send(req *rpc.RpcMessageRequest) {
	if req != nil {
		isMatchReceiver := sender.msgReceiver.isMatch(req.TarHost, int(req.TarPort))
		if isMatchReceiver {
			sender.msgReceiver.Receive(req)
		} else {
			sender.sendQueue <- req
		}
	}
}

func (sender *MsgSender) start() {
	for {
		req := <-sender.sendQueue
		strKey := getTargetSign(req.TarHost, int(req.TarPort))
		actual, _ := sender.clientMap.LoadOrStore(strKey, NewRpcClient(strKey))
		act := actual.(*RpcClient)
		act.Send(req)
	}
}

func getTargetSign(host string, port int) string {
	sign := host + ":" + strconv.Itoa(port)
	return sign
}

func (sender *MsgSender) Stop() {

}

package main

import (
	"fmt"
	"time"

	"im-server/commons/errs"
	"im-server/simulator/examples"
	"im-server/simulator/utils"
	"im-server/simulator/wsclients"
)

func main() {
	wsClient := wsclients.NewWsImClient("ws://127.0.0.1:9002", "appkey", examples.Token1, examples.OnMessage, examples.OnDisconnect)
	code, connAck := wsClient.Connect("network", "num")
	if code == utils.ClientErrorCode_Success && connAck.Code == int32(errs.IMErrorCode_SUCCESS) {
		fmt.Println("connect success.", connAck.UserId, connAck.Session)

		time.Sleep(1000 * time.Second)
		wsClient.Disconnect()
	} else {
		fmt.Println("result:", code)
	} //CYXf6GNeM
}

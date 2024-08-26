package main

import (
	"fmt"
	"time"

	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/simulator/examples"
	"im-server/simulator/utils"
	"im-server/simulator/wsclients"
)

func main() {
	token := "CgdrZWZ1a2V5GiCkfWmqjspTHh3dsWtk_f9lN0TCCo0PDiEEvyEYoD3SeQ=="
	wsClient := wsclients.NewWsImClient("wss://ws.routechat.im", "kefukey", token, examples.OnMessage, examples.OnDisconnect)
	code, connAck := wsClient.Connect("network", "num")
	if code == utils.ClientErrorCode_Success && connAck.Code == int32(errs.IMErrorCode_SUCCESS) {
		fmt.Println("connect success.", connAck.UserId, connAck.Session)

		examples.QryHistoryMsgs(wsClient, 0, "1826895576622534657", pbobjs.ChannelType_Group, 20)
		// examples.SyncMsgs(wsClient)
		// examples.SendPrivateMsg(wsClient, "userid2")
		// code, resp := wsClient.QryFirstUnreadMsg(&pbobjs.QryFirstUnreadMsgReq{
		// 	TargetId:    "userid2",
		// 	ChannelType: pbobjs.ChannelType_Private,
		// })
		// fmt.Println(code)
		// fmt.Println(tools.ToJson(resp))

		time.Sleep(1000 * time.Second)
		wsClient.Disconnect()
	} else {
		fmt.Println("result:", code)
	} //CYXf6GNeM
}

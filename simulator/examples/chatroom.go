package examples

import (
	"fmt"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/simulator/utils"
	"im-server/simulator/wsclients"
)

func SendChatMsg(wsClient *wsclients.WsImClient, chatId string) {
	if wsClient != nil && wsClient.GetState() == utils.State_Connected {
		upMsg := pbobjs.UpMsg{
			MsgType:    "txtmsg",
			MsgContent: []byte(`{"content":"msg_content"}`),
		}
		code, sendAck := wsClient.SendChatMsg(chatId, &upMsg)
		fmt.Println(code)
		if code == utils.ClientErrorCode_Success {
			fmt.Println(sendAck.Code, sendAck.MsgId, sendAck.Timestamp)
		}
	}
}
func AddChatAtt(wsClient *wsclients.WsImClient, chatId string) {
	if wsClient != nil && wsClient.GetState() == utils.State_Connected {
		att := &pbobjs.ChatAttReq{
			Key:     "key1",
			Value:   "value1",
			IsForce: true,
		}
		code, ack := wsClient.AddChatAtt(chatId, att)
		fmt.Println(code)
		if code == utils.ClientErrorCode_Success {
			fmt.Println(ack.Code, ack.MsgId, ack.Timestamp)
		}
	}
}

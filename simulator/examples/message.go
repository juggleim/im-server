package examples

import (
	"encoding/json"
	"fmt"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices/msgdefines"
	"im-server/simulator/utils"
	"im-server/simulator/wsclients"
)

func SendPrivateMsg(wsClient *wsclients.WsImClient, targetId string) {
	if wsClient != nil && wsClient.GetState() == utils.State_Connected {
		flag := msgdefines.SetStoreMsg(0)
		flag = msgdefines.SetCountMsg(flag)
		upMsg := pbobjs.UpMsg{
			MsgType:    "txtMsg",
			MsgContent: []byte(`{"content":"msg content"}`),
			Flags:      flag,
		}
		code, sendAck := wsClient.SendPrivateMsg(targetId, &upMsg)
		if code == utils.ClientErrorCode_Success {
			fmt.Println(sendAck.Code, sendAck.MsgId, sendAck.Timestamp, sendAck.MsgSeqNo)
		} else {
			fmt.Println("ResultCode:", sendAck.Code)
		}
	}
}

func SyncMsgs(wsClient *wsclients.WsImClient) {
	if wsClient != nil && wsClient.GetState() == utils.State_Connected {
		req := pbobjs.SyncMsgReq{
			SyncTime:        0,
			SendBoxSyncTime: 0,
			ContainsSendBox: true,
		}
		code, msgs := wsClient.SyncMsgs(&req)
		fmt.Print("sync_msg_code:", code, "\t")
		if msgs != nil {
			fmt.Println("msg_count:", len(msgs.Msgs), "is_finished:", msgs.IsFinished)
		} else {
			fmt.Println("msgs is nil")
		}

		if code == utils.ClientErrorCode_Success && msgs != nil {
			for _, m := range msgs.Msgs {
				fmt.Println("*****************************")
				bs, _ := json.Marshal(m)
				fmt.Println(string(bs))
				fmt.Println(string(m.MsgContent))
			}
			fmt.Println("count:", len(msgs.Msgs), msgs.IsFinished)
		}
	}
}

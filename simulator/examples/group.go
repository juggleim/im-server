package examples

import (
	"fmt"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices/msgdefines"
	"im-server/simulator/utils"
	"im-server/simulator/wsclients"
)

func SendGroupMsg(wsClient *wsclients.WsImClient, groupId string) {
	if wsClient != nil && wsClient.GetState() == utils.State_Connected {
		flag := msgdefines.SetCountMsg(0)
		upMsg := pbobjs.UpMsg{
			MsgType:    "text",
			MsgContent: []byte(`{"content":"msg content"}`),
			Flags:      msgdefines.SetStoreMsg(flag),

			MentionInfo: &pbobjs.MentionInfo{
				MentionType: pbobjs.MentionType_Someone,
				TargetUsers: []*pbobjs.UserInfo{},
			},
		}
		upMsg.MentionInfo.TargetUsers = append(upMsg.MentionInfo.TargetUsers, &pbobjs.UserInfo{UserId: "userid2"})
		code, sendAck := wsClient.SendGroupMsg(groupId, &upMsg)
		fmt.Println(code)
		if code == utils.ClientErrorCode_Success {
			fmt.Println(sendAck.Code, sendAck.MsgId, sendAck.Timestamp)
		}
	}
}
func SendGroupMsgWithMention(wsClient *wsclients.WsImClient, groupId string) {
	if wsClient != nil && wsClient.GetState() == utils.State_Connected {
		upMsg := pbobjs.UpMsg{
			MsgType:    "txtmsg",
			MsgContent: []byte(`{"content":"msg content"}`),
			Flags:      msgdefines.SetStoreMsg(0),
			MentionInfo: &pbobjs.MentionInfo{
				MentionType: pbobjs.MentionType_All,
			},
		}
		code, sendAck := wsClient.SendGroupMsg(groupId, &upMsg)
		fmt.Println(code)
		if code == utils.ClientErrorCode_Success {
			fmt.Println(sendAck.Code, sendAck.MsgId, sendAck.Timestamp)
		}
	}
}

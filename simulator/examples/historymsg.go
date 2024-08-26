package examples

import (
	"encoding/json"
	"fmt"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/simulator/utils"
	"im-server/simulator/wsclients"
	"time"
)

func ModifyMsg(wsClient *wsclients.WsImClient, targetId string, channelType pbobjs.ChannelType, msgId string, msgTime int64, newContent string) {
	if wsClient != nil && wsClient.GetState() == utils.State_Connected {
		code, ack := wsClient.ModifyMsg(&pbobjs.ModifyMsgReq{
			TargetId:    targetId,
			ChannelType: channelType,
			MsgId:       msgId,
			MsgTime:     msgTime,
			MsgContent:  []byte(newContent),
		})
		fmt.Println("code:", code, "ack:", tools.ToJson(ack))
	}
}

func RecallMsg(wsClient *wsclients.WsImClient, targetId string, channelType pbobjs.ChannelType, msgId string, msgTime int64) {
	if wsClient != nil && wsClient.GetState() == utils.State_Connected {
		code, recallAck := wsClient.RecallMsg(&pbobjs.RecallMsgReq{
			TargetId:    targetId,
			MsgId:       msgId,
			ChannelType: channelType,
			MsgTime:     msgTime,
		})
		fmt.Println("code:", code, "ack:", recallAck)
	}
}

func QryHistoryMsgs(wsClient *wsclients.WsImClient, startTime int64, targetId string, channelType pbobjs.ChannelType, count int32) {
	if wsClient != nil && wsClient.GetState() == utils.State_Connected {
		qryHisMsgs := pbobjs.QryHisMsgsReq{
			TargetId:    targetId,
			ChannelType: channelType,
			StartTime:   startTime,
			Count:       count,
			Order:       0,
		}

		code, downMsgSet := wsClient.QryHistoryMsgs(&qryHisMsgs)
		if code == utils.ClientErrorCode_Success && downMsgSet != nil {
			if len(downMsgSet.Msgs) > 0 {
				for _, downMsg := range downMsgSet.Msgs {

					fmt.Println("********************")
					bs, _ := json.Marshal(downMsg)
					fmt.Println(string(bs))
					fmt.Println(string(downMsg.MsgContent), time.UnixMilli(downMsg.MsgTime))

				}
			}
		} else {
			fmt.Println(code)
		}
	}
}

func MarkReadMsg(wsClient *wsclients.WsImClient, targetId string, channelType pbobjs.ChannelType, msgs []*pbobjs.SimpleMsg) {
	if wsClient != nil && wsClient.GetState() == utils.State_Connected {
		markReadReq := pbobjs.MarkReadReq{
			TargetId:    targetId,
			ChannelType: channelType,
			Msgs:        msgs,
		}
		code, ack := wsClient.MarkReadMsg(&markReadReq)
		fmt.Println("code:", code, "ack:", ack)
	}
}

func DelHisMsgs(wsClient *wsclients.WsImClient, targetId string, channelType pbobjs.ChannelType, msgs []*pbobjs.SimpleMsg) {
	if wsClient != nil && wsClient.GetState() == utils.State_Connected {
		delMsgs := &pbobjs.DelHisMsgsReq{
			TargetId:    targetId,
			ChannelType: channelType,
			Msgs:        []*pbobjs.SimpleMsg{},
		}
		for _, msg := range msgs {
			delMsgs.Msgs = append(delMsgs.Msgs, &pbobjs.SimpleMsg{
				MsgId:   msg.MsgId,
				MsgTime: msg.MsgTime,
			})
		}
		code := wsClient.DelHisMsgs(delMsgs)
		fmt.Println("code:", code)
	}
}

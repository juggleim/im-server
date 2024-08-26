package examples

import (
	"encoding/json"
	"fmt"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/simulator/utils"
	"im-server/simulator/wsclients"
)

func SetTopConvers(wsClient *wsclients.WsImClient) {
	if wsClient != nil && wsClient.GetState() == utils.State_Connected {
		req := &pbobjs.ConversationsReq{
			Conversations: []*pbobjs.Conversation{},
		}
		req.Conversations = append(req.Conversations, &pbobjs.Conversation{
			TargetId:    "userid2",
			ChannelType: pbobjs.ChannelType_Private,
			IsTop:       1,
		})
		req.Conversations = append(req.Conversations, &pbobjs.Conversation{
			TargetId:    "groupid1",
			ChannelType: pbobjs.ChannelType_Group,
			IsTop:       1,
		})
		code := wsClient.SetTopConvers(req)
		fmt.Println("code:", code)
	}
}

func UndisturbConversation(wsClient *wsclients.WsImClient, targetId string, channelType pbobjs.ChannelType, undisturb int32) {
	if wsClient != nil && wsClient.GetState() == utils.State_Connected {
		req := &pbobjs.UndisturbConversReq{
			Items: []*pbobjs.UndisturbConverItem{},
		}
		req.Items = append(req.Items, &pbobjs.UndisturbConverItem{
			TargetId:      targetId,
			ChannelType:   channelType,
			UndisturbType: undisturb,
		})
		code := wsClient.UndisturbConvers(req)
		fmt.Println("code:", code)
	}
}

func ClearUnread(wsClient *wsclients.WsImClient, targetId string, channelType pbobjs.ChannelType, latestReadIndex int64) {
	if wsClient != nil && wsClient.GetState() == utils.State_Connected {
		req := &pbobjs.ClearUnreadReq{
			Conversations: []*pbobjs.Conversation{},
		}
		req.Conversations = append(req.Conversations, &pbobjs.Conversation{
			TargetId:        targetId,
			ChannelType:     channelType,
			LatestReadIndex: latestReadIndex,
		})
		code := wsClient.ClearUnread(req)
		fmt.Println("code:", code)
	}
}

func QryConversation(wsClient *wsclients.WsImClient, targetId string, channelType pbobjs.ChannelType) {
	if wsClient != nil && wsClient.GetState() == utils.State_Connected {
		qry := &pbobjs.QryConverReq{
			TargetId:    targetId,
			ChannelType: channelType,
		}
		code, resp := wsClient.QryConversation(qry)
		fmt.Println("code:", code, "resp:", tools.ToJson(resp))
	}
}

func QryConversations(wsClient *wsclients.WsImClient) {
	if wsClient != nil && wsClient.GetState() == utils.State_Connected {
		qryConversations := pbobjs.QryConversationsReq{
			StartTime: 0,
			Count:     1,
			Order:     0,
		}
		code, conversations := wsClient.QryConversations(&qryConversations)
		if code == utils.ClientErrorCode_Success && conversations != nil {
			for _, conver := range conversations.Conversations {
				fmt.Println("++++++++++++++++++++++++++++++")
				fmt.Println(tools.ToJson(conver))
				fmt.Println(conver.TargetId, conver.SortTime, conver.UnreadCount, conver.LatestReadIndex, conver.LatestUnreadIndex)
				bs, _ := json.Marshal(conver.Msg)
				fmt.Println(string(bs))
			}
		} else {
			fmt.Println(code)
		}
	}
}

func SyncConversations(wsClient *wsclients.WsImClient) {
	code, resp := wsClient.SyncConversations(&pbobjs.QryConversationsReq{
		StartTime: 0,
		Count:     100,
	})
	fmt.Println("Result_code:", code)
	if resp != nil {
		for _, conver := range resp.Conversations {

			fmt.Println(tools.ToJson(conver))
			fmt.Println("*************")
			fmt.Println(string(conver.Msg.MsgContent))

		}
	}
}

func QryMentionMsgs(wsClient *wsclients.WsImClient) {
	code, resp := wsClient.QryMentionMsgs(&pbobjs.QryMentionMsgsReq{
		TargetId:    "Tp6nLyUKX",
		ChannelType: pbobjs.ChannelType_Group,
		StartTime:   0,
		Count:       10,
	})
	fmt.Println("code:", code)
	for _, msg := range resp.MentionMsgs {
		fmt.Println("************************")
		fmt.Println(tools.ToJson(msg))
	}
}

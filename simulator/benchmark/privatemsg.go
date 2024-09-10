package benchmark

import (
	"encoding/json"
	"fmt"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/commonservices"
	"im-server/simulator/utils"
	"im-server/simulator/wsclients"
	"time"
)

type Txt struct {
	Content string `json:"content"`
	Time    int64  `json:"time"`
}

func OnMessage(msg *pbobjs.DownMsg) {
	var txt Txt
	json.Unmarshal(msg.MsgContent, &txt)
	msgTime := time.UnixMilli(txt.Time)
	fmt.Println("Received_Msg sender:", msg.SenderId, "msg_id:", msg.MsgId, "cost:", time.Since(msgTime))
}

func Connecting(count int) map[string]*wsclients.WsImClient {
	connectMap := make(map[string]*wsclients.WsImClient)
	for i := 1; i <= count; i++ {
		userId := fmt.Sprintf("userid%d", i)
		token := createToken(Appkey, SecureKey, userId)
		client := wsclients.NewWsImClient(WsAddress, Appkey, token, OnMessage, nil, nil)
		code, ack := client.Connect("", "")
		if code != utils.ClientErrorCode_Success {
			fmt.Println("Failed to connect. user_id:", userId, "code:", code, "msg:", tools.ToJson(ack))
			continue
		}
		connectMap[userId] = client
	}
	return connectMap
}

func PrivateMsg1000() {
	connectMap := Connecting(1000)
	for i := 1; i <= 1000; i++ {
		index := i
		go func() {
			senderId := fmt.Sprintf("userid%d", index)
			tarInt := index + 1
			if tarInt > 1000 {
				tarInt = tarInt % 1000
			}
			targetId := fmt.Sprintf("userid%d", tarInt)

			client := connectMap[senderId]
			if client != nil {
				flag := commonservices.SetStoreMsg(0)
				flag = commonservices.SetCountMsg(flag)
				start := time.Now()
				code, resp := client.SendPrivateMsg(targetId, &pbobjs.UpMsg{
					MsgType:    "jg:text",
					MsgContent: []byte(fmt.Sprintf("{\"content\":\"hello\",\"time\":%d}", time.Now().UnixMilli())),
					Flags:      flag,
				})

				fmt.Println("sender:", senderId, "target:", targetId, "code:", code, "resp:", tools.ToJson(resp), "cost:", time.Since(start))
			}
		}()
	}
	time.Sleep(2 * time.Minute)
	for _, v := range connectMap {
		v.Disconnect()
	}
}

func PrivateMsg3000() {
	connectMap := Connecting(3000)
	for i := 1; i <= 3000; i++ {
		index := i
		go func() {
			senderId := fmt.Sprintf("userid%d", index)
			tarInt := index + 1
			if tarInt > 3000 {
				tarInt = tarInt % 3000
			}
			targetId := fmt.Sprintf("userid%d", tarInt)

			client := connectMap[senderId]
			if client != nil {
				flag := commonservices.SetStoreMsg(0)
				flag = commonservices.SetCountMsg(flag)
				start := time.Now()
				code, resp := client.SendPrivateMsg(targetId, &pbobjs.UpMsg{
					MsgType:    "jg:text",
					MsgContent: []byte(fmt.Sprintf("{\"content\":\"hello\",\"time\":%d}", time.Now().UnixMilli())),
					Flags:      flag,
				})

				fmt.Println("sender:", senderId, "target:", targetId, "code:", code, "resp:", tools.ToJson(resp), "cost:", time.Since(start))
			}
		}()
	}
	time.Sleep(10 * time.Minute)
	for _, v := range connectMap {
		v.Disconnect()
	}
}

func PrivateMsg3000_5() {
	connectMap := Connecting(3000)
	for i := 1; i <= 3000; i++ {
		index := i
		go func() {
			senderId := fmt.Sprintf("userid%d", index)
			tarInt := index + 1
			if tarInt > 3000 {
				tarInt = tarInt % 3000
			}
			targetId := fmt.Sprintf("userid%d", tarInt)

			client := connectMap[senderId]
			if client != nil {
				flag := commonservices.SetStoreMsg(0)
				flag = commonservices.SetCountMsg(flag)

				for y := 0; y < 20; y++ {
					for x := 0; x < 5; x++ {
						start := time.Now()
						code, resp := client.SendPrivateMsg(targetId, &pbobjs.UpMsg{
							MsgType:    "jg:text",
							MsgContent: []byte(fmt.Sprintf("{\"content\":\"hello\",\"time\":%d}", time.Now().UnixMilli())),
							Flags:      flag,
						})
						fmt.Println("sender:", senderId, "target:", targetId, "code:", code, "resp:", tools.ToJson(resp), "cost:", time.Since(start))
					}
					time.Sleep(500 * time.Millisecond)
				}
			}
		}()
	}
	time.Sleep(30 * time.Minute)
	for _, v := range connectMap {
		v.Disconnect()
	}
}

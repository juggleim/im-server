package benchmark

import (
	"fmt"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/commonservices/msgdefines"
	"time"
)

func GroupMsg1000() {
	connectMap := Connecting(1000)
	groupId := "group1000"
	for i := 1; i <= 100; i++ {
		senderId := fmt.Sprintf("userid%d", i)
		client := connectMap[senderId]
		flag := msgdefines.SetStoreMsg(0)
		flag = msgdefines.SetCountMsg(flag)
		start := time.Now()
		code, resp := client.SendGroupMsg(groupId, &pbobjs.UpMsg{
			MsgType:    msgdefines.InnerMsgType_Text,
			MsgContent: []byte(fmt.Sprintf("{\"content\":\"hello\",\"time\":%d}", time.Now().UnixMilli())),
			Flags:      flag,
		})
		fmt.Println("sender:", senderId, "target:", groupId, "code:", code, "resp:", tools.ToJson(resp), "cost:", time.Since(start))
		time.Sleep(1 * time.Second)
	}

	time.Sleep(30 * time.Minute)
	for _, v := range connectMap {
		v.Disconnect()
	}
}

func GroupMsg2000() {
	connectMap := Connecting(2000)
	groupId := "group2000"
	for i := 1; i <= 100; i++ {
		senderId := fmt.Sprintf("userid%d", i)
		client := connectMap[senderId]
		flag := msgdefines.SetStoreMsg(0)
		flag = msgdefines.SetCountMsg(flag)
		start := time.Now()
		code, resp := client.SendGroupMsg(groupId, &pbobjs.UpMsg{
			MsgType:    msgdefines.InnerMsgType_Text,
			MsgContent: []byte(fmt.Sprintf("{\"content\":\"hello\",\"time\":%d}", time.Now().UnixMilli())),
			Flags:      flag,
		})
		fmt.Println("sender:", senderId, "target:", groupId, "code:", code, "resp:", tools.ToJson(resp), "cost:", time.Since(start))
		time.Sleep(1 * time.Second)
	}

	time.Sleep(30 * time.Minute)
	for _, v := range connectMap {
		v.Disconnect()
	}
}

func GroupMsg3000() {
	connectMap := Connecting(3000)
	groupId := "group3000"
	for i := 1; i <= 100; i++ {
		senderId := fmt.Sprintf("userid%d", i)
		client := connectMap[senderId]
		flag := msgdefines.SetStoreMsg(0)
		flag = msgdefines.SetCountMsg(flag)
		start := time.Now()
		code, resp := client.SendGroupMsg(groupId, &pbobjs.UpMsg{
			MsgType:    msgdefines.InnerMsgType_Text,
			MsgContent: []byte(fmt.Sprintf("{\"content\":\"hello\",\"time\":%d}", time.Now().UnixMilli())),
			Flags:      flag,
		})
		fmt.Println("sender:", senderId, "target:", groupId, "code:", code, "resp:", tools.ToJson(resp), "cost:", time.Since(start))
		time.Sleep(1 * time.Second)
	}
	time.Sleep(30 * time.Minute)
	for _, v := range connectMap {
		v.Disconnect()
	}
}

func GroupMsg5000() {
	connectMap := Connecting(5000)
	groupId := "group5000"
	for i := 1; i <= 10; i++ {
		senderId := fmt.Sprintf("userid%d", i)
		client := connectMap[senderId]
		flag := msgdefines.SetStoreMsg(0)
		flag = msgdefines.SetCountMsg(flag)
		start := time.Now()
		code, resp := client.SendGroupMsg(groupId, &pbobjs.UpMsg{
			MsgType:    msgdefines.InnerMsgType_Text,
			MsgContent: []byte(fmt.Sprintf("{\"content\":\"hello\",\"time\":%d}", time.Now().UnixMilli())),
			Flags:      flag,
		})
		fmt.Println("sender:", senderId, "target:", groupId, "code:", code, "resp:", tools.ToJson(resp), "cost:", time.Since(start))
		time.Sleep(time.Second)
	}
	time.Sleep(10 * time.Minute)
	for _, v := range connectMap {
		v.Disconnect()
	}
}

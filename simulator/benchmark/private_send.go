package main

import (
	"fmt"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices"
	"im-server/simulator/wsclients"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

func privateSend() {
	var sendClients = make(map[int]*wsclients.WsImClient, 1000)
	var clientLocker sync.Mutex

	var connectWg sync.WaitGroup
	var sendWg sync.WaitGroup

	var (
		msgMap     = make(map[string]*time.Duration, 10000)
		msgMapLock sync.Mutex
	)

	var (
		sendClientNum = 500
		turnMsgCount  = 20
	)

	for i := 0; i < sendClientNum*2; i++ {
		connectWg.Add(1)

		go func(i int) {
			time.Sleep(time.Duration(rand.Int31n(1000)) * time.Millisecond)
			client := wsclients.NewWsImClient(wsUrl, appKey, generateUserTokenStr(userPrefix+strconv.Itoa(i)), func(msg *pbobjs.DownMsg) {
				msgMapLock.Lock()
				if _, ok := msgMap[msg.MsgId]; ok {
					sendTime := bytesToTime(msg.MsgContent)
					duration := time.Since(sendTime)
					msgMap[msg.MsgId] = &duration
				}
				msgMapLock.Unlock()

				//fmt.Printf("msgId:%s, senderId:%s, receiverId:%s, msgTime:%v, msgContent:%s\n", msg.MsgId, msg.SenderId, msg.TargetId, msg.MsgTime, msg.MsgContent)

			}, func(msg *pbobjs.StreamDownMsg) {
				fmt.Println("stream down msg:", msg)
			}, nil)
			code, _ := client.Connect("", "")
			if code != 0 {
				fmt.Printf("connect code: %d\n", code)
			}
			clientLocker.Lock()
			if _, ok := sendClients[i]; ok {
				panic("client already exists")
			}
			sendClients[i] = client
			clientLocker.Unlock()
			connectWg.Done()
		}(i)
	}
	connectWg.Wait()
	defer func() {
		for _, client := range sendClients {
			client.Disconnect()
		}
	}()
	fmt.Printf("连接创建完成，连接数 %d\n", len(sendClients))

	sendStart := time.Now()

	fmt.Printf("开始发送消息, 发送客户端数量 %d, 发送消息数量 %d\n", sendClientNum, turnMsgCount)
	for i := 0; i < sendClientNum; i++ {
		sendWg.Add(1)
		go func(i int) {
			client := sendClients[i]
			flag := commonservices.SetStoreMsg(0)
			flag = commonservices.SetCountMsg(flag)

			for j := 0; j < turnMsgCount; j++ {
				upMsg := pbobjs.UpMsg{
					MsgType:    "txtMsg",
					MsgContent: timeToBytes(time.Now()),
					Flags:      flag,
				}
				code, ack := client.SendPrivateMsg(userPrefix+strconv.Itoa(i+sendClientNum), &upMsg)
				if code != 0 {
					fmt.Printf("send upMsg failed, code: %d\n", code)
					return
				}

				msgMapLock.Lock()
				msgMap[ack.MsgId] = nil
				msgMapLock.Unlock()
			}
			sendWg.Done()
		}(i)
	}
	WaitTimeout(&sendWg, 10*time.Second)
	fmt.Printf("发送消息数量 %d, time used:%v\n", len(msgMap), time.Now().Sub(sendStart))

	s := time.Now()
	ticker := time.NewTicker(time.Second)

	prevTotal := 0
	for {
		select {
		case t := <-ticker.C:
			if t.Sub(s).Seconds() > 10 {
				return
			}

			msgMapLock.Lock()
			total, maxDelay, minDelay, avgDelay := statisticsMsgMap(msgMap)
			if total > prevTotal {
				fmt.Printf("收到消息数量 %d, 平均延迟 %v, 最大延迟 %v, 最小延迟 %v\n", total, avgDelay, maxDelay, minDelay)
				prevTotal = total
			}
			if total >= len(msgMap) {
				return
			}
			msgMapLock.Unlock()
		}
	}

}

func timeToBytes(t time.Time) []byte {
	milli := t.UnixMilli()
	return []byte(strconv.FormatInt(milli, 10))
}

func bytesToTime(b []byte) time.Time {
	milli, _ := strconv.ParseInt(string(b), 10, 64)
	return time.UnixMilli(milli)
}

func statisticsMsgMap(msgMap map[string]*time.Duration) (total int, maxDelay time.Duration, minDelay time.Duration, avgDelay time.Duration) {
	var totalDuration time.Duration
	for _, v := range msgMap {
		if v == nil {
			continue
		}
		total += 1
		if *v > maxDelay {
			maxDelay = *v
		}
		if *v < minDelay {
			minDelay = *v
		}
		totalDuration += *v
	}
	if total == 0 {
		return
	}
	avgDelay = totalDuration / time.Duration(total)
	return
}

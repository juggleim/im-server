package main

import (
	"fmt"
	"go.uber.org/atomic"
	"gopkg.in/yaml.v3"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices"
	"im-server/services/commonservices/tokens"
	"im-server/simulator/serversdk"
	"im-server/simulator/wsclients"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
)

const userPrefix = "benchmark_user_"

var appKey, secret, secureKey, apiUrl, wsUrl string

type Config struct {
	AppKey    string `yaml:"appKey"`
	Secret    string `yaml:"secret"`
	SecureKey string `yaml:"secureKey"`
	ApiUrl    string `yaml:"apiUrl"`
	WsUrl     string `yaml:"wsUrl"`
}

var config Config

func loadConfig() {
	file, err := os.ReadFile("config.yml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		panic(err)
	}
	appKey = config.AppKey
	secret = config.Secret
	secureKey = config.SecureKey
	apiUrl = config.ApiUrl
	wsUrl = config.WsUrl
}
func init() {
	loadConfig()
}

func main() {

	//RegisterUsers(userPrefix, 1)
	privateSend()
	//CreateGroup()
}

func calcTimeUsed(fn func()) {
	start := time.Now()
	fn()
	fmt.Println("Time used:", time.Now().Sub(start))
}

func generateUserTokenStr(userId string) string {
	token := tokens.ImToken{
		AppKey:    appKey,
		UserId:    userId,
		DeviceId:  "deviceid",
		TokenTime: time.Now().UnixMilli(),
	}
	secureKey := []byte(secureKey)
	tokenStr, _ := token.ToTokenString(secureKey)
	return tokenStr
}

func CreateGroup() {
	sdk := serversdk.NewJuggleIMSdk(appKey, secret, apiUrl)
	req := serversdk.GroupMembersReq{
		GroupId:       "benchmark_group1",
		GroupName:     "benchmark_group1",
		GroupPortrait: "",
		MemberIds:     nil,
	}
	for i := 0; i < 10000; i++ {
		req.MemberIds = append(req.MemberIds, userPrefix+strconv.Itoa(i))
	}
	code, _, err := sdk.CreateGroup(req)
	fmt.Printf("Create group, code: %d, err: %v\n", code, err)
}

func RegisterUsers(prefix string, count int) error {
	sdk := serversdk.NewJuggleIMSdk(appKey, secret, apiUrl)
	for i := 0; i < count; i++ {
		_, code, _, err := sdk.Register(serversdk.User{
			UserId:       prefix + strconv.Itoa(i),
			Nickname:     prefix + strconv.Itoa(i),
			UserPortrait: "",
			ExtFields:    nil,
		})
		fmt.Printf("Register user %s, code: %d, err: %v\n", prefix+strconv.Itoa(i), code, err)
	}

	return nil
}

func privateSend() {
	var sendClients = make(map[int]*wsclients.WsImClient, 1000)
	var clientLocker sync.Mutex

	var recvMsgs = make([]time.Duration, 0, 10000)
	var recvLocker sync.Mutex

	var sendCount = atomic.NewInt32(0)

	var connectWg sync.WaitGroup
	var sendWg sync.WaitGroup
	var recvWg sync.WaitGroup

	for i := 0; i < 500; i++ {
		connectWg.Add(1)

		go func(i int) {
			time.Sleep(time.Duration(rand.Int31n(1000)) * time.Millisecond)
			client := wsclients.NewWsImClient(wsUrl, appKey, generateUserTokenStr(userPrefix+strconv.Itoa(i)), func(msg *pbobjs.DownMsg) {
				sendTime := bytesToTime(msg.MsgContent)
				recvLocker.Lock()
				recvMsgs = append(recvMsgs, time.Since(sendTime))
				recvLocker.Unlock()

				recvWg.Done()
			}, nil, nil)
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

	lastSendNo := 50
	for i := 0; i < lastSendNo; i++ {
		var turnMsgCount = 300
		sendWg.Add(1)
		recvWg.Add(turnMsgCount)
		go func(i int) {
			client := sendClients[i]
			flag := commonservices.SetStoreMsg(0)
			flag = commonservices.SetStateMsg(flag)
			flag = commonservices.SetCountMsg(flag)

			for i := 0; i < turnMsgCount; i++ {
				upMsg := pbobjs.UpMsg{
					MsgType:    "txtMsg",
					MsgContent: timeToBytes(time.Now()),
					Flags:      flag,
				}
				code, _ := client.SendPrivateMsg(userPrefix+strconv.Itoa(i+lastSendNo), &upMsg)
				if code != 0 {
					fmt.Printf("send upMsg failed, code: %d\n", code)
				}
				sendCount.Add(1)
			}
			sendWg.Done()
		}(i)
	}
	WaitTimeout(&sendWg, 10*time.Second)
	fmt.Printf("发送消息数量 %d, time used:%v\n", sendCount.Load(), time.Now().Sub(sendStart))

	WaitTimeout(&recvWg, 10*time.Second)
	fmt.Printf("收到消息数量 %d, max:%v, min:%v, avg:%v\n", len(recvMsgs), maxDuration(recvMsgs), minDuration(recvMsgs), avgDuration(recvMsgs))

}

func timeToBytes(t time.Time) []byte {
	milli := t.UnixMilli()
	return []byte(strconv.FormatInt(milli, 10))
}

func bytesToTime(b []byte) time.Time {
	milli, _ := strconv.ParseInt(string(b), 10, 64)
	return time.UnixMilli(milli)
}

func maxDuration(durations []time.Duration) time.Duration {
	duration := time.Duration(0)
	for _, d := range durations {
		if duration < d {
			duration = d
		}
	}
	return duration
}

func minDuration(durations []time.Duration) time.Duration {
	duration := time.Duration(0)
	for _, d := range durations {
		if duration > d {
			duration = d
		}
	}
	return duration
}

func avgDuration(durations []time.Duration) time.Duration {
	if len(durations) == 0 {
		return 0
	}
	total := time.Duration(0)
	for _, d := range durations {
		total += d
	}
	return total / time.Duration(len(durations))
}

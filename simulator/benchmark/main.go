package main

import (
	"fmt"
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

func createGroup(groupId string, memberCount int) {
	sdk := serversdk.NewJuggleIMSdk(appKey, secret, apiUrl)
	req := serversdk.GroupMembersReq{
		GroupId:       groupId,
		GroupName:     groupId,
		GroupPortrait: "",
		MemberIds:     nil,
	}
	for i := 0; i < memberCount; i++ {
		req.MemberIds = append(req.MemberIds, userPrefix+strconv.Itoa(i))
	}
	code, _, err := sdk.CreateGroup(req)
	fmt.Printf("Create group, code: %d, err: %v\n", code, err)
}

func dissolveGroup(groupId string) {
	sdk := serversdk.NewJuggleIMSdk(appKey, secret, apiUrl)
	code, _, err := sdk.DissolveGroup(groupId)
	fmt.Printf("Disolve group, code: %d, err: %v\n", code, err)
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

func statisticsMsgMap(msgMap map[string]*groupMessageData) (total int, maxDelay time.Duration, minDelay time.Duration, avgDelay time.Duration) {
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
	avgDelay = totalDuration / time.Duration(total)
	return
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

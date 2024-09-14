package benchmark

import (
	"fmt"
	"im-server/commons/tools"
	"im-server/services/commonservices/tokens"
	"im-server/services/connectmanager/server/codec"
	"im-server/simulator/utils"
	"im-server/simulator/wsclients"
	"time"
)

var (
	Appkey    string = "appkey"
	SecureKey string = "abcdefghijklmnop"
	WsAddress string = "ws://172.25.149.199:9002"
)

func OnDisconnect(code utils.ClientErrorCode, disMsg *codec.DisconnectMsgBody) {
	fmt.Println("disconnect:", tools.ToJson(disMsg))
}

func CreateToken(appkey, secureKey, userid string) string {
	imToken := &tokens.ImToken{
		AppKey:    appkey,
		UserId:    userid,
		TokenTime: time.Now().UnixMilli(),
	}
	tokenStr, err := imToken.ToTokenString([]byte(secureKey))
	if err != nil {
		return ""
	}
	return tokenStr
}
func connect(count int) {
	for i := 1; i <= count; i++ {
		userId := fmt.Sprintf("userid%d", i)
		token := CreateToken(Appkey, SecureKey, userId)
		if token != "" {
			client := wsclients.NewWsImClient(WsAddress, Appkey, token, nil, nil, OnDisconnect)
			start := time.Now()
			code, connectAckMsg := client.Connect("nettwork", "ispNum")
			if code != utils.ClientErrorCode_Success {
				fmt.Println("Failed to connect. user_id:", userId, "code:", code, "msg:", tools.ToJson(connectAckMsg))
			} else {
				fmt.Println("Success to connect. user_id:", userId, time.Since(start))
			}
		} else {
			fmt.Println("Failed to create connect token. user_id:", userId)
		}
	}
}

func Connect4000() {
	connect(4000)
	time.Sleep(10 * time.Minute)
}

func Connect1000() {
	connect(1000)
	time.Sleep(1 * time.Minute)
}

func Connect5000() {
	connect(5000)
	time.Sleep(10 * time.Minute)
}

func Connect10000() {
	connect(10000)
	time.Sleep(10 * time.Minute)
}

func Connect50000() {
	connect(50000)
	time.Sleep(10 * time.Minute)
}

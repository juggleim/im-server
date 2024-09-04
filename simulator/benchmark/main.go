package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"im-server/services/commonservices/tokens"
	"im-server/simulator/serversdk"
	"os"
	"strconv"
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
	//groupSend()
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

func registerUsers(prefix string, count int) error {
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

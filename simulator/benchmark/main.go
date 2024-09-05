package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"im-server/services/commonservices/tokens"
	"im-server/simulator/serversdk"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/urfave/cli/v2"
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
	app := &cli.App{
		Name:  "benchmark",
		Usage: "sdk压测",
		Commands: []*cli.Command{
			{
				Name:  "registerUsers",
				Usage: "注册压测用户",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:     "userNum",
						Usage:    "用户数量",
						Required: true,
					},
				},
				Action: func(c *cli.Context) error {
					userNum := c.Int("userNum")
					_ = registerUsers(userPrefix, userNum)
					return nil
				},
			},
			{
				Name:  "sendPrivateMsg",
				Usage: "发送私聊消息",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:     "connNum",
						Usage:    "在线连接数",
						Required: true,
					},
					&cli.IntFlag{
						Name:     "sendNum",
						Usage:    "发送用户数量",
						Required: true,
					},
					&cli.IntFlag{
						Name:     "turnCount",
						Usage:    "每个用户发送消息数量",
						Required: true,
					},
					&cli.IntFlag{
						Name:     "timeout",
						Usage:    "超时时间",
						Value:    10,
						Required: false,
					},
				},
				Action: func(c *cli.Context) error {
					connNum := c.Int("connNum")
					sendNum := c.Int("sendNum")
					turnCount := c.Int("turnCount")
					timeout := c.Int("timeout")
					privateSend(connNum, sendNum, turnCount, timeout)
					return nil
				},
			},
			{
				Name:  "sendGroupMsg",
				Usage: "发送群组消息",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:     "groupMemberCount",
						Usage:    "群组成员数量",
						Required: true,
					},
					&cli.IntFlag{
						Name:     "sendNum",
						Usage:    "发送用户数量",
						Required: true,
					},
					&cli.IntFlag{
						Name:     "turnCount",
						Usage:    "每个用户发送消息数量",
						Required: true,
					},
					&cli.IntFlag{
						Name:     "timeout",
						Usage:    "超时时间",
						Value:    30,
						Required: false,
					},
				},
				Action: func(c *cli.Context) error {
					groupMemberCount := c.Int("groupMemberCount")
					sendNum := c.Int("sendNum")
					turnCount := c.Int("turnCount")
					timeout := c.Int("timeout")
					groupSend(groupMemberCount, sendNum, turnCount, timeout)
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
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

package interceptors

import (
	"context"
	"fmt"
	"im-server/commons/bases"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"strings"
	"time"
)

type CustomInterceptor struct {
	AppKey     string
	AppSecret  string
	RequestUrl string
	Conditions []*Condition
}

func (inter *CustomInterceptor) GetConditions() []*Condition {
	return inter.Conditions
}

func (inter *CustomInterceptor) CheckMsgInterceptor(ctx context.Context, senderId, receiverId string, channelType pbobjs.ChannelType, msg *pbobjs.UpMsg) bool {
	if bases.GetIsFromApiFromCtx(ctx) {
		return false
	}
	appkey := bases.GetAppKeyFromCtx(ctx)
	nonce := tools.RandStr(8)
	tsStr := fmt.Sprintf("%d", time.Now().UnixMilli())
	headers := map[string]string{
		"Content-Type": "application/json",
		"appkey":       appkey,
		"nonce":        nonce,
		"timestamp":    tsStr,
		"signature":    tools.SHA1(fmt.Sprintf("%s%s%s", inter.AppSecret, nonce, tsStr)),
	}
	msgEvent := &MsgEvent{
		Platform:    bases.GetPlatformFromCtx(ctx),
		Sender:      senderId,
		Receiver:    receiverId,
		ChannelType: int(channelType),
		MsgType:     msg.MsgType,
		MsgContent:  string(msg.MsgContent),
	}
	body := tools.ToJson(msgEvent)
	respBs, code, err := tools.HttpDoBytes("POST", inter.RequestUrl, headers, body)
	if err != nil {
		fmt.Println("xxx:", err)
		return false
	}
	if code != 200 {
		fmt.Println("xxx:", code)
		return false
	}
	resp := &CustomInterceptorResp{}
	err = tools.JsonUnMarshal(respBs, resp)
	if err != nil {
		return false
	}
	if strings.ToLower(resp.Result) == "pass" {
		return false
	} else {
		return true
	}
}

type CustomInterceptorResp struct {
	Result string `json:"result"`
}

type MsgEvent struct {
	Platform    string `json:"platform"`
	Sender      string `json:"sender"`
	Receiver    string `json:"receiver"`
	ChannelType int    `json:"channel_type"`
	MsgType     string `json:"msg_type"`
	MsgContent  string `json:"msg_content"`
}

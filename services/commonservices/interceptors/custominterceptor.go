package interceptors

import (
	"context"
	"fmt"
	"im-server/commons/bases"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/commonservices/msgdefines"
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

func (inter *CustomInterceptor) CheckMsgInterceptor(ctx context.Context, senderId, receiverId string, channelType pbobjs.ChannelType, msg *pbobjs.UpMsg) (InterceptorResult, int64) {
	if bases.GetIsFromApiFromCtx(ctx) {
		return InterceptorResult_Pass, 0
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
		return InterceptorResult_Pass, 0
	}
	if code != 200 {
		return InterceptorResult_Pass, 0
	}
	resp := &CustomInterceptorResp{}
	err = tools.JsonUnMarshal(respBs, resp)
	if err != nil {
		return InterceptorResult_Pass, 0
	}
	result := strings.ToLower(resp.Result)
	customCode := resp.CustomCode
	if result == "pass" {
		return InterceptorResult_Pass, 0
	} else if result == "replace" {
		if msg != nil {
			if resp.MsgType == "" && resp.MsgContent == "" {
				return InterceptorResult_Pass, customCode
			}
			if resp.MsgType != "" {
				msg.MsgType = resp.MsgType
				msg.Flags = msgdefines.SetModifiedMsg(msg.Flags)
			}
			if resp.MsgContent != "" {
				msg.MsgContent = []byte(resp.MsgContent)
				msg.Flags = msgdefines.SetModifiedMsg(msg.Flags)
			}
			return InterceptorResult_Replace, customCode
		}
		return InterceptorResult_Pass, customCode
	} else if result == "reject" {
		return InterceptorResult_Reject, customCode
	} else if result == "silent" {
		return InterceptorResult_Silent, customCode
	} else {
		return InterceptorResult_Pass, customCode
	}
}

type CustomInterceptorResp struct {
	Result     string `json:"result"`
	MsgType    string `json:"msg_type"`
	MsgContent string `json:"msg_content"`
	CustomCode int64  `json:"custom_code"`
}

type MsgEvent struct {
	Platform    string `json:"platform"`
	Sender      string `json:"sender"`
	Receiver    string `json:"receiver"`
	ChannelType int    `json:"channel_type"`
	MsgType     string `json:"msg_type"`
	MsgContent  string `json:"msg_content"`
}

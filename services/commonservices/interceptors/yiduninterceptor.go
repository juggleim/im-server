package interceptors

import (
	"context"
	"fmt"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/commonservices/msgdefines"

	"github.com/yidun/yidun-golang-sdk/yidun/core/http"
	v5 "github.com/yidun/yidun-golang-sdk/yidun/service/antispam/text"
	"github.com/yidun/yidun-golang-sdk/yidun/service/antispam/text/v5/check/sync/single"
)

type YidunInterceptorConf struct {
	SecretId   string `json:"secret_id"`
	SecretKey  string `json:"secret_key"`
	BusinessId string `json:"business_id"`
}

type YidunInterceptor struct {
	Conf       *YidunInterceptorConf
	Conditions []*Condition

	textCheckClient *v5.TextClient
}

func (inter *YidunInterceptor) GetConditions() []*Condition {
	return inter.Conditions
}

func (inter *YidunInterceptor) CheckMsgInterceptor(ctx context.Context, senderId, receiverId string, channelType pbobjs.ChannelType, msg *pbobjs.UpMsg) (InterceptorResult, int64) {
	if msg.MsgType == msgdefines.InnerMsgType_Text {
		textMsg := &msgdefines.TextMsg{}
		err := tools.JsonUnMarshal(msg.MsgContent, textMsg)
		if err == nil {
			textCheckClient := inter.getTextCheckClient()
			if textCheckClient != nil {
				request := single.NewTextCheckRequest(inter.Conf.BusinessId)
				request.SetDataID(tools.GenerateUUIDShort11())
				request.SetContent(textMsg.Content)
				request.SetProtocol(http.ProtocolEnumHTTPS)
				response, err := textCheckClient.SyncCheckText(request)
				if err == nil {
					if response != nil {
						fmt.Println(tools.ToJson(response.Result))
						if response.GetCode() == 200 {
							if response.Result != nil && response.Result.Antispam != nil {
								suggestion := response.Result.Antispam.Suggestion
								if *suggestion == 0 || *suggestion == 1 {
									return InterceptorResult_Pass, 0
								} else {
									return InterceptorResult_Reject, 0
								}
							}
						} else {
							fmt.Println("error code: ", response.GetCode())
							fmt.Println("error msg: ", response.GetMsg())
						}
					}
				} else {
					fmt.Println(err)
				}
			}
		} else {
			fmt.Println(err)
		}
	}
	return InterceptorResult_Pass, 0
}

func (inter *YidunInterceptor) getTextCheckClient() *v5.TextClient {
	if inter.textCheckClient == nil {
		inter.textCheckClient = v5.NewTextClientWithAccessKey(inter.Conf.SecretId, inter.Conf.SecretKey)
	}
	return inter.textCheckClient
}

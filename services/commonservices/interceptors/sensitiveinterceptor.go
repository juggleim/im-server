package interceptors

import (
	"context"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/sensitivemanager/sensitivecall"
)

type SensitiveInterceptor struct {
}

func (inter *SensitiveInterceptor) GetConditions() []*Condition {
	return []*Condition{
		{
			ChannelTypeChecker: CreateMatcher("*"),
			MsgTypeChecker:     CreateMatcher("jg:text"),
			SenderIdChecker:    CreateMatcher("*"),
			ReceiverIdChecker:  CreateMatcher("*"),
		},
	}
}

func (inter *SensitiveInterceptor) CheckMsgInterceptor(ctx context.Context, senderId, receiverId string, channelType pbobjs.ChannelType, msg *pbobjs.UpMsg) bool {
	txtMsg := make(map[string]interface{})
	err := tools.JsonUnMarshal(msg.MsgContent, txtMsg)
	if err != nil {
		return false
	}
	contentVal, ok := txtMsg["content"]
	if !ok {
		return false
	}
	content, ok := contentVal.(string)
	if !ok || content == "" {
		return false
	}
	filterResp, code, err := sensitivecall.FilterCall(ctx, content)
	if err != nil || code != errs.IMErrorCode_SUCCESS {
		return false
	}
	if filterResp.HandlerType == pbobjs.SensitiveHandlerType_deny {
		return true
	} else if filterResp.HandlerType == pbobjs.SensitiveHandlerType_replace {
		txtMsg["content"] = filterResp.FilteredText
		bs, _ := tools.JsonMarshal(txtMsg)
		msg.MsgContent = bs
		return false
	}
	return false
}

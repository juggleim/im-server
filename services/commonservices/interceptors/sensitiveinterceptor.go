package interceptors

import (
	"context"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/commonservices/msgdefines"
	"im-server/services/sensitivemanager/sensitivecall"
)

type SensitiveInterceptor struct {
}

func (inter *SensitiveInterceptor) GetConditions() []*Condition {
	return []*Condition{
		{
			ChannelTypeChecker: CreateMatcher("*"),
			MsgTypeChecker:     CreateMatcher(msgdefines.InnerMsgType_Text),
			SenderIdChecker:    CreateMatcher("*"),
			ReceiverIdChecker:  CreateMatcher("*"),
		},
	}
}

func (inter *SensitiveInterceptor) CheckMsgInterceptor(ctx context.Context, senderId, receiverId string, channelType pbobjs.ChannelType, msg *pbobjs.UpMsg) (InterceptorResult, int64) {
	txtMsg := make(map[string]interface{})
	err := tools.JsonUnMarshal(msg.MsgContent, txtMsg)
	if err != nil {
		return InterceptorResult_Pass, 0
	}
	contentVal, ok := txtMsg["content"]
	if !ok {
		return InterceptorResult_Pass, 0
	}
	content, ok := contentVal.(string)
	if !ok || content == "" {
		return InterceptorResult_Pass, 0
	}
	filterResp, code, err := sensitivecall.FilterCall(ctx, content)
	if err != nil || code != errs.IMErrorCode_SUCCESS {
		return InterceptorResult_Pass, 0
	}
	if filterResp.HandlerType == pbobjs.SensitiveHandlerType_deny {
		return InterceptorResult_Reject, 0
	} else if filterResp.HandlerType == pbobjs.SensitiveHandlerType_replace {
		txtMsg["content"] = filterResp.FilteredText
		bs, _ := tools.JsonMarshal(txtMsg)
		msg.MsgContent = bs
		return InterceptorResult_Replace, 0
	}
	return InterceptorResult_Pass, 0
}

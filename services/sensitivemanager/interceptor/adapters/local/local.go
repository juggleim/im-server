package local

import (
	"context"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/sensitivemanager/sensitivecall"
)

type LocalInterceptor struct{}

func (i *LocalInterceptor) CheckMsgInterceptor(ctx context.Context, upMsg *pbobjs.UpMsg) (intercept bool, err error) {
	return CheckSensitive(ctx, upMsg), nil
}

func CheckSensitive(ctx context.Context, msg *pbobjs.UpMsg) bool {
	if msg.MsgType == "jg:text" {
		txtMsg := make(map[string]interface{})
		err := tools.JsonUnMarshal(msg.MsgContent, &txtMsg)
		if err != nil {
			return true
		}
		contentVal, ok := txtMsg["content"]
		if !ok {
			return true
		}
		content, ok := contentVal.(string)
		if !ok {
			return true
		}
		if content == "" {
			return true
		}
		filterResp, code, err := sensitivecall.FilterCall(ctx, content)
		if err != nil {
			return false
		}
		if code != errs.IMErrorCode_SUCCESS {
			return false
		}
		if filterResp.HandlerType == pbobjs.SensitiveHandlerType_deny {
			return true
		}
		if filterResp.HandlerType == pbobjs.SensitiveHandlerType_replace {
			txtMsg["content"] = filterResp.FilteredText
			bs, _ := tools.JsonMarshal(txtMsg)
			msg.MsgContent = bs
			return false
		}
	}
	return false
}

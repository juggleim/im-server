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
		txtMsg := &struct {
			Content string `json:"content"`
		}{}
		err := tools.JsonUnMarshal(msg.MsgContent, txtMsg)
		if err != nil {
			return true
		}
		filterResp, code, err := sensitivecall.FilterCall(ctx, txtMsg.Content)
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
			txtMsg.Content = filterResp.FilteredText
			bs, _ := tools.JsonMarshal(txtMsg)
			msg.MsgContent = bs
			return false
		}
	}
	return false
}

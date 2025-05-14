package interceptors

import (
	"context"
	"im-server/commons/pbdefines/pbobjs"
)

type InterceptorResult int32

var (
	InterceptorResult_Pass    InterceptorResult = 0
	InterceptorResult_Reject  InterceptorResult = 1
	InterceptorResult_Replace InterceptorResult = 2
	InterceptorResult_Silent  InterceptorResult = 3
)

type IInterceptor interface {
	GetConditions() []*Condition
	CheckMsgInterceptor(ctx context.Context, senderId, receiverId string, channelType pbobjs.ChannelType, msg *pbobjs.UpMsg) (InterceptorResult, int64)
}

type MsgInterceptor struct {
	Interceptor IInterceptor
}

func (inter *MsgInterceptor) CheckMsgInterceptor(ctx context.Context, senderId, receiverId string, channelType pbobjs.ChannelType, msg *pbobjs.UpMsg) (InterceptorResult, int64) {
	if inter.Interceptor == nil {
		return InterceptorResult_Pass, 0
	}
	if !ConditionMatchs(inter.Interceptor.GetConditions(), senderId, receiverId, channelType, msg.MsgType, msg.MsgContent) {
		return InterceptorResult_Pass, 0
	}
	return inter.Interceptor.CheckMsgInterceptor(ctx, senderId, receiverId, channelType, msg)
}

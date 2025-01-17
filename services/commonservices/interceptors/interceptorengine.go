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
)

type IInterceptor interface {
	GetConditions() []*Condition
	CheckMsgInterceptor(ctx context.Context, senderId, receiverId string, channelType pbobjs.ChannelType, msg *pbobjs.UpMsg) InterceptorResult
}

type MsgInterceptor struct {
	Interceptor IInterceptor
}

func (inter *MsgInterceptor) CheckMsgInterceptor(ctx context.Context, senderId, receiverId string, channelType pbobjs.ChannelType, msg *pbobjs.UpMsg) InterceptorResult {
	if inter.Interceptor == nil {
		return InterceptorResult_Pass
	}
	if !ConditionMatchs(inter.Interceptor.GetConditions(), senderId, receiverId, channelType, msg.MsgType, msg.MsgContent) {
		return InterceptorResult_Pass
	}
	return inter.Interceptor.CheckMsgInterceptor(ctx, senderId, receiverId, channelType, msg)
}

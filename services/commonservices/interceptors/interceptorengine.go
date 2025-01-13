package interceptors

import (
	"context"
	"fmt"
	"im-server/commons/pbdefines/pbobjs"
)

type IInterceptor interface {
	GetConditions() []*Condition
	CheckMsgInterceptor(ctx context.Context, senderId, receiverId string, channelType pbobjs.ChannelType, msg *pbobjs.UpMsg) bool
}

type MsgInterceptor struct {
	Interceptor IInterceptor
}

func (inter *MsgInterceptor) CheckMsgInterceptor(ctx context.Context, senderId, receiverId string, channelType pbobjs.ChannelType, msg *pbobjs.UpMsg) bool {
	fmt.Println("do check:")
	if inter.Interceptor == nil {
		return false
	}
	fmt.Println("no interceptor")
	if !ConditionMatchs(inter.Interceptor.GetConditions(), senderId, receiverId, channelType, msg.MsgType, msg.MsgContent) {
		fmt.Println("match:", false)
		return false
	}
	fmt.Println("have match")
	return inter.Interceptor.CheckMsgInterceptor(ctx, senderId, receiverId, channelType, msg)
}

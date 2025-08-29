package bases

import (
	"context"
	"im-server/commons/tools"
)

type BaseActorOption interface {
	HandleCtx(ctx context.Context) context.Context
}

type OnlySendboxOption struct {
}

func (opt *OnlySendboxOption) HandleCtx(ctx context.Context) context.Context {
	retCtx := setCtxValue(ctx, CtxKey_OnlySendbox, true)
	return retCtx
}

type NoNotifySenderOption struct{}

func (opt *NoNotifySenderOption) HandleCtx(ctx context.Context) context.Context {
	retCtx := setCtxValue(ctx, CtxKey_NoSendbox, true)
	return retCtx
}

type WithMsgIdOption struct {
	MsgId string
}

func (opt *WithMsgIdOption) HandleCtx(ctx context.Context) context.Context {
	retCtx := ctx
	if opt.MsgId != "" {
		retCtx = setCtxValue(retCtx, CtxKey_MsgId, opt.MsgId)
	}
	return retCtx
}

type TargetIdsOption struct {
	TargetIds []string
}

func (opt *TargetIdsOption) HandleCtx(ctx context.Context) context.Context {
	retCtx := ctx
	if len(opt.TargetIds) > 0 {
		retCtx = setCtxValue(retCtx, CtxKey_TargetIds, opt.TargetIds)
	}
	return retCtx
}

type ExtsOption struct {
	Exts map[string]string
}

func (opt *ExtsOption) HandleCtx(ctx context.Context) context.Context {
	retCtx := ctx
	if len(opt.Exts) > 0 {
		exts := GetExtsFromCtx(ctx)
		for k, v := range opt.Exts {
			exts[k] = v
		}
		retCtx = setCtxValue(retCtx, CtxKey_Exts, exts)
	}
	return retCtx
}

type MarkFromApiOption struct{}

func (opt *MarkFromApiOption) HandleCtx(ctx context.Context) context.Context {
	retCtx := setCtxValue(ctx, CtxKey_IsFromApi, true)
	return retCtx
}

type ReGenerateSessionOption struct{}

func (opt *ReGenerateSessionOption) HandleCtx(ctx context.Context) context.Context {
	retCtx := setCtxValue(ctx, CtxKey_Session, tools.GenerateUUIDShort11())
	return retCtx
}

type WithDelMsgOption struct {
	MsgId string
}

func (opt *WithDelMsgOption) HandleCtx(ctx context.Context) context.Context {
	retCtx := setCtxValue(ctx, CtxKey_DelMsgId, opt.MsgId)
	return retCtx
}

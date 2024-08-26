package bases

import "context"

type BaseActorOption interface {
	HandleCtx(ctx context.Context) context.Context
}

type OnlySendboxOption struct {
}

func (opt *OnlySendboxOption) HandleCtx(ctx context.Context) context.Context {
	retCtx := setCtxValue(ctx, CtxKey_OnlySendbox, true)
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

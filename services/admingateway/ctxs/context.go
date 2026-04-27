package ctxs

import (
	"context"
	"im-server/commons/bases"

	"github.com/gin-gonic/gin"
)

const (
	CtxKey_Account  bases.CtxKey = "CtxKey_Account"
	CtxKey_RoleType bases.CtxKey = "CtxKey_RoleType"
)

func ToCtx(ginCtx *gin.Context) context.Context {
	rpcCtx := context.Background()
	appkey := ginCtx.GetString(string(bases.CtxKey_AppKey))
	if appkey != "" {
		rpcCtx = context.WithValue(rpcCtx, bases.CtxKey_AppKey, appkey)
	}
	session := ginCtx.GetString(string(bases.CtxKey_Session))
	if session != "" {
		rpcCtx = context.WithValue(rpcCtx, bases.CtxKey_Session, session)
	}
	currentUserId := ginCtx.GetString(string(bases.CtxKey_RequesterId))
	if currentUserId != "" {
		rpcCtx = context.WithValue(rpcCtx, bases.CtxKey_RequesterId, currentUserId)
	}
	account := ginCtx.GetString(string(CtxKey_Account))
	if account != "" {
		rpcCtx = context.WithValue(rpcCtx, CtxKey_Account, account)
	}
	return rpcCtx
}

func ToCtxWithRequester(ginCtx *gin.Context, requestId string) context.Context {
	rpcCtx := ToCtx(ginCtx)
	if requestId != "" {
		rpcCtx = context.WithValue(rpcCtx, bases.CtxKey_RequesterId, requestId)
	}
	return rpcCtx
}

func GetAppKeyFromCtx(ctx context.Context) string {
	if appKey, ok := ctx.Value(bases.CtxKey_AppKey).(string); ok {
		return appKey
	}
	return ""
}

func SetAppKeyToCtx(ctx context.Context, appkey string) context.Context {
	if appkey != "" {
		ctx = context.WithValue(ctx, bases.CtxKey_AppKey, appkey)
	}
	return ctx
}

func GetRequesterIdFromCtx(ctx context.Context) string {
	if requesterId, ok := ctx.Value(bases.CtxKey_RequesterId).(string); ok {
		return requesterId
	}
	return ""
}

func GetSessionFromCtx(ctx context.Context) string {
	if session, ok := ctx.Value(bases.CtxKey_Session).(string); ok {
		return session
	}
	return ""
}

func GetAccountFromCtx(ctx context.Context) string {
	if requesterId, ok := ctx.Value(CtxKey_Account).(string); ok {
		return requesterId
	}
	return ""
}

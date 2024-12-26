package services

import (
	"context"
	"im-server/commons/bases"

	"github.com/gin-gonic/gin"
)

const (
	CtxKey_AppKey  string = "CtxKey_AppKey"
	CtxKey_Session string = "CtxKey_Session"
)

func ToRpcCtx(ginCtx *gin.Context, requestId string) context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, bases.CtxKey_AppKey, GetCtxString(ginCtx, CtxKey_AppKey))
	ctx = context.WithValue(ctx, bases.CtxKey_Session, GetCtxString(ginCtx, CtxKey_Session))
	ctx = context.WithValue(ctx, bases.CtxKey_IsFromApi, true)
	if requestId != "" {
		ctx = context.WithValue(ctx, bases.CtxKey_RequesterId, requestId)
	}
	return ctx
}

func GetCtxString(ctx *gin.Context, key string) string {
	val, exist := ctx.Get(key)
	if exist {
		return val.(string)
	} else {
		return ""
	}
}

func GetAppkeyFromCtx(ctx *gin.Context) string {
	val, exist := ctx.Get(CtxKey_AppKey)
	if exist {
		return val.(string)
	} else {
		return ""
	}
}

package services

import (
	"context"
	"im-server/commons/bases"
	"im-server/services/commonservices"

	"github.com/gin-gonic/gin"
)

// ToRpcCtx snapshots request-scoped values into a goroutine-safe context.
// Call it before the Gin handler returns; *gin.Context may be reused afterwards.
func ToRpcCtx(ginCtx *gin.Context, requestId string) context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, bases.CtxKey_AppKey, GetCtxString(ginCtx, string(bases.CtxKey_AppKey)))
	ctx = context.WithValue(ctx, bases.CtxKey_Session, GetCtxString(ginCtx, string(bases.CtxKey_Session)))
	ctx = context.WithValue(ctx, bases.CtxKey_IsFromApi, true)
	ctx = context.WithValue(ctx, bases.CtxKey_Platform, string(commonservices.Platform_Server))
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
	val, exist := ctx.Get(string(bases.CtxKey_AppKey))
	if exist {
		return val.(string)
	} else {
		return ""
	}
}

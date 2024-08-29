package imcontext

import (
	"context"
	"im-server/commons/bases"
	"sync"

	"golang.org/x/time/rate"
)

func SetContextAttr(ctx WsHandleContext, key string, value interface{}) {
	if ctx.Attachment() == nil {
		attMap := &sync.Map{}
		ctx.SetAttachment(attMap)
	}
	attMap := ctx.Attachment().(*sync.Map)
	attMap.Store(key, value)
	ctx.SetAttachment(attMap)
}

func GetContextAttr(ctx WsHandleContext, key string) interface{} {
	if ctx.Attachment() != nil {
		attMap := ctx.Attachment().(*sync.Map)
		val, ok := attMap.Load(key)
		if ok {
			return val
		}
	}
	return nil
}
func GetContextAttrString(ctx WsHandleContext, key string) string {
	ret := GetContextAttr(ctx, key)
	if ret != nil {
		str, ok := ret.(string)
		if ok {
			return str
		}
	}
	return ""
}

func GetConnSession(ctx WsHandleContext) string {
	return GetContextAttrString(ctx, StateKey_ConnectSession)
}

func GetDeviceId(ctx WsHandleContext) string {
	return GetContextAttrString(ctx, StateKey_DeviceID)
}

func GetInstanceId(ctx WsHandleContext) string {
	return GetContextAttrString(ctx, StateKey_InstanceId)
}

func GetPlatform(ctx WsHandleContext) string {
	return GetContextAttrString(ctx, StateKey_Platform)
}

func CheckConnected(ctx WsHandleContext) bool {
	str := GetContextAttrString(ctx, StateKey_Connected)
	return str == "1"
}

func GetLimiter(ctx WsHandleContext) *rate.Limiter {
	limiterObj := GetContextAttr(ctx, StateKey_Limiter)
	if limiterObj != nil {
		return limiterObj.(*rate.Limiter)
	}
	return nil
}

func GetCtxLocker(ctx WsHandleContext) *sync.Mutex {
	obj := GetContextAttr(ctx, StateKey_CtxLocker)
	if obj == nil {
		lock := &sync.Mutex{}
		SetContextAttr(ctx, StateKey_CtxLocker, lock)
		return lock
	} else {
		return obj.(*sync.Mutex)
	}
}
func GetServerIndexAfterIncrease(ctx WsHandleContext) uint16 {
	lock := GetCtxLocker(ctx)
	lock.Lock()
	defer lock.Unlock()
	var index uint16 = 0
	indexObj := GetContextAttr(ctx, StateKey_ServerMsgIndex)
	if indexObj != nil {
		index = indexObj.(uint16)
	}
	index = index + 1
	SetContextAttr(ctx, StateKey_ServerMsgIndex, index)
	return index
}

func PutServerPubCallback(ctx WsHandleContext, index int32, callback func()) {
	lock := GetCtxLocker(ctx)
	lock.Lock()
	defer lock.Unlock()
	obj := GetContextAttr(ctx, StateKey_ServerPubCallbackMap)
	var callbackMap *sync.Map
	if obj == nil {
		callbackMap = &sync.Map{}
		SetContextAttr(ctx, StateKey_ServerPubCallbackMap, callbackMap)
	} else {
		callbackMap = obj.(*sync.Map)
	}
	callbackMap.Store(index, callback)
}

func GetAndDeleteServerPubCallback(ctx WsHandleContext, index int32) func() {
	obj := GetContextAttr(ctx, StateKey_ServerPubCallbackMap)
	if obj != nil {
		callbackMap := obj.(*sync.Map)
		callbackObj, ok := callbackMap.LoadAndDelete(index)
		if ok {
			callback := callbackObj.(func())
			return callback
		}
	}
	return nil
}

func RemoveServerPubCallback(ctx WsHandleContext, index int32) {
	obj := GetContextAttr(ctx, StateKey_ServerPubCallbackMap)
	if obj != nil {
		callbackMap := obj.(*sync.Map)
		callbackMap.Delete(index)
	}
}

func PutQueryAckCallback(ctx WsHandleContext, index int32, callback func()) {
	lock := GetCtxLocker(ctx)
	lock.Lock()
	defer lock.Unlock()
	obj := GetContextAttr(ctx, StateKey_QueryConfirmMap)
	var callbackMap *sync.Map
	if obj == nil {
		callbackMap = &sync.Map{}
		SetContextAttr(ctx, StateKey_QueryConfirmMap, callbackMap)
	} else {
		callbackMap = obj.(*sync.Map)
	}
	callbackMap.Store(index, callback)
}

func GetAndDeleteQueryAckCallback(ctx WsHandleContext, index int32) func() {
	obj := GetContextAttr(ctx, StateKey_QueryConfirmMap)
	if obj != nil {
		callbackMap := obj.(*sync.Map)
		callbackObj, ok := callbackMap.LoadAndDelete(index)
		if ok {
			callback := callbackObj.(func())
			return callback
		}
	}
	return nil
}

func RemoveQueryAckCallback(ctx WsHandleContext, index int32) {
	obj := GetContextAttr(ctx, StateKey_QueryConfirmMap)
	if obj != nil {
		callbackMap := obj.(*sync.Map)
		callbackMap.Delete(index)
	}
}

func SendCtxFromNettyCtx(inboundCtx WsHandleContext) context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, bases.CtxKey_AppKey, GetContextAttrString(inboundCtx, StateKey_Appkey))
	ctx = context.WithValue(ctx, bases.CtxKey_Session, GetConnSession(inboundCtx))
	ctx = context.WithValue(ctx, bases.CtxKey_RequesterId, 0)
	ctx = context.WithValue(ctx, bases.CtxKey_Qos, 0)
	return ctx
}

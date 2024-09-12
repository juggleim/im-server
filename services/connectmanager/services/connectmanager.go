package services

import (
	"errors"
	"hash/crc32"
	"strings"
	"sync"
	"time"

	"im-server/commons/errs"
	"im-server/commons/logs"
	"im-server/services/commonservices"
	"im-server/services/connectmanager/server/codec"
	"im-server/services/connectmanager/server/imcontext"
)

var OnlineUserConnectMap sync.Map    //map[useridentifier]map[session]netty.HandlerContext
var OnlineSessionConnectMap sync.Map // map[session]netty.HandlerContext
var lockArray [512]*sync.Mutex

func init() {
	for i := 0; i < 512; i++ {
		lockArray[i] = &sync.Mutex{}
	}
}
func GetConnectCtxBySession(session string) imcontext.WsHandleContext {
	if obj, ok := OnlineSessionConnectMap.Load(session); ok {
		ctx := obj.(imcontext.WsHandleContext)
		return ctx
	}
	return nil
}
func GetConnectCtxByUser(appkey, userid string) map[string]imcontext.WsHandleContext {
	identifier := getUserIdentifier(appkey, userid)
	if ctxMapObj, ok := OnlineUserConnectMap.Load(identifier); ok {
		ctxMap := ctxMapObj.(map[string]imcontext.WsHandleContext)
		return ctxMap
	}
	return map[string]imcontext.WsHandleContext{}
}
func PutInContextCache(ctx imcontext.WsHandleContext) {
	session := imcontext.GetConnSession(ctx)
	if session != "" {
		OnlineSessionConnectMap.Store(session, ctx)

		appkey := imcontext.GetContextAttrString(ctx, imcontext.StateKey_Appkey)
		userid := imcontext.GetContextAttrString(ctx, imcontext.StateKey_UserID)
		identifier := getUserIdentifier(appkey, userid)

		lock := GetLock(identifier)
		lock.Lock()
		defer lock.Unlock()
		var userSessionMap map[string]imcontext.WsHandleContext
		if tmpUserSessionMap, ok := OnlineUserConnectMap.Load(identifier); ok {
			userSessionMap = tmpUserSessionMap.(map[string]imcontext.WsHandleContext)
		} else {
			userSessionMap = map[string]imcontext.WsHandleContext{}
			OnlineUserConnectMap.Store(identifier, userSessionMap)
		}
		userSessionMap[session] = ctx

		appinfo, exist := commonservices.GetAppInfo(appkey)
		if exist && appinfo != nil && appinfo.KickMode == 0 {
			//check other device and kick off
			currentPlatform := imcontext.GetContextAttrString(ctx, imcontext.StateKey_Platform)
			needRemoveSessions := []string{}
			for sessionKey, clientCtx := range userSessionMap {
				if sessionKey == session {
					continue
				}
				platform := imcontext.GetContextAttrString(clientCtx, imcontext.StateKey_Platform)
				if currentPlatform == string(commonservices.Platform_Android) || currentPlatform == string(commonservices.Platform_IOS) {
					if platform == string(commonservices.Platform_Android) || platform == string(commonservices.Platform_IOS) {
						needRemoveSessions = append(needRemoveSessions, sessionKey)
					}
				} else if currentPlatform == string(commonservices.Platform_PC) && platform == string(commonservices.Platform_PC) {
					needRemoveSessions = append(needRemoveSessions, sessionKey)
				}
			}
			//remove from cache
			for _, sessionKey := range needRemoveSessions {
				delete(userSessionMap, sessionKey)
				kickCtxObj, exist := OnlineSessionConnectMap.LoadAndDelete(sessionKey)
				if exist && kickCtxObj != nil {
					kickCtx := kickCtxObj.(imcontext.WsHandleContext)
					go func() {
						disconnectMsg := codec.NewDisconnectMessage(&codec.DisconnectMsgBody{
							Code:      int32(errs.IMErrorCode_CONNECT_KICKED_OFF),
							Timestamp: time.Now().UnixMilli(),
						})
						kickCtx.Write(disconnectMsg)
						logs.Infof("session:%s\taction:%s\tcode:%d", imcontext.GetConnSession(kickCtx), imcontext.Action_Disconnect, disconnectMsg.MsgBody.Code)
						Offline(kickCtx, errs.IMErrorCode_CONNECT_KICKED_OFF)
						time.Sleep(time.Millisecond * 50)
						kickCtx.Close(errors.New("kick off by other login"))
					}()
				}
			}
		}
	}
}

func RemoveFromContextCache(ctx imcontext.WsHandleContext) {
	session := imcontext.GetContextAttrString(ctx, imcontext.StateKey_ConnectSession)
	if session != "" {
		OnlineSessionConnectMap.Delete(session)
	}
	appkey := imcontext.GetContextAttrString(ctx, imcontext.StateKey_Appkey)
	userid := imcontext.GetContextAttrString(ctx, imcontext.StateKey_UserID)
	identifier := getUserIdentifier(appkey, userid)

	lock := GetLock(identifier)
	lock.Lock()
	defer lock.Unlock()
	if userSessionMapObj, ok := OnlineUserConnectMap.Load(identifier); ok {
		userSessionMap := userSessionMapObj.(map[string]imcontext.WsHandleContext)
		delete(userSessionMap, session)
		if len(userSessionMap) <= 0 {
			OnlineUserConnectMap.Delete(identifier)
		}
	}
}
func GetLock(identifier string) *sync.Mutex {
	v := int(crc32.ChecksumIEEE([]byte(identifier)))
	if v < 0 {
		v = -v
	}
	return lockArray[v%512]
}
func getUserIdentifier(appkey, userid string) string {
	return strings.Join([]string{appkey, userid}, "_")
}

func GetConnectCountByUser(appkey, userid string) int32 {
	identifier := getUserIdentifier(appkey, userid)
	if ctxMapObj, ok := OnlineUserConnectMap.Load(identifier); ok {
		ctxMap := ctxMapObj.(map[string]imcontext.WsHandleContext)
		return int32(len(ctxMap))
	}
	return 0
}

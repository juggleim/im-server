package services

import (
	"errors"
	"hash/crc32"
	"strings"
	"sync"
	"sync/atomic"
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
var onlineUserCount atomic.Int64
var userConnectCount atomic.Int64
var sessionConnectCount atomic.Int64
var appConnectCountMap sync.Map // map[appkey]*atomic.Int64

func init() {
	for i := 0; i < 512; i++ {
		lockArray[i] = &sync.Mutex{}
	}
	commonservices.RegisterClientConnectMetricsProvider(GetOnlineConnectMetrics)
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
		storeOnlineSession(session, ctx)

		appkey := imcontext.GetContextAttrString(ctx, imcontext.StateKey_Appkey)
		userid := imcontext.GetContextAttrString(ctx, imcontext.StateKey_UserID)
		identifier := getUserIdentifier(appkey, userid)
		deviceId := imcontext.GetDeviceId(ctx)

		lock := GetLock(identifier)
		lock.Lock()
		defer lock.Unlock()
		userSessionMap := addUserSessionLocked(identifier, appkey, session, ctx)

		appinfo, exist := commonservices.GetAppInfo(appkey)
		if exist && appinfo != nil && appinfo.KickMode == 0 {
			//check other device and kick off
			currentPlatform := imcontext.GetContextAttrString(ctx, imcontext.StateKey_Platform)
			needRemoveSessions := []string{}
			for sessionKey, clientCtx := range userSessionMap {
				if sessionKey == session {
					continue
				}
				did := imcontext.GetDeviceId(clientCtx)
				if did == deviceId {
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
				removeUserSessionLocked(identifier, appkey, sessionKey, userSessionMap)
				kickCtxObj, exist := loadAndDeleteOnlineSession(sessionKey)
				if exist && kickCtxObj != nil {
					kickCtx := kickCtxObj.(imcontext.WsHandleContext)
					go func() {
						code := errs.IMErrorCode_CONNECT_KICKED_OFF
						kickDeviceId := imcontext.GetDeviceId(kickCtx)
						if deviceId == kickDeviceId {
							code = errs.IMErrorCode_CONNECT_KICKED_BY_SELF
						}
						disconnectMsg := codec.NewDisconnectMessage(&codec.DisconnectMsgBody{
							Code:      int32(code),
							Timestamp: time.Now().UnixMilli(),
						})
						kickCtx.Write(disconnectMsg)
						logs.Infof("session:%s\taction:%s\tcode:%d", imcontext.GetConnSession(kickCtx), imcontext.Action_Disconnect, disconnectMsg.MsgBody.Code)
						RemoveFromContextCache(kickCtx)
						Offline(kickCtx, code)
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
		deleteOnlineSession(session)
	}
	appkey := imcontext.GetContextAttrString(ctx, imcontext.StateKey_Appkey)
	userid := imcontext.GetContextAttrString(ctx, imcontext.StateKey_UserID)
	identifier := getUserIdentifier(appkey, userid)

	lock := GetLock(identifier)
	lock.Lock()
	defer lock.Unlock()
	if userSessionMapObj, ok := OnlineUserConnectMap.Load(identifier); ok {
		userSessionMap := userSessionMapObj.(map[string]imcontext.WsHandleContext)
		removeUserSessionLocked(identifier, appkey, session, userSessionMap)
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

func GetOnlineConnectMetrics() commonservices.ClientConnectMetrics {
	return commonservices.ClientConnectMetrics{
		OnlineUserCount:     onlineUserCount.Load(),
		UserConnectCount:    userConnectCount.Load(),
		SessionConnectCount: sessionConnectCount.Load(),
	}
}

func storeOnlineSession(session string, ctx imcontext.WsHandleContext) bool {
	if session == "" {
		return false
	}
	if _, loaded := OnlineSessionConnectMap.LoadOrStore(session, ctx); loaded {
		OnlineSessionConnectMap.Store(session, ctx)
		return false
	}
	sessionConnectCount.Add(1)
	return true
}

func loadAndDeleteOnlineSession(session string) (interface{}, bool) {
	if session == "" {
		return nil, false
	}
	obj, exist := OnlineSessionConnectMap.LoadAndDelete(session)
	if exist {
		sessionConnectCount.Add(-1)
	}
	return obj, exist
}

func deleteOnlineSession(session string) bool {
	_, exist := loadAndDeleteOnlineSession(session)
	return exist
}

func addUserSessionLocked(identifier, appkey, session string, ctx imcontext.WsHandleContext) map[string]imcontext.WsHandleContext {
	var userSessionMap map[string]imcontext.WsHandleContext
	if tmpUserSessionMap, ok := OnlineUserConnectMap.Load(identifier); ok {
		userSessionMap = tmpUserSessionMap.(map[string]imcontext.WsHandleContext)
	} else {
		userSessionMap = map[string]imcontext.WsHandleContext{}
		OnlineUserConnectMap.Store(identifier, userSessionMap)
		onlineUserCount.Add(1)
	}
	if _, exist := userSessionMap[session]; !exist {
		userConnectCount.Add(1)
		incrAppConnectCount(appkey, 1)
	}
	userSessionMap[session] = ctx
	return userSessionMap
}

func removeUserSessionLocked(identifier, appkey, session string, userSessionMap map[string]imcontext.WsHandleContext) bool {
	if session == "" {
		return false
	}
	if _, exist := userSessionMap[session]; !exist {
		return false
	}
	delete(userSessionMap, session)
	userConnectCount.Add(-1)
	incrAppConnectCount(appkey, -1)
	if len(userSessionMap) <= 0 {
		if _, loaded := OnlineUserConnectMap.LoadAndDelete(identifier); loaded {
			onlineUserCount.Add(-1)
		}
	}
	return true
}

func incrAppConnectCount(appkey string, delta int64) {
	if appkey == "" || delta == 0 {
		return
	}
	counterObj, _ := appConnectCountMap.LoadOrStore(appkey, &atomic.Int64{})
	counter := counterObj.(*atomic.Int64)
	if counter.Add(delta) <= 0 {
		appConnectCountMap.Delete(appkey)
	}
}

func foreachAppConnectCount(f func(appkey string, count int64)) {
	appConnectCountMap.Range(func(key, value any) bool {
		appkey, ok := key.(string)
		if !ok {
			return true
		}
		counter, ok := value.(*atomic.Int64)
		if !ok {
			return true
		}
		count := counter.Load()
		if count > 0 {
			f(appkey, count)
		}
		return true
	})
}

func resetOnlineConnectStateForTest() {
	OnlineUserConnectMap = sync.Map{}
	OnlineSessionConnectMap = sync.Map{}
	appConnectCountMap = sync.Map{}
	onlineUserCount.Store(0)
	userConnectCount.Store(0)
	sessionConnectCount.Store(0)
}

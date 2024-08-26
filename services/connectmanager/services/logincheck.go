package services

import (
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices"
	"im-server/services/commonservices/tokens"
	"im-server/services/connectmanager/dbs"
	"im-server/services/connectmanager/server/codec"
	"im-server/services/connectmanager/server/imcontext"
	"strings"
	"time"
)

var supportPlatforms map[string]bool

func init() {
	supportPlatforms = make(map[string]bool)
	supportPlatforms[string(commonservices.Platform_Android)] = true
	supportPlatforms[string(commonservices.Platform_IOS)] = true
	supportPlatforms[string(commonservices.Platform_Web)] = true
	supportPlatforms[string(commonservices.Platform_PC)] = true
}

func CheckLogin(ctx imcontext.WsHandleContext, msg *codec.ConnectMsgBody) (int32, string) {
	appkey := msg.Appkey
	tokenStr := msg.Token
	//check platform
	if _, exist := supportPlatforms[msg.Platform]; !exist {
		return int32(errs.IMErrorCode_CONNECT_UNSUPPROTEDPLATFORM), ""
	}
	//check token
	tokenWrap, err := tokens.ParseTokenString(tokenStr)
	if err != nil {
		return int32(errs.IMErrorCode_CONNECT_TOKEN_AUTHFAIL), ""
	}
	appinfo, exist := commonservices.GetAppInfo(appkey)
	if !exist || appinfo == nil {
		return int32(errs.IMErrorCode_CONNECT_APP_NOT_EXISTED), ""
	}
	token, err := tokens.ParseToken(tokenWrap, []byte(appinfo.AppSecureKey))
	if err != nil || token.UserId == "" {
		return int32(errs.IMErrorCode_CONNECT_TOKEN_AUTHFAIL), ""
	}

	if appinfo.TokenEffectiveMinute > 0 && (token.TokenTime+int64(appinfo.TokenEffectiveMinute)*60*1000) < time.Now().UnixMilli() {
		return int32(errs.IMErrorCode_CONNECT_TOKEN_EXPIRED), ""
	}
	imcontext.SetContextAttr(ctx, imcontext.StateKey_UserID, token.UserId)
	imcontext.SetContextAttr(ctx, imcontext.StateKey_DeviceID, msg.DeviceId)

	//check ban user
	banUser, exist := GetBanUserFromCache(appkey, token.UserId)
	if exist {
		currentTime := time.Now().UnixMilli()
		if ban, exist := banUser.Items[string(dbs.UserBanScopeDefault)]; exist {
			if ban.BanType == pbobjs.BanType_Permanent || (ban.BanType == pbobjs.BanType_Temporary && ban.EndTime > currentTime) {
				return int32(errs.IMErrorCode_CONNECT_USER_BLOCK), ban.Ext
			}
		}
		if ban, exist := banUser.Items[string(dbs.UserBanScopePlatform)]; exist && strings.Contains(ban.ScopeValue, msg.Platform) {
			if ban.BanType == pbobjs.BanType_Permanent || (ban.BanType == pbobjs.BanType_Temporary && ban.EndTime > currentTime) {
				return int32(errs.IMErrorCode_CONNECT_USER_BLOCK), ban.Ext
			}
		}
		if ban, exist := banUser.Items[string(dbs.UserBanScopeDevice)]; exist && strings.Contains(ban.ScopeValue, msg.DeviceId) {
			if ban.BanType == pbobjs.BanType_Permanent || (ban.BanType == pbobjs.BanType_Temporary && ban.EndTime > currentTime) {
				return int32(errs.IMErrorCode_CONNECT_USER_BLOCK), ban.Ext
			}
		}
		if ban, exist := banUser.Items[string(dbs.UserBanScopeIp)]; exist && strings.Contains(ban.ScopeValue, msg.ClientIp) {
			if ban.BanType == pbobjs.BanType_Permanent || (ban.BanType == pbobjs.BanType_Temporary && ban.EndTime > currentTime) {
				return int32(errs.IMErrorCode_CONNECT_USER_BLOCK), ban.Ext
			}
		}
	}

	return int32(errs.IMErrorCode_SUCCESS), ""
}

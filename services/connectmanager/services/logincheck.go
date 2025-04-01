package services

import (
	"im-server/commons/errs"
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
	supportPlatforms[string(commonservices.Platform_Harmony)] = true
}

func CheckLogin(ctx imcontext.WsHandleContext, msg *codec.ConnectMsgBody) (int32, string) {
	appkey := msg.Appkey
	tokenStr := msg.Token
	//check platform
	if _, exist := supportPlatforms[msg.Platform]; !exist {
		return int32(errs.IMErrorCode_CONNECT_UNSUPPROTEDPLATFORM), ""
	}
	//check security domain
	if msg.Platform == string(commonservices.Platform_Web) {
		referer := imcontext.GetReferer(ctx)
		appInfo, exist := commonservices.GetAppInfo(appkey)
		if exist && appInfo != nil && appInfo.SecurityDomainsObj != nil && len(appInfo.SecurityDomainsObj.Domains) > 0 {
			isContains := false
			for _, domain := range appInfo.SecurityDomainsObj.Domains {
				if domain == referer {
					isContains = true
					break
				}
			}
			if !isContains {
				return int32(errs.IMErrorCode_CONNECT_UNSECURITYDOMAIN), ""
			}
		}
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
			if ban.EndTime == 0 || ban.EndTime > currentTime {
				return int32(errs.IMErrorCode_CONNECT_USER_BLOCK), ban.Ext
			}
		}
		if ban, exist := banUser.Items[string(dbs.UserBanScopePlatform)]; exist && strings.Contains(ban.ScopeValue, msg.Platform) {
			if ban.EndTime == 0 || ban.EndTime > currentTime {
				return int32(errs.IMErrorCode_CONNECT_USER_BLOCK), ban.Ext
			}
		}
		if ban, exist := banUser.Items[string(dbs.UserBanScopeDevice)]; exist && strings.Contains(ban.ScopeValue, msg.DeviceId) {
			if ban.EndTime == 0 || ban.EndTime > currentTime {
				return int32(errs.IMErrorCode_CONNECT_USER_BLOCK), ban.Ext
			}
		}
		if ban, exist := banUser.Items[string(dbs.UserBanScopeIp)]; exist && strings.Contains(ban.ScopeValue, msg.ClientIp) {
			if ban.EndTime == 0 || ban.EndTime > currentTime {
				return int32(errs.IMErrorCode_CONNECT_USER_BLOCK), ban.Ext
			}
		}
	}

	return int32(errs.IMErrorCode_SUCCESS), ""
}

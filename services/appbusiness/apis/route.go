package apis

import (
	"im-server/commons/errs"
	"im-server/services/appbusiness/httputils"
	"im-server/services/commonservices"
	"im-server/services/commonservices/tokens"
	"net/http"
	"net/url"
)

func LoadAppApis(mux *http.ServeMux) {
	RouteRegiste(mux, http.MethodPost, "/jim/login", Login)
	RouteRegiste(mux, http.MethodPost, "/jim/sms/send", func(ctx *httputils.HttpContext) {
		ctx.ResponseSucc(nil)
	})
	RouteRegiste(mux, http.MethodPost, "/jim/sms_login", SmsLogin)
	RouteRegiste(mux, http.MethodPost, "/jim/file_cred", GetFileCred)

	RouteRegiste(mux, http.MethodGet, "/jim/bots/list", QryBots)

	RouteRegiste(mux, http.MethodPost, "/jim/users/update", UpdateUser)
	RouteRegiste(mux, http.MethodPost, "/jim/users/updsettings", UpdateUserSettings)
	RouteRegiste(mux, http.MethodPost, "/jim/users/search", SearchByPhone)
	RouteRegiste(mux, http.MethodGet, "/jim/users/info", QryUserInfo)

	RouteRegiste(mux, http.MethodPost, "/jim/groups/add", CreateGroup)
	RouteRegiste(mux, http.MethodPost, "/jim/groups/create", CreateGroup)
	RouteRegiste(mux, http.MethodPost, "/jim/groups/update", UpdateGroup)
	RouteRegiste(mux, http.MethodPost, "/jim/groups/dissolve", DissolveGroup)
	RouteRegiste(mux, http.MethodPost, "/jim/groups/members/add", AddGrpMembers)
	RouteRegiste(mux, http.MethodPost, "/jim/groups/apply", GroupApply)
	RouteRegiste(mux, http.MethodPost, "/jim/groups/invite", GroupInvite)
	RouteRegiste(mux, http.MethodPost, "/jim/groups/quit", QuitGroup)
	RouteRegiste(mux, http.MethodPost, "/jim/groups/members/del", DelGrpMembers)
	RouteRegiste(mux, http.MethodGet, "/jim/groups/members/list", QryGrpMembers)
	RouteRegiste(mux, http.MethodGet, "/jim/groups/info", QryGroupInfo)
	RouteRegiste(mux, http.MethodPost, "/jim/groups/setgrpannouncement", SetGrpAnnouncement)
	RouteRegiste(mux, http.MethodGet, "/jim/groups/getgrpannouncement", GetGrpAnnouncement)
	RouteRegiste(mux, http.MethodPost, "/jim/groups/setdisplayname", SetGrpDisplayName)
	//group manage
	RouteRegiste(mux, http.MethodPost, "/jim/groups/management/chgowner", ChgGroupOwner)
	RouteRegiste(mux, http.MethodPost, "/jim/groups/management/administrators/add", AddGrpAdministrator)
	RouteRegiste(mux, http.MethodPost, "/jim/groups/management/administrators/del", DelGrpAdministrator)
	RouteRegiste(mux, http.MethodGet, "/jim/groups/management/administrators/list", QryGrpAdministrators)
	RouteRegiste(mux, http.MethodPost, "/jim/groups/management/setmute", SetGroupMute)
	RouteRegiste(mux, http.MethodPost, "/jim/groups/management/setgrpverifytype", SetGrpVerifyType)
	RouteRegiste(mux, http.MethodPost, "/jim/groups/management/sethismsgvisible", SetGrpHisMsgVisible)
	RouteRegiste(mux, http.MethodGet, "/jim/groups/mygroups", QryMyGroups)
	// grp application
	RouteRegiste(mux, http.MethodGet, "/jim/groups/myapplications", QryMyGrpApplications)
	RouteRegiste(mux, http.MethodGet, "/jim/groups/mypendinginvitations", QryMyPendingGrpInvitations)
	RouteRegiste(mux, http.MethodGet, "/jim/groups/grpinvitations", QryGrpInvitations)
	RouteRegiste(mux, http.MethodGet, "/jim/groups/grppendingapplications", QryGrpPendingApplications)

	RouteRegiste(mux, http.MethodGet, "/jim/friends/list", QryFriends)
	RouteRegiste(mux, http.MethodPost, "/jim/friends/add", AddFriend)
	RouteRegiste(mux, http.MethodPost, "/jim/friends/apply", ApplyFriend)
	RouteRegiste(mux, http.MethodPost, "/jim/friends/confirm", ConfirmFriend)
	RouteRegiste(mux, http.MethodPost, "/jim/friends/del", DelFriend)
	RouteRegiste(mux, http.MethodGet, "/jim/friends/applications", FriendApplications)
	RouteRegiste(mux, http.MethodGet, "/jim/friends/myapplications", MyFriendApplications)
	RouteRegiste(mux, http.MethodGet, "/jim/friends/mypendingapplications", MyPendingFriendApplications)
}

func RouteRegiste(mux *http.ServeMux, method, path string, handler func(ctx *httputils.HttpContext)) {
	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Content-Type", "application/json")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		if r.Method != method {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		qryParams, _ := url.ParseQuery(r.URL.RawQuery)
		ctx := &httputils.HttpContext{
			Writer:      w,
			Request:     r,
			QueryParams: qryParams,
		}
		//check appkey
		appkey := r.Header.Get("appkey")
		if appkey == "" {
			ctx.ResponseErr(errs.IMErrorCode_APP_APPKEY_REQUIRED)
			return
		}
		ctx.AppKey = appkey
		appInfo, exist := commonservices.GetAppInfo(appkey)
		if !exist || appInfo == nil {
			ctx.ResponseErr(errs.IMErrorCode_APP_NOT_EXISTED)
			return
		}
		urlPath := r.URL.Path
		if urlPath != "/jim/login" && urlPath != "/jim/sms/send" && urlPath != "/jim/sms_login" {
			//current userId
			tokenStr := r.Header.Get("Authorization")
			if tokenStr == "" {
				ctx.ResponseErr(errs.IMErrorCode_APP_NOT_LOGIN)
				return
			}
			if tokenStr != "" {
				tokenWrap, err := tokens.ParseTokenString(tokenStr)
				if err != nil {
					ctx.ResponseErr(errs.IMErrorCode_APP_NOT_LOGIN)
					return
				}
				token, err := tokens.ParseToken(tokenWrap, []byte(appInfo.AppSecureKey))
				if err != nil {
					ctx.ResponseErr(errs.IMErrorCode_APP_NOT_LOGIN)
					return
				}
				ctx.CurrentUserId = token.UserId
			}
		}
		handler(ctx)
	})
}

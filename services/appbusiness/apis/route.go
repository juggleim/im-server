package apis

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/juggleim/jugglechat-server/apimodels"
	"github.com/juggleim/jugglechat-server/errs"
	"github.com/juggleim/jugglechat-server/services"
	"github.com/juggleim/jugglechat-server/utils"
)

func LoadAppApis(mux *http.ServeMux) {
	RouteRegiste(mux, http.MethodPost, "/jim/login", Login)
	RouteRegiste(mux, http.MethodGet, "/jim/login/qrcode", GenerateQrCode)
	RouteRegiste(mux, http.MethodPost, "/jim/login/qrcode/check", CheckQrCode)
	RouteRegiste(mux, http.MethodPost, "/jim/sms/send", SmsSend)
	RouteRegiste(mux, http.MethodPost, "/jim/sms_login", SmsLogin)
	RouteRegiste(mux, http.MethodPost, "/jim/sms/login", SmsLogin)
	RouteRegiste(mux, http.MethodPost, "/jim/email/send", EmailSend)
	RouteRegiste(mux, http.MethodPost, "/jim/email/login", EmailLogin)
	RouteRegiste(mux, http.MethodPost, "/jim/login/qrcode/confirm", ConfirmQrCode)
	RouteRegiste(mux, http.MethodPost, "/jim/file_cred", GetFileCred)

	RouteRegiste(mux, http.MethodGet, "/jim/bots/list", QryBots)

	RouteRegiste(mux, http.MethodPost, "/jim/assistants/answer", AssistantAnswer)
	RouteRegiste(mux, http.MethodPost, "/jim/assistants/prompts/add", PromptAdd)
	RouteRegiste(mux, http.MethodPost, "/jim/assistants/prompts/update", PromptUpdate)
	RouteRegiste(mux, http.MethodPost, "/jim/assistants/prompts/del", PromptDel)
	RouteRegiste(mux, http.MethodPost, "/jim/assistants/prompts/batchdel", PromptBatchDel)
	RouteRegiste(mux, http.MethodGet, "/jim/assistants/prompts/list", QryPrompts)

	RouteRegiste(mux, http.MethodPost, "/jim/bots/messages/listener", BotMsgListener)

	RouteRegiste(mux, http.MethodPost, "/jim/users/update", UpdateUser)
	RouteRegiste(mux, http.MethodPost, "/jim/users/updsettings", UpdateUserSettings)
	RouteRegiste(mux, http.MethodPost, "/jim/users/search", SearchByPhone)
	RouteRegiste(mux, http.MethodGet, "/jim/users/info", QryUserInfo)
	RouteRegiste(mux, http.MethodGet, "/jim/users/qrcode", QryUserQrCode)

	RouteRegiste(mux, http.MethodPost, "/jim/telegrambots/add", TelegramBotAdd)
	RouteRegiste(mux, http.MethodPost, "/jim/telegrambots/del", TelegramBotDel)
	RouteRegiste(mux, http.MethodPost, "/jim/telegrambots/batchdel", TelegramBotBatchDel)
	RouteRegiste(mux, http.MethodGet, "/jim/telegrambots/list", TelegramBotList)

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
	RouteRegiste(mux, http.MethodPost, "/jim/groups/members/check", CheckGroupMembers)
	RouteRegiste(mux, http.MethodGet, "/jim/groups/info", QryGroupInfo)
	RouteRegiste(mux, http.MethodGet, "/jim/groups/qrcode", QryGrpQrCode)
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

	RouteRegiste(mux, http.MethodGet, "/jim/friends/list", QryFriendsWithPage)
	RouteRegiste(mux, http.MethodPost, "/jim/friends/add", AddFriend)
	RouteRegiste(mux, http.MethodPost, "/jim/friends/apply", ApplyFriend)
	RouteRegiste(mux, http.MethodPost, "/jim/friends/confirm", ConfirmFriend)
	RouteRegiste(mux, http.MethodPost, "/jim/friends/del", DelFriend)
	RouteRegiste(mux, http.MethodGet, "/jim/friends/applications", FriendApplications)
	RouteRegiste(mux, http.MethodGet, "/jim/friends/myapplications", MyFriendApplications)
	RouteRegiste(mux, http.MethodGet, "/jim/friends/mypendingapplications", MyPendingFriendApplications)

	//post
	RouteRegiste(mux, http.MethodGet, "/jim/posts/list", QryPosts)
	RouteRegiste(mux, http.MethodGet, "/jim/posts/info", PostInfo)
	RouteRegiste(mux, http.MethodPost, "/jim/posts/add", PostAdd)
	// RouteRegiste(mux, http.MethodPost, "/jim/posts/update", nil)
	// RouteRegiste(mux, http.MethodPost, "/jim/posts/del", nil)
	// RouteRegiste(mux, http.MethodPost, "/jim/posts/reactions/add", nil)
	// RouteRegiste(mux, http.MethodGet, "/jim/posts/reactions/list", nil)

	RouteRegiste(mux, http.MethodGet, "/jim/posts/comments/list", QryPostComments)
	RouteRegiste(mux, http.MethodPost, "/jim/posts/comments/add", PostCommentAdd)
	// RouteRegiste(mux, http.MethodPost, "/jim/posts/comments/update", nil)
	// RouteRegiste(mux, http.MethodPost, "/jim/posts/comments/del", nil)
}

func RouteRegiste(mux *http.ServeMux, method, path string, handler func(ctx *HttpContext)) {
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
		ctx := &HttpContext{
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
		appInfo, exist := services.GetAppInfo(appkey)
		if !exist || appInfo == nil {
			ctx.ResponseErr(errs.IMErrorCode_APP_NOT_EXISTED)
			return
		}
		urlPath := r.URL.Path
		if urlPath != "/jim/login" && urlPath != "/jim/sms/send" && urlPath != "/jim/sms_login" && urlPath != "/jim/sms/login" && urlPath != "/jim/email/send" && urlPath != "/jim/email/login" && urlPath != "/jim/login/qrcode" && urlPath != "/jim/login/qrcode/check" {
			//current userId
			tokenStr := r.Header.Get("Authorization")
			if tokenStr == "" {
				ctx.ResponseErr(errs.IMErrorCode_APP_NOT_LOGIN)
				return
			}
			if tokenStr != "" {
				if strings.HasPrefix(tokenStr, "Bearer ") {
					tokenStr = tokenStr[7:]
					if !services.CheckApiKey(tokenStr, appkey, appInfo.AppSecureKey) {
						ctx.ResponseErr(errs.IMErrorCode_APP_NOT_LOGIN)
						return
					}
				} else {
					tokenWrap, err := services.ParseTokenString(tokenStr)
					if err != nil {
						ctx.ResponseErr(errs.IMErrorCode_APP_NOT_LOGIN)
						return
					}
					token, err := services.ParseToken(tokenWrap, []byte(appInfo.AppSecureKey))
					if err != nil {
						ctx.ResponseErr(errs.IMErrorCode_APP_NOT_LOGIN)
						return
					}
					ctx.CurrentUserId = token.UserId
				}
			}
		}
		handler(ctx)
	})
}

type HttpContext struct {
	Writer      http.ResponseWriter
	Request     *http.Request
	QueryParams url.Values

	AppKey        string
	CurrentUserId string
}

func (ctx *HttpContext) BindJSON(req interface{}) error {
	return Body2Obj(ctx.Request.Body, req)
}

func (ctx *HttpContext) Query(key string) string {
	if ctx.QueryParams != nil {
		return ctx.QueryParams.Get(key)
	}
	return ""
}

func (ctx *HttpContext) ResponseErr(code errs.IMErrorCode) {
	appErr := errs.GetApiErrorByCode(code)
	ctx.Writer.WriteHeader(appErr.HttpCode)
	bs, _ := utils.JsonMarshal(appErr)
	ctx.Writer.Write(bs)
}

func (ctx *HttpContext) ResponseSucc(resp interface{}) {
	connonResp := &apimodels.CommonResp{
		CommonError: apimodels.CommonError{
			ErrorMsg: "success",
		},
		Data: resp,
	}
	ctx.Writer.WriteHeader(http.StatusOK)
	bs, _ := utils.JsonMarshal(connonResp)
	ctx.Writer.Write(bs)
}

func (ctx *HttpContext) ToRpcCtx() context.Context {
	rpcCtx := context.Background()
	rpcCtx = context.WithValue(rpcCtx, services.CtxKey_AppKey, ctx.AppKey)
	rpcCtx = context.WithValue(rpcCtx, services.CtxKey_Session, fmt.Sprintf("app_%s", utils.GenerateUUIDShort11()))
	if ctx.CurrentUserId != "" {
		rpcCtx = context.WithValue(rpcCtx, services.CtxKey_RequesterId, ctx.CurrentUserId)
	}
	return rpcCtx
}

func Read2String(read io.ReadCloser) string {
	buf := bytes.NewBuffer([]byte{})
	for {
		bs := make([]byte, 1024)
		c, err := read.Read(bs)
		buf.Write(bs)
		if err != nil || c < 1024 {
			break
		}
	}
	return buf.String()
}

func Read2Bytes(read io.ReadCloser) []byte {
	bs, err := io.ReadAll(read)
	if err == nil {
		return bs
	}
	return []byte{}
}

func Body2Obj(read io.ReadCloser, obj interface{}) error {
	bs := Read2Bytes(read)
	if len(bs) <= 0 {
		return errors.New("no value")
	}
	return utils.JsonUnMarshal(bs, obj)
}

func ErrorHttpResp(ctx *HttpContext, code errs.IMErrorCode) {
	ctx.ResponseErr(code)
}

func SuccessHttpResp(ctx *HttpContext, resp interface{}) {
	ctx.ResponseSucc(resp)
}

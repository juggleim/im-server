package navigator

import (
	"errors"
	"im-server/commons/errs"
	"im-server/commons/gmicro/utils"
	"im-server/services/commonservices"
	"im-server/services/commonservices/tokens"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Result struct {
	Content string `json:"content"`
}

func LoadClientLogUploadApis(mux *http.ServeMux) {
	routeRegiste(mux, http.MethodPost, "/navigator/test", func(ctx *NavHttpContext) {
		ctx.ResponseSucc(&Result{
			Content: "test_content",
		})
	})
	routeRegiste(mux, http.MethodPost, "/navigator/upload-log-plain", UploadLogPlain)
	routeRegiste(mux, http.MethodPost, "/navigator/log-status", LogStatus)
}

const (
	CtxKey_AppKey string = "CtxKey_AppKey"
	CtxKey_UserId string = "CtxKey_UserId"
)

type NavHttpContext struct {
	Writer      http.ResponseWriter
	Request     *http.Request
	QueryParams url.Values

	AppKey        string
	CurrentUserId string
}

func (ctx *NavHttpContext) BindJSON(req interface{}) error {
	return Body2Obj(ctx.Request.Body, req)
}

func (ctx *NavHttpContext) Query(key string) string {
	if ctx.QueryParams != nil {
		return ctx.QueryParams.Get(key)
	}
	return ""
}

func (ctx *NavHttpContext) ResponseErr(code errs.IMErrorCode) {
	appErr := errs.GetApiErrorByCode(code)
	ctx.Writer.WriteHeader(appErr.HttpCode)
	bs, _ := utils.JsonMarshal(appErr)
	ctx.Writer.Write(bs)
}

type CommonResp struct {
	errs.ApiErrorMsg
	Data interface{} `json:"data,omitempty"`
}

func (ctx *NavHttpContext) ResponseSucc(resp interface{}) {
	connonResp := &CommonResp{
		ApiErrorMsg: errs.ApiErrorMsg{
			Msg: "success",
		},
		Data: resp,
	}
	ctx.Writer.WriteHeader(http.StatusOK)
	bs, _ := utils.JsonMarshal(connonResp)
	ctx.Writer.Write(bs)
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

func routeRegiste(mux *http.ServeMux, method, path string, handler func(ctx *NavHttpContext)) {
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
		ctx := &NavHttpContext{
			Writer:      w,
			Request:     r,
			QueryParams: qryParams,
		}
		urlPath := r.URL.Path
		if urlPath != "/navigator/test" {
			//check appkey
			appkey := r.Header.Get("x-appkey")
			if appkey == "" {
				ctx.ResponseErr(errs.IMErrorCode_CONNECT_APPKEY_REQUIRED)
				return
			}
			ctx.AppKey = appkey

			//token
			tokenStr := r.Header.Get("x-token")
			userId, code := parseToken(appkey, tokenStr)
			if code != errs.IMErrorCode_SUCCESS {
				ctx.ResponseErr(code)
				return
			}
			ctx.CurrentUserId = userId
		}
		handler(ctx)
	})
}

func parseToken(appkey, tokenStr string) (string, errs.IMErrorCode) {
	tokenWrap, err := tokens.ParseTokenString(tokenStr)
	if err != nil {
		return "", errs.IMErrorCode_CONNECT_TOKEN_ILLEGAL
	}
	if tokenWrap.AppKey != appkey {
		return "", errs.IMErrorCode_CONNECT_TOKEN_AUTHFAIL
	}
	appInfo, exist := commonservices.GetAppInfo(appkey)
	if !exist || appInfo == nil {
		return "", errs.IMErrorCode_CONNECT_APP_NOT_EXISTED
	}
	token, err := tokens.ParseToken(tokenWrap, []byte(appInfo.AppSecureKey))
	if err != nil {
		return "", errs.IMErrorCode_CONNECT_TOKEN_AUTHFAIL
	}
	if appInfo.TokenEffectiveMinute > 0 && (token.TokenTime+int64(appInfo.TokenEffectiveMinute)*60*1000) < time.Now().UnixMilli() {
		return "", errs.IMErrorCode_CONNECT_TOKEN_EXPIRED
	}
	return token.UserId, errs.IMErrorCode_SUCCESS
}

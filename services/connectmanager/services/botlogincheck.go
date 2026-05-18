package services

import (
	"im-server/commons/errs"
	"im-server/services/commonservices"
	"im-server/services/commonservices/tokens"
	"im-server/services/connectmanager/server/codec"
	"im-server/services/connectmanager/server/imcontext"
)

func CheckBotLogin(ctx imcontext.WsHandleContext, msg *codec.ConnectMsgBody) errs.IMErrorCode {
	appkey := msg.Appkey
	tokenStr := msg.Token
	//check token
	tokenWrap, err := tokens.ParseTokenString(tokenStr)
	if err != nil {
		return errs.IMErrorCode_CONNECT_TOKEN_AUTHFAIL
	}
	appinfo, exist := commonservices.GetAppInfo(appkey)
	if !exist || appinfo == nil {
		return errs.IMErrorCode_CONNECT_APP_NOT_EXISTED
	}
	token, err := tokens.ParseToken(tokenWrap, []byte(appinfo.AppSecureKey))
	if err != nil || token.UserId == "" {
		return errs.IMErrorCode_CONNECT_TOKEN_AUTHFAIL
	}
	imcontext.SetContextAttr(ctx, imcontext.StateKey_UserID, token.UserId)
	imcontext.SetContextAttr(ctx, imcontext.StateKey_UserType, token.UserType)
	return errs.IMErrorCode_SUCCESS
}

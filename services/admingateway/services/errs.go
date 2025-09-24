package services

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type AdminErrorCode int

var (
	AdminErrorCode_Success           AdminErrorCode = 0
	AdminErrorCode_Default           AdminErrorCode = 1000
	AdminErrorCode_AuthFail          AdminErrorCode = 1001
	AdminErrorCode_ParamError        AdminErrorCode = 1002
	AdminErrorCode_LoginFail         AdminErrorCode = 1003
	AdminErrorCode_LicenseNotIllegal AdminErrorCode = 1004
	AdminErrorCode_LicenseExpired    AdminErrorCode = 1005
	AdminErrorCode_AppHasExisted     AdminErrorCode = 1006
	AdminErrorCode_AddAppFail        AdminErrorCode = 1007
	AdminErrorCode_AppkeyNotMatch    AdminErrorCode = 1008
	AdminErrorCode_AppNotExist       AdminErrorCode = 1009
	AdminErrorCode_UpdAppFail        AdminErrorCode = 1010
	AdminErrorCode_NotSupportField   AdminErrorCode = 1011
	AdminErrorCode_AccountExisted    AdminErrorCode = 1012
	AdminErrorCode_UpdPwdFail        AdminErrorCode = 1013
	AdminErrorCode_ServerErr         AdminErrorCode = 1014
	AdminErrorCode_NoFileEngine      AdminErrorCode = 1015
	AdminErrorCode_AccountForbidden  AdminErrorCode = 1016
	AdminErrorCode_AccountNotExist   AdminErrorCode = 1017
	AdminErrorCode_NotPermission     AdminErrorCode = 1018
)

type ApiErrorMsg struct {
	HttpCode int            `json:"-"`
	Code     AdminErrorCode `json:"code"`
	Msg      string         `json:"msg"`
}

type SuccHttpResp struct {
	ApiErrorMsg
	Data interface{} `json:"data"`
}

func SuccessHttpResp(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, SuccHttpResp{
		ApiErrorMsg: ApiErrorMsg{
			Code: 0,
			Msg:  "success",
		},
		Data: data,
	})
}

func FailHttpResp(ctx *gin.Context, code AdminErrorCode, msgs ...string) {
	var msg = "fail"
	if len(msgs) > 0 {
		msg = msgs[0]
	}
	ctx.JSON(http.StatusOK, SuccHttpResp{
		ApiErrorMsg: ApiErrorMsg{
			Code: code,
			Msg:  msg,
		},
	})
}

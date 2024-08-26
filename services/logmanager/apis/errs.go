package apis

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ApiErrorCode int

var (
	ApiErrorCode_Success    ApiErrorCode = 0
	ApiErrorCode_Default    ApiErrorCode = 1000
	ApiErrorCode_AuthFail   ApiErrorCode = 1001
	ApiErrorCode_ParamError ApiErrorCode = 1002
	ApiErrorCode_ServerErr  ApiErrorCode = 1003
)

type ApiErrorMsg struct {
	HttpCode int          `json:"-"`
	Code     ApiErrorCode `json:"code"`
	Msg      string       `json:"msg"`
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

func FailHttpResp(ctx *gin.Context, code ApiErrorCode, msgs ...string) {
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

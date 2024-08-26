package tools

import (
	"net/http"

	"im-server/commons/errs"

	"github.com/gin-gonic/gin"
)

func ErrorHttpResp(ctx *gin.Context, code errs.IMErrorCode) {
	apiErr := errs.GetApiErrorByCode(code)
	ctx.JSON(int(apiErr.HttpCode), apiErr)
}

type SuccHttpResp struct {
	errs.ApiErrorMsg
	Data interface{} `json:"data"`
}

func SuccessHttpResp(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, SuccHttpResp{
		ApiErrorMsg: errs.ApiErrorMsg{
			Code: 0,
			Msg:  "success",
		},
		Data: data,
	})
}

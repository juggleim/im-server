package apis

import (
	"im-server/services/admingateway/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetEventSubConfig(ctx *gin.Context) {
	var req services.EventSubConfigReq
	if err := ctx.ShouldBindJSON(&req); err != nil || req.AppKey == "" {
		ctx.JSON(http.StatusBadRequest, &services.ApiErrorMsg{
			Code: services.AdminErrorCode_ParamError,
			Msg:  "param illegal",
		})
		return
	}
	code := services.SetEventSubConfig(&req)
	ctx.JSON(http.StatusOK, &services.ApiErrorMsg{
		Code: code,
	})
}

func GetEventSubConfig(ctx *gin.Context) {
	appkey := ctx.Query("app_key")
	if appkey == "" {
		ctx.JSON(http.StatusBadRequest, &services.ApiErrorMsg{
			Code: services.AdminErrorCode_ParamError,
			Msg:  "param illegal",
		})
		return
	}
	code, config := services.GetEventSubConfig(appkey)
	if code != services.AdminErrorCode_Success {
		ctx.JSON(http.StatusOK, &services.ApiErrorMsg{
			Code: code,
		})
	} else {
		services.SuccessHttpResp(ctx, config)
	}
}

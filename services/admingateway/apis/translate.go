package apis

import (
	"im-server/services/admingateway/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetTranslateConf(ctx *gin.Context) {
	var req services.TranslateConf
	if err := ctx.ShouldBindJSON(&req); err != nil || req.AppKey == "" || req.Conf == nil {
		ctx.JSON(http.StatusBadRequest, &services.ApiErrorMsg{
			Code: services.AdminErrorCode_ParamError,
			Msg:  "param illegal",
		})
		return
	}
	code := services.SetTranslateConf(req.AppKey, req.Conf)
	ctx.JSON(http.StatusOK, &services.ApiErrorMsg{
		Code: code,
	})
}

func GetTranslateConf(ctx *gin.Context) {
	appkey := ctx.Query("app_key")
	if appkey == "" {
		ctx.JSON(http.StatusBadRequest, &services.ApiErrorMsg{
			Code: services.AdminErrorCode_ParamError,
			Msg:  "param illegal",
		})
		return
	}
	code, conf := services.GetTranslateConf(appkey)
	if code != services.AdminErrorCode_Success {
		ctx.JSON(http.StatusOK, &services.ApiErrorMsg{
			Code: code,
		})
	} else {
		services.SuccessHttpResp(ctx, conf)
	}
}

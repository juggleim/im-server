package apis

import (
	"im-server/services/admingateway/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetRtcConf(ctx *gin.Context) {
	var req services.RtcConfReq
	if err := ctx.ShouldBindJSON(&req); err != nil || req.AppKey == "" || req.Conf == nil {
		ctx.JSON(http.StatusBadRequest, &services.ApiErrorMsg{
			Code: services.AdminErrorCode_ParamError,
			Msg:  "param illegal",
		})
		return
	}
	code := services.SetRtcConf(req.AppKey, req.Conf)
	ctx.JSON(http.StatusOK, &services.ApiErrorMsg{
		Code: code,
	})
}

func GetRtcConf(ctx *gin.Context) {
	appkey := ctx.Query("app_key")
	if appkey == "" {
		ctx.JSON(http.StatusBadRequest, &services.ApiErrorMsg{
			Code: services.AdminErrorCode_ParamError,
			Msg:  "param illegal",
		})
		return
	}
	code, conf := services.GetRtcConf(appkey)
	if code != services.AdminErrorCode_Success {
		ctx.JSON(http.StatusOK, &services.ApiErrorMsg{
			Code: code,
		})
	} else {
		services.SuccessHttpResp(ctx, conf)
	}
}

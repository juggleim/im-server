package apis

import (
	"im-server/commons/tools"
	"im-server/services/admingateway/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func QryAppInfo(ctx *gin.Context) {
	appkey := ctx.Query("app_key")
	appinfo := services.QryApp(appkey)
	services.SuccessHttpResp(ctx, appinfo)
}

func CreateApp(ctx *gin.Context) {
	var req services.AppInfo
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, &services.ApiErrorMsg{
			Code: services.AdminErrorCode_ParamError,
			Msg:  "param illegal",
		})
		return
	}
	code, appinfo := services.CreateApp(req)
	if code != services.AdminErrorCode_Success {
		ctx.JSON(http.StatusOK, &services.ApiErrorMsg{
			Code: code,
			Msg:  "",
		})
	} else {
		services.SuccessHttpResp(ctx, appinfo)
	}
}

type CreateAppReq struct {
}

func QryApps(ctx *gin.Context) {
	offsetStr := ctx.Query("offset")
	limitStr := ctx.Query("limit")
	var limit int64 = 50
	if limitStr != "" {
		intVal, err := tools.String2Int64(limitStr)
		if err == nil && intVal > 0 && intVal <= 100 {
			limit = intVal
		}
	}
	apps := services.QryApps(limit, offsetStr)
	services.SuccessHttpResp(ctx, apps)
}

func UpdateAppConfigs(ctx *gin.Context) {
	var req services.AppConfigs
	if err := ctx.ShouldBindJSON(&req); err != nil || req.AppKey == "" {
		ctx.JSON(http.StatusBadRequest, &services.ApiErrorMsg{
			Code: services.AdminErrorCode_ParamError,
			Msg:  "param illegal",
		})
		return
	}
	code := services.UpdateAppConfigs(req.AppKey, req.Configs)
	ctx.JSON(http.StatusOK, &services.ApiErrorMsg{
		Code: code,
	})
}

func QryAppConfigs(ctx *gin.Context) {
	var req QryConfigsReq
	if err := ctx.ShouldBindJSON(&req); err != nil || req.AppKey == "" {
		ctx.JSON(http.StatusBadRequest, &services.ApiErrorMsg{
			Code: services.AdminErrorCode_ParamError,
			Msg:  "param illegal",
		})
		return
	}
	code, resp := services.QryAppConfigs(req.AppKey, req.ConfigKeys)
	if code != services.AdminErrorCode_Success {
		ctx.JSON(http.StatusOK, &services.ApiErrorMsg{
			Code: code,
		})
	} else {
		services.SuccessHttpResp(ctx, resp)
	}
}

type QryConfigsReq struct {
	AppKey     string   `json:"app_key"`
	ConfigKeys []string `json:"config_keys"`
}

package apis

import (
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/admingateway/services"
	"im-server/services/commonservices"
	"net/http"

	"github.com/gin-gonic/gin"
)

func QryMsgStatistic(ctx *gin.Context) {
	appkey := ctx.Query("app_key")
	statTypeStr := ctx.Query("stat_type")
	if appkey == "" {
		ctx.JSON(http.StatusBadRequest, &services.ApiErrorMsg{
			Code: services.AdminErrorCode_ParamError,
			Msg:  "param illegal",
		})
		return
	}
	var statType int64 = 0
	if statTypeStr != "" {
		intVal, err := tools.String2Int64(statTypeStr)
		if err == nil && intVal > 0 {
			statType = intVal
		}
	}
	channelTypeStr := ctx.Query("channel_type")
	var channelType int64 = 0
	if channelTypeStr != "" {
		intVal, err := tools.String2Int64(channelTypeStr)
		if err == nil && intVal > 0 {
			channelType = intVal
		}
	}
	startStr := ctx.Query("start")
	var start int64 = 0
	if startStr != "" {
		intVal, err := tools.String2Int64(startStr)
		if err == nil && intVal > 0 {
			start = intVal
		}
	}
	endStr := ctx.Query("end")
	var end int64 = 0
	if endStr != "" {
		intVal, err := tools.String2Int64(endStr)
		if err == nil && intVal > 0 {
			end = intVal
		}
	}
	items := commonservices.QryMsgStatistic(appkey, commonservices.StatType(statType), pbobjs.ChannelType(channelType), start, end)
	services.SuccessHttpResp(ctx, items)
}

func QryUserActivities(ctx *gin.Context) {
	appkey := ctx.Query("app_key")
	if appkey == "" {
		ctx.JSON(http.StatusBadRequest, &services.ApiErrorMsg{
			Code: services.AdminErrorCode_ParamError,
			Msg:  "param illegal",
		})
		return
	}
	startStr := ctx.Query("start")
	var start int64 = 0
	if startStr != "" {
		intVal, err := tools.String2Int64(startStr)
		if err == nil && intVal > 0 {
			start = intVal
		}
	}
	endStr := ctx.Query("end")
	var end int64 = 0
	if endStr != "" {
		intVal, err := tools.String2Int64(endStr)
		if err == nil && intVal > 0 {
			end = intVal
		}
	}
	items := commonservices.QryUserActivities(appkey, start, end)
	services.SuccessHttpResp(ctx, items)
}

func QryUserRegiste(ctx *gin.Context) {
	appkey := ctx.Query("app_key")
	if appkey == "" {
		ctx.JSON(http.StatusBadRequest, &services.ApiErrorMsg{
			Code: services.AdminErrorCode_ParamError,
			Msg:  "param illegal",
		})
		return
	}
	startStr := ctx.Query("start")
	var start int64 = 0
	if startStr != "" {
		intVal, err := tools.String2Int64(startStr)
		if err == nil && intVal > 0 {
			start = intVal
		}
	}
	endStr := ctx.Query("end")
	var end int64 = 0
	if endStr != "" {
		intVal, err := tools.String2Int64(endStr)
		if err == nil && intVal > 0 {
			end = intVal
		}
	}
	ret := commonservices.QryUserRegiste(appkey, start, end)
	services.SuccessHttpResp(ctx, ret)
}

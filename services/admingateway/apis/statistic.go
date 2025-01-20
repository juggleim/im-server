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
	if appkey == "" {
		ctx.JSON(http.StatusBadRequest, &services.ApiErrorMsg{
			Code: services.AdminErrorCode_ParamError,
			Msg:  "param illegal",
		})
		return
	}
	statTypeStrArr := ctx.QueryArray("stat_type")
	statTypes := []commonservices.StatType{}
	for _, statTypeStr := range statTypeStrArr {
		intVal, err := tools.String2Int64(statTypeStr)
		if err == nil && intVal > 0 {
			statTypes = append(statTypes, commonservices.StatType(intVal))
		}
	}
	if len(statTypes) <= 0 {
		statTypes = append(statTypes, commonservices.StatType_Up)
		statTypes = append(statTypes, commonservices.StatType_Down)
		statTypes = append(statTypes, commonservices.StatType_Dispatch)
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
			start = intVal / 1000
		}
	}
	endStr := ctx.Query("end")
	var end int64 = 0
	if endStr != "" {
		intVal, err := tools.String2Int64(endStr)
		if err == nil && intVal > 0 {
			end = intVal / 1000
		}
	}
	items := commonservices.QryMsgStatistic(appkey, statTypes, pbobjs.ChannelType(channelType), start, end)
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

func QryConnectCount(ctx *gin.Context) {
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
	ret := commonservices.QryConnect(appkey, start, end)
	services.SuccessHttpResp(ctx, ret)
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

func QryMaxConnectCount(ctx *gin.Context) {
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
	ret := commonservices.QryMaxConnect(appkey, start, end)
	services.SuccessHttpResp(ctx, ret)
}

package apis

import (
	"github.com/gin-gonic/gin"
	apiServices "im-server/services/admingateway/services"
	"im-server/services/logmanager/services"
	"strconv"
)

func Logs(ctx *gin.Context) {
	appKey := ctx.Query("appKey")
	lastKey := ctx.Query("lastKey")
	prev := ctx.Query("prev")
	logType := ctx.Query("logType")
	userId := ctx.Query("userId")
	session := ctx.Query("session")
	startTimeQ := ctx.Query("startTime")

	var table string
	if logType == "connection" {
		table = "connect"
	} else if logType == "sdk" {
		table = "session"
	} else if logType == "business" {
		table = "bus"
	}
	startTime, _ := strconv.ParseInt(startTimeQ, 10, 64)
	options := services.FetchLogOptions{
		AppKey:    appKey,
		Table:     table,
		UserId:    userId,
		StartTime: startTime,
		Session:   session,
		LastKey:   lastKey,
		Count:     10,
		Prev:      prev == "true",
	}
	logs, err := services.FetchAppServerLogs(options)

	if err != nil {
		apiServices.FailHttpResp(ctx, apiServices.AdminErrorCode_ServerErr)
		return
	}

	apiServices.SuccessHttpResp(ctx, map[string]interface{}{
		"list": logs,
	})
}

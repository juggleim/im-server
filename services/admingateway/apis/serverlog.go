package apis

import (
	"im-server/commons/tools"
	"im-server/services/admingateway/services"
	logService "im-server/services/logmanager/services"

	"github.com/gin-gonic/gin"
)

type ServerLogs struct {
	Logs []map[string]interface{} `json:"logs"`
}

func QryUserConnectLogs(ctx *gin.Context) {
	qryServerLogs(ctx, "user_connect")
}

func QryConnectLogs(ctx *gin.Context) {
	qryServerLogs(ctx, "connect")
}

func qryServerLogs(ctx *gin.Context, logType string) {
	appkey := ctx.Query("app_key")
	userId := ctx.Query("user_id")
	session := ctx.Query("session")
	startStr := ctx.Query("start")
	var start int64 = 0
	if startStr != "" {
		intVal, err := tools.String2Int64(startStr)
		if err == nil && intVal > 0 {
			start = intVal
		}
	}
	countStr := ctx.Query("count")
	var count int64 = 100
	if countStr != "" {
		intVal, err := tools.String2Int64(countStr)
		if err == nil && intVal > 0 {
			count = intVal
		}
	}
	ret := &ServerLogs{
		Logs: []map[string]interface{}{},
	}
	var logs []logService.LogEntity
	var err error
	if logType == "user_connect" {
		logs, err = logService.QryUserConnectLogs(appkey, userId, start, count)
	} else if logType == "connect" {
		logs, err = logService.QryConnectLogs(appkey, session, start, count)
	} else {
		services.SuccessHttpResp(ctx, ret)
		return
	}
	if err == nil {
		for _, log := range logs {
			var item map[string]interface{}
			err := tools.JsonUnMarshal([]byte(log.Value), &item)
			if err == nil {
				ret.Logs = append(ret.Logs, item)
			}
		}
	}
	services.SuccessHttpResp(ctx, ret)
}

package apis

import (
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/admingateway/services"
	logService "im-server/services/logmanager/services"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/proto"
)

type ServerLogs struct {
	Logs []map[string]interface{} `json:"logs"`
}

func QryUserConnectLogs(ctx *gin.Context) {
	qryServerLogs(ctx, string(logService.ServerLogType_UserConnect))
}

func QryConnectLogs(ctx *gin.Context) {
	qryServerLogs(ctx, string(logService.ServerLogType_Connect))
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

	var targetId string
	if logType == string(logService.ServerLogType_UserConnect) {
		targetId = userId
	} else if logType == string(logService.ServerLogType_Connect) {
		targetId = session
	} else {
		services.FailHttpResp(ctx, services.AdminErrorCode_ParamError)
		return
	}
	services.SetCtxString(ctx, services.CtxKey_AppKey, appkey)
	code, resp, err := services.SyncApiCall(ctx, "qry_vlog", "", targetId, &pbobjs.QryServerLogsReq{
		LogType: logType,
		UserId:  userId,
		Session: session,
		Start:   start,
		Count:   count,
	}, func() proto.Message {
		return &pbobjs.QryServerLogsResp{}
	})
	if err != nil {
		services.FailHttpResp(ctx, services.AdminErrorCode_ServerErr, err.Error())
		return
	}
	if code != services.AdminErrorCode_Success {
		services.FailHttpResp(ctx, services.AdminErrorCode(code), "")
		return
	}
	logsResp := resp.(*pbobjs.QryServerLogsResp)

	ret := &ServerLogs{
		Logs: []map[string]interface{}{},
	}
	for _, logStr := range logsResp.Logs {
		var item map[string]interface{}
		err := tools.JsonUnMarshal([]byte(logStr), &item)
		if err == nil {
			ret.Logs = append(ret.Logs, item)
		}
	}
	services.SuccessHttpResp(ctx, ret)
}

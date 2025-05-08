package apis

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/admingateway/services"
	"im-server/services/commonservices"
	logService "im-server/services/logmanager/services"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/proto"
)

var hisMsgRedirectMethods map[string]bool

func init() {
	hisMsgRedirectMethods = make(map[string]bool)
	hisMsgRedirectMethods["clean_hismsg"] = true
	hisMsgRedirectMethods["del_msg"] = true
	hisMsgRedirectMethods["del_hismsg"] = true
	hisMsgRedirectMethods["mark_read"] = true
	hisMsgRedirectMethods["modify_msg"] = true
	hisMsgRedirectMethods["qry_hismsgs"] = true
	hisMsgRedirectMethods["qry_first_unread_msg"] = true
	hisMsgRedirectMethods["qry_hismsg_by_ids"] = true
	hisMsgRedirectMethods["qry_read_infos"] = true
	hisMsgRedirectMethods["recall_msg"] = true
	hisMsgRedirectMethods["msg_search"] = true
}

type ServerLogs struct {
	Logs []map[string]interface{} `json:"logs"`
}

func QryUserConnectLogs(ctx *gin.Context) {
	qryServerLogs(ctx, string(logService.ServerLogType_UserConnect))
}

func QryConnectLogs(ctx *gin.Context) {
	qryServerLogs(ctx, string(logService.ServerLogType_Connect))
}

func QryBusinessLogs(ctx *gin.Context) {
	qryServerLogs(ctx, string(logService.ServerLogType_Business))
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

	targetIdStr := ctx.Query("target_id")
	method := ctx.Query("method")
	seqIndexStr := ctx.Query("index")
	var seqIndex int64 = 0
	if seqIndexStr != "" {
		intVal, err := tools.String2Int64(seqIndexStr)
		if err == nil && intVal > 0 {
			seqIndex = intVal
		}
	}

	targetIds := []string{}
	if logType == string(logService.ServerLogType_UserConnect) {
		targetIds = append(targetIds, userId)
	} else if logType == string(logService.ServerLogType_Connect) {
		targetIds = append(targetIds, userId)
	} else if logType == string(logService.ServerLogType_Business) {
		if targetIdStr == "" || method == "" || userId == "" {
			services.FailHttpResp(ctx, services.AdminErrorCode_ParamError)
			return
		}
		targetIds = append(targetIds, targetIdStr)
		if _, ok := hisMsgRedirectMethods[method]; ok {
			targetIds = append(targetIds, commonservices.GetConversationId(userId, targetIdStr, pbobjs.ChannelType_Private))
		}
	} else {
		services.FailHttpResp(ctx, services.AdminErrorCode_ParamError)
		return
	}
	services.SetCtxString(ctx, services.CtxKey_AppKey, appkey)
	logs := []string{}
	for _, targetId := range targetIds {
		logs = qryLogsByRpc(services.ToRpcCtx(ctx, ""), targetId, &pbobjs.QryServerLogsReq{
			LogType: logType,
			UserId:  userId,
			Session: session,
			Start:   start,
			Count:   count,
			Index:   int32(seqIndex),
		})
		if len(logs) > 0 {
			break
		}
	}
	ret := &ServerLogs{
		Logs: []map[string]interface{}{},
	}
	for _, logStr := range logs {
		var item map[string]interface{}
		err := tools.JsonUnMarshal([]byte(logStr), &item)
		if err == nil {
			ret.Logs = append(ret.Logs, item)
		}
	}
	services.SuccessHttpResp(ctx, ret)
}

func qryLogsByRpc(ctx context.Context, targetId string, req *pbobjs.QryServerLogsReq) []string {
	code, resp, err := bases.SyncRpcCall(ctx, "qry_vlog", targetId, req, func() proto.Message {
		return &pbobjs.QryServerLogsResp{}
	})
	if err != nil {
		return []string{}
	}
	if code != errs.IMErrorCode_SUCCESS {
		return []string{}
	}
	logsResp := resp.(*pbobjs.QryServerLogsResp)
	return logsResp.Logs
}

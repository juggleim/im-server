package apis

import (
	"fmt"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/admingateway/services"
	apiService "im-server/services/apigateway/services"
	"im-server/services/commonservices"
	"im-server/services/commonservices/dbs"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ClientLogNtfReq struct {
	AppKey      string `json:"app_key"`
	UserId      string `json:"user_id"`
	Start       int64  `json:"start"`
	End         int64  `json:"end"`
	Description string `json:"description"`
	Platform    string `json:"platform"`
}

type LogCmd struct {
	Start    int64  `json:"start"`
	End      int64  `json:"end"`
	Platform string `json:"platform"`
}

func ClientLogNtf(ctx *gin.Context) {
	var req ClientLogNtfReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, &services.ApiErrorMsg{
			Code: services.AdminErrorCode_ParamError,
			Msg:  "param illegal",
		})
		return
	}
	traceId := fmt.Sprintf("api_%s", tools.GenerateUUIDShort11())
	item := dbs.ClientLogDao{
		AppKey:      req.AppKey,
		UserId:      req.UserId,
		CreatedTime: time.Now(),
		Start:       req.Start,
		End:         req.End,
		State:       dbs.ClientLogState_Default,
		TraceId:     traceId,
		Description: req.Description,
		Platform:    req.Platform,
	}
	ctx.Set(apiService.CtxKey_AppKey, req.AppKey)
	ctx.Set(apiService.CtxKey_Session, traceId)

	logCmd := &LogCmd{
		Start:    req.Start,
		End:      req.End,
		Platform: req.Platform,
	}
	code, resp, err := apiService.SyncSendMsg(ctx, "s_msg", "clientlog", req.UserId, &pbobjs.UpMsg{
		MsgType:    "jg:logcmd",
		MsgContent: []byte(tools.ToJson(logCmd)),
		Flags:      commonservices.SetCmdMsg(0),
	}, false)
	if err != nil || code != errs.IMErrorCode_SUCCESS {
		item.State = dbs.ClientLogState_SendFail
		if err != nil {
			item.FailReason = err.Error()
		} else {
			item.FailReason = fmt.Sprintf("code:%d", code)
		}
	} else {
		item.State = dbs.ClientLogState_SendOK
		item.MsgId = resp.MsgId
	}
	dao := &dbs.ClientLogDao{}
	dao.Create(item)
	services.SuccessHttpResp(ctx, nil)
}

type ClientLogItem struct {
	ID          string `json:"id"`
	UserId      string `json:"user_id"`
	Start       int64  `json:"start"`
	End         int64  `json:"end"`
	CreatedTime int64  `json:"created_time"`
	State       int    `json:"state"`
	Platform    string `json:"platform"`
	DeviceId    string `json:"device_id"`
	TraceId     string `json:"trace_id"`
	MsgId       string `json:"msg_id"`
	FailReason  string `json:"fail_reason"`
}
type ClientLogItems struct {
	Items  []*ClientLogItem `json:"items"`
	Offset string           `json:"offset"`
}

func ClientLogList(ctx *gin.Context) {
	appkey := ctx.Query("app_key")
	userId := ctx.Query("user_id")
	if appkey == "" {
		ctx.JSON(http.StatusBadRequest, &services.ApiErrorMsg{
			Code: services.AdminErrorCode_ParamError,
		})
		return
	}
	offset := ctx.Query("offset")
	var startId int64 = 0
	if offset != "" {
		intVal, err := tools.DecodeInt(offset)
		if err == nil && intVal > 0 {
			startId = intVal
		}
	}
	limitStr := ctx.Query("limit")
	var limit int64 = 100
	if limitStr != "" {
		intVal, err := tools.String2Int64(limitStr)
		if err == nil && intVal > 0 {
			limit = intVal
		}
	}
	if limit > 100 {
		limit = 100
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
	ret := &ClientLogItems{
		Items: []*ClientLogItem{},
	}
	dao := dbs.ClientLogDao{}
	list, err := dao.QryLogs(appkey, userId, start, end, startId, limit)
	if err == nil {
		for _, item := range list {
			idStr, _ := tools.EncodeInt(item.ID)
			ret.Offset = idStr
			ret.Items = append(ret.Items, &ClientLogItem{
				ID:          idStr,
				UserId:      item.UserId,
				Start:       item.Start,
				End:         item.End,
				CreatedTime: item.CreatedTime.UnixMilli(),
				State:       int(item.State),
				Platform:    item.Platform,
				DeviceId:    item.DeviceId,
				TraceId:     item.TraceId,
				MsgId:       item.MsgId,
				FailReason:  item.FailReason,
			})
		}
	}
	services.SuccessHttpResp(ctx, ret)
}

func ClientLogDownload(ctx *gin.Context) {
	appkey := ctx.Query("app_key")
	idStr := ctx.Query("id")
	if appkey == "" || idStr == "" {
		ctx.JSON(http.StatusBadRequest, &services.ApiErrorMsg{
			Code: services.AdminErrorCode_ParamError,
			Msg:  "param illegal",
		})
		return
	}
	id, err := tools.DecodeInt(idStr)
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusBadRequest, &services.ApiErrorMsg{
			Code: services.AdminErrorCode_ParamError,
			Msg:  "param illegal",
		})
		return
	}
	data := []byte{}
	dao := dbs.ClientLogDao{}
	log := dao.FindById(appkey, id)

	if log != nil {
		data = log.Log
	}
	ctx.Header("Content-Disposition", "attachement;filename=client.log")
	ctx.Data(http.StatusOK, "application/octet-stream", data)
}

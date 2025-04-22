package navigator

import (
	"compress/gzip"
	"im-server/commons/errs"
	"im-server/services/commonservices/dbs"
)

func UploadLogPlain(ctx *NavHttpContext) {
	//gzip
	if ctx.Request.Header.Get("Content-Encoding") == "gzip" {
		reader, err := gzip.NewReader(ctx.Request.Body)
		if err != nil {
			ctx.ResponseErr(errs.IMErrorCode_OTHER_CLIENT_LOG_PARAM_ILLEGAL)
			return
		}
		ctx.Request.Body = reader
	}
	var req = struct {
		Log   string `json:"log"`
		MsgId string `json:"msg_id"`
	}{}
	if err := ctx.BindJSON(&req); err != nil || req.Log == "" {
		ctx.ResponseErr(errs.IMErrorCode_OTHER_CLIENT_LOG_PARAM_ILLEGAL)
		return
	}
	if req.MsgId != "" && req.Log != "" {
		dao := dbs.ClientLogDao{}
		dao.Update(ctx.AppKey, req.MsgId, []byte(req.Log), dbs.ClientLogState_Uploaded)
	}
	ctx.ResponseSucc(nil)
}

type UploadLogStatusReq struct {
	MsgId  string `json:"msg_id"`
	LogUrl string `json:"log_url"`
	State  int    `json:"state"`
}

func LogStatus(ctx *NavHttpContext) {
	var req UploadLogStatusReq
	if err := ctx.BindJSON(&req); err != nil || req.MsgId == "" {
		ctx.ResponseErr(errs.IMErrorCode_OTHER_CLIENT_LOG_PARAM_ILLEGAL)
		return
	}
	dao := dbs.ClientLogDao{}
	dao.UpdateLogUrl(ctx.AppKey, req.MsgId, req.LogUrl, dbs.ClientLogState(req.State))
	ctx.ResponseSucc(nil)
}

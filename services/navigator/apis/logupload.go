package apis

import (
	"im-server/commons/errs"
	"im-server/commons/gmicro/utils"
	"im-server/commons/tools"
	"im-server/services/commonservices/dbs"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func UploadClientLog(ctx *gin.Context) {
	file, err := ctx.FormFile("log")
	if err != nil || file == nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_OTHER_CLIENT_LOG_PARAM_ILLEGAL)
		return
	}

	src, err := file.Open()
	if err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_OTHER_CLIENT_LOG_PARAM_ILLEGAL)
		return
	}
	defer src.Close()
	data := make([]byte, file.Size)
	_, err = src.Read(data)

	if err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_OTHER_CLIENT_LOG_PARAM_ILLEGAL)
		return
	}

	msgId := ctx.PostForm("msg_id")
	if msgId != "" && len(data) > 0 {
		dao := dbs.ClientLogDao{}
		dao.Update(ctx.GetString(CtxKey_AppKey), msgId, data, dbs.ClientLogState_Uploaded)
	}

	fileExt := filepath.Ext(file.Filename)
	filePath := "./client-logs/" + ctx.GetString(CtxKey_AppKey) + "-" + utils.GenerateUUIDShortString() + fileExt

	err = ctx.SaveUploadedFile(file, filePath)
	if err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_PARAM_REQUIRED)
		return
	}
	tools.SuccessHttpResp(ctx, nil)
}

func UploadClientLogPlain(ctx *gin.Context) {
	var req = struct {
		Log   string `json:"log"`
		MsgId string `json:"msg_id"`
	}{}
	if err := ctx.ShouldBindJSON(&req); err != nil || req.Log == "" {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_OTHER_CLIENT_LOG_PARAM_ILLEGAL)
		return
	}
	if req.MsgId != "" && req.Log != "" {
		dao := dbs.ClientLogDao{}
		dao.Update(ctx.GetString(CtxKey_AppKey), req.MsgId, []byte(req.Log), dbs.ClientLogState_Uploaded)
	}
	tools.SuccessHttpResp(ctx, nil)
}

package apis

import (
	"github.com/gin-gonic/gin"
	"im-server/commons/gmicro/utils"
	"os"
	"path/filepath"
)

func UploadClientLog(c *gin.Context) {

	file, err := c.FormFile("log")
	if err != nil {
		FailHttpResp(c, ApiErrorCode_ParamError, "No file found in request")
		return
	}
	if file == nil {
		FailHttpResp(c, ApiErrorCode_ParamError, "No file found in request")
		return
	}

	fileExt := filepath.Ext(file.Filename)
	filePath := "./client-logs/" + c.GetString(CtxKey_AppKey) + "-" + utils.GenerateUUIDShortString() + fileExt
	err = c.SaveUploadedFile(file, filePath)
	if err != nil {
		FailHttpResp(c, ApiErrorCode_ServerErr, "Error while saving file")
		return
	}
	SuccessHttpResp(c, nil)
}

func UploadClientLogPlain(c *gin.Context) {
	var req = struct {
		Log string `json:"log"`
	}{}
	if err := c.ShouldBindJSON(&req); err != nil || req.Log == "" {
		FailHttpResp(c, ApiErrorCode_ParamError, "log is empty")
		return
	}

	filePath := "./client-logs/" + c.GetString(CtxKey_AppKey) + "-" + utils.GenerateUUIDShortString() + ".log"

	err := os.MkdirAll("./client-logs", os.ModePerm)
	if err != nil {
		FailHttpResp(c, ApiErrorCode_ParamError, "save log file failed.")
		return
	}
	err = os.WriteFile(filePath, []byte(req.Log), 0666)
	if err != nil {
		FailHttpResp(c, ApiErrorCode_ParamError, "save log file failed.")
		return
	}
	SuccessHttpResp(c, nil)
}

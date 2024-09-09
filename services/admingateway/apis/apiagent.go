package apis

import (
	"im-server/services/admingateway/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ApiBody struct {
	Method string `json:"method"`
	AppKey string `json:"app_key"`
	Path   string `json:"path"`
	Body   string `json:"body"`
}

func ApiAgent(ctx *gin.Context) {
	var req ApiBody
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, services.AdminErrorCode_ParamError)
		return
	}
	appInfo := services.QryApp(req.AppKey)
	if appInfo == nil {
		ctx.JSON(http.StatusForbidden, services.AdminErrorCode_ParamError)
		return
	}
	httpCode, resp := services.ApiAgent(req.Method, req.Path, req.Body, req.AppKey, appInfo.AppSecret)
	ctx.JSON(httpCode, resp)
}

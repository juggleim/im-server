package apis

import (
	"encoding/json"
	"im-server/commons/tools"
	"im-server/services/admingateway/services"
	"im-server/services/commonservices"
	"im-server/services/pushmanager/storages/dbs"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetAndroidPushConf(ctx *gin.Context) {
	appkey := ctx.Query("app_key")
	pushChannel := ctx.Query("push_channel")
	if appkey == "" || pushChannel == "" {
		ctx.JSON(http.StatusBadRequest, &services.ApiErrorMsg{
			Code: services.AdminErrorCode_Default,
		})
		return
	}
	ret := &commonservices.AndroidPushConf{
		AppKey:      appkey,
		PushChannel: pushChannel,
	}
	dao := dbs.AndroidPushConfDao{}
	androidDao, err := dao.Find(appkey, pushChannel)
	if err == nil && androidDao != nil {
		ret.Package = androidDao.Package
		_ = json.Unmarshal([]byte(androidDao.PushConf), &ret.Extra)

	}
	services.SuccessHttpResp(ctx, ret)
}

func SetAndroidPushConf(ctx *gin.Context) {
	var req commonservices.AndroidPushConf
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, &services.ApiErrorMsg{
			Code: services.AdminErrorCode_ParamError,
			Msg:  "param illegal",
		})
		return
	}
	pushConf := ""
	var conf any
	if strings.EqualFold(req.PushChannel, string(commonservices.PushChannel_Huawei)) {
		conf = tools.MapToStruct[commonservices.HuaweiPushConf](req.Extra)
	} else if strings.EqualFold(req.PushChannel, string(commonservices.PushChannel_Xiaomi)) {
		conf = tools.MapToStruct[commonservices.XiaomiPushConf](req.Extra)
	} else if strings.EqualFold(req.PushChannel, string(commonservices.PushChannel_OPPO)) {
		conf = tools.MapToStruct[commonservices.OppoPushConf](req.Extra)
	} else if strings.EqualFold(req.PushChannel, string(commonservices.PushChannel_VIVO)) {
		conf = tools.MapToStruct[commonservices.VivoPushConf](req.Extra)
	} else if strings.EqualFold(req.PushChannel, string(commonservices.PushChannel_Jpush)) {
		conf = tools.MapToStruct[commonservices.JPushConf](req.Extra)
	}
	pushConf = tools.ToJson(conf)

	//save to db
	dao := dbs.AndroidPushConfDao{}
	dao.Upsert(dbs.AndroidPushConfDao{
		AppKey:      req.AppKey,
		PushChannel: req.PushChannel,
		Package:     req.Package,
		PushConf:    pushConf,
	})
	services.SuccessHttpResp(ctx, nil)
}

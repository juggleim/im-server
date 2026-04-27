package apis

import (
	"errors"
	"im-server/commons/dbcommons"
	"im-server/services/admingateway/services"
	"im-server/services/commonservices/dbs"
	"im-server/services/commonservices/logs"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type InterceptorReq struct {
	ID              int64  `json:"id,omitempty"`
	AppKey          string `json:"app_key,omitempty"`
	Name            string `json:"name,omitempty"`
	Sort            int    `json:"sort"`
	RequestUrl      string `json:"request_url,omitempty"`
	RequestTemplate string `json:"request_template,omitempty"`
	SuccTemplate    string `json:"succ_template,omitempty"`
	IsAsync         int    `json:"is_async"`
	Conf            string `json:"conf,omitempty"`
	InterceptType   int    `json:"intercept_type"`
}

type InterceptorConditionReq struct {
	ID            int64  `json:"id,omitempty"`
	AppKey        string `json:"app_key,omitempty"`
	InterceptorId int64  `json:"interceptor_id"`
	ChannelType   string `json:"channel_type,omitempty"`
	MsgType       string `json:"msg_type,omitempty"`
	SenderId      string `json:"sender_id,omitempty"`
	ReceiverId    string `json:"receiver_id,omitempty"`
}

func AddInterceptor(ctx *gin.Context) {
	var req InterceptorReq
	if err := ctx.ShouldBindJSON(&req); err != nil || req.AppKey == "" || req.Name == "" {
		services.FailHttpResp(ctx, services.AdminErrorCode_ParamError, "param illegal")
		return
	}
	dao := dbs.InterceptorDao{}
	err := dao.Create(dbs.InterceptorDao{
		AppKey:          req.AppKey,
		Name:            req.Name,
		Sort:            req.Sort,
		RequestUrl:      req.RequestUrl,
		RequestTemplate: req.RequestTemplate,
		SuccTemplate:    req.SuccTemplate,
		IsAsync:         req.IsAsync,
		Conf:            req.Conf,
		InterceptType:   req.InterceptType,
	})
	if err != nil {
		logs.NewLogEntity().Error(err.Error())
		services.FailHttpResp(ctx, services.AdminErrorCode_Default)
		return
	}
	services.SuccessHttpResp(ctx, nil)
}

func DeleteInterceptor(ctx *gin.Context) {
	var req InterceptorReq
	if err := ctx.ShouldBindJSON(&req); err != nil || req.ID <= 0 || req.AppKey == "" {
		services.FailHttpResp(ctx, services.AdminErrorCode_ParamError, "param illegal")
		return
	}
	db := dbcommons.GetDb()
	err := db.Where("id=? and app_key=?", req.ID, req.AppKey).Delete(&dbs.InterceptorDao{}).Error
	if err != nil {
		logs.NewLogEntity().Error(err.Error())
		services.FailHttpResp(ctx, services.AdminErrorCode_Default)
		return
	}
	_ = db.Where("interceptor_id=? and app_key=?", req.ID, req.AppKey).Delete(&dbs.IcConditionDao{}).Error
	services.SuccessHttpResp(ctx, nil)
}

func UpdateInterceptor(ctx *gin.Context) {
	var req InterceptorReq
	if err := ctx.ShouldBindJSON(&req); err != nil || req.ID <= 0 || req.AppKey == "" {
		services.FailHttpResp(ctx, services.AdminErrorCode_ParamError, "param illegal")
		return
	}
	updateVals := map[string]interface{}{
		"name":             req.Name,
		"sort":             req.Sort,
		"request_url":      req.RequestUrl,
		"request_template": req.RequestTemplate,
		"succ_template":    req.SuccTemplate,
		"is_async":         req.IsAsync,
		"conf":             req.Conf,
		"intercept_type":   req.InterceptType,
	}
	err := dbcommons.GetDb().Model(&dbs.InterceptorDao{}).Where("id=? and app_key=?", req.ID, req.AppKey).Updates(updateVals).Error
	if err != nil {
		logs.NewLogEntity().Error(err.Error())
		services.FailHttpResp(ctx, services.AdminErrorCode_Default)
		return
	}
	services.SuccessHttpResp(ctx, nil)
}

func ListInterceptors(ctx *gin.Context) {
	appkey := ctx.Query("app_key")
	if appkey == "" {
		services.FailHttpResp(ctx, services.AdminErrorCode_ParamError, "param illegal")
		return
	}
	dao := dbs.InterceptorDao{}
	items, err := dao.QryInterceptors(appkey)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			services.SuccessHttpResp(ctx, map[string]interface{}{"items": []*InterceptorReq{}})
			return
		}
		logs.NewLogEntity().Error(err.Error())
		services.FailHttpResp(ctx, services.AdminErrorCode_Default)
		return
	}
	resp := make([]*InterceptorReq, 0, len(items))
	for _, item := range items {
		resp = append(resp, &InterceptorReq{
			ID:              item.ID,
			AppKey:          item.AppKey,
			Name:            item.Name,
			Sort:            item.Sort,
			RequestUrl:      item.RequestUrl,
			RequestTemplate: item.RequestTemplate,
			SuccTemplate:    item.SuccTemplate,
			IsAsync:         item.IsAsync,
			Conf:            item.Conf,
			InterceptType:   item.InterceptType,
		})
	}
	services.SuccessHttpResp(ctx, map[string]interface{}{"items": resp})
}

func AddInterceptorConditions(ctx *gin.Context) {
	var req InterceptorConditionReq
	if err := ctx.ShouldBindJSON(&req); err != nil || req.AppKey == "" || req.InterceptorId <= 0 {
		services.FailHttpResp(ctx, services.AdminErrorCode_ParamError, "param illegal")
		return
	}
	dao := dbs.IcConditionDao{}
	err := dao.Create(dbs.IcConditionDao{
		AppKey:        req.AppKey,
		InterceptorId: req.InterceptorId,
		ChannelType:   req.ChannelType,
		MsgType:       req.MsgType,
		SenderId:      req.SenderId,
		ReceiverId:    req.ReceiverId,
	})
	if err != nil {
		logs.NewLogEntity().Error(err.Error())
		services.FailHttpResp(ctx, services.AdminErrorCode_Default)
		return
	}
	services.SuccessHttpResp(ctx, nil)
}

func DeleteInterceptorConditions(ctx *gin.Context) {
	var req InterceptorConditionReq
	if err := ctx.ShouldBindJSON(&req); err != nil || req.AppKey == "" || req.ID <= 0 {
		services.FailHttpResp(ctx, services.AdminErrorCode_ParamError, "param illegal")
		return
	}
	err := dbcommons.GetDb().Where("id=? and app_key=?", req.ID, req.AppKey).Delete(&dbs.IcConditionDao{}).Error
	if err != nil {
		logs.NewLogEntity().Error(err.Error())
		services.FailHttpResp(ctx, services.AdminErrorCode_Default)
		return
	}
	services.SuccessHttpResp(ctx, nil)
}

func UpdateInterceptorConditions(ctx *gin.Context) {
	var req InterceptorConditionReq
	if err := ctx.ShouldBindJSON(&req); err != nil || req.ID <= 0 || req.AppKey == "" || req.InterceptorId <= 0 {
		services.FailHttpResp(ctx, services.AdminErrorCode_ParamError, "param illegal")
		return
	}
	updateVals := map[string]interface{}{
		"interceptor_id": req.InterceptorId,
		"channel_type":   req.ChannelType,
		"msg_type":       req.MsgType,
		"sender_id":      req.SenderId,
		"receiver_id":    req.ReceiverId,
	}
	err := dbcommons.GetDb().Model(&dbs.IcConditionDao{}).Where("id=? and app_key=?", req.ID, req.AppKey).Updates(updateVals).Error
	if err != nil {
		logs.NewLogEntity().Error(err.Error())
		services.FailHttpResp(ctx, services.AdminErrorCode_Default)
		return
	}
	services.SuccessHttpResp(ctx, nil)
}

func ListInterceptorConditions(ctx *gin.Context) {
	appkey := ctx.Query("app_key")
	interceptorID := ctx.Query("interceptor_id")
	if appkey == "" || interceptorID == "" {
		services.FailHttpResp(ctx, services.AdminErrorCode_ParamError, "param illegal")
		return
	}
	var rows []*dbs.IcConditionDao
	err := dbcommons.GetDb().Where("app_key=? and interceptor_id=?", appkey, interceptorID).Find(&rows).Error
	if err != nil {
		logs.NewLogEntity().Error(err.Error())
		services.FailHttpResp(ctx, services.AdminErrorCode_Default)
		return
	}
	resp := make([]*InterceptorConditionReq, 0, len(rows))
	for _, row := range rows {
		resp = append(resp, &InterceptorConditionReq{
			ID:            row.ID,
			AppKey:        row.AppKey,
			InterceptorId: row.InterceptorId,
			ChannelType:   row.ChannelType,
			MsgType:       row.MsgType,
			SenderId:      row.SenderId,
			ReceiverId:    row.ReceiverId,
		})
	}
	services.SuccessHttpResp(ctx, map[string]interface{}{"items": resp})
}

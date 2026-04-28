package apis

import (
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/apigateway/models"
	"im-server/services/apigateway/services"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/proto"
)

func AddBindDevice(ctx *gin.Context) {
	var req models.BindDevice
	if err := ctx.BindJSON(&req); err != nil || req.UserId == "" || req.DeviceId == "" {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}
	code, _, err := bases.SyncRpcCall(services.ToRpcCtx(ctx, ""), "add_bind_device", req.UserId, &pbobjs.BindDevice{
		DeviceId:      req.DeviceId,
		Platform:      req.Platform,
		DeviceCompany: req.DeviceCompany,
		DeviceModel:   req.DeviceModel,
	}, nil)
	if err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_INTERNAL_TIMEOUT)
		return
	}
	if code != errs.IMErrorCode_SUCCESS {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode(code))
		return
	}
	tools.SuccessHttpResp(ctx, nil)
}

func DelBindDevice(ctx *gin.Context) {
	var req models.BindDevice
	if err := ctx.BindJSON(&req); err != nil || req.UserId == "" || req.DeviceId == "" {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}
	code, _, err := bases.SyncRpcCall(services.ToRpcCtx(ctx, ""), "del_bind_device", req.UserId, &pbobjs.BindDevice{
		DeviceId:      req.DeviceId,
		Platform:      req.Platform,
		DeviceCompany: req.DeviceCompany,
		DeviceModel:   req.DeviceModel,
	}, nil)
	if err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_INTERNAL_TIMEOUT)
		return
	}
	if code != errs.IMErrorCode_SUCCESS {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode(code))
		return
	}
	tools.SuccessHttpResp(ctx, nil)
}

func QryBindDevices(ctx *gin.Context) {
	userId := ctx.Query("user_id")
	if userId == "" {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_PARAM_REQUIRED)
		return
	}
	code, respObj, err := bases.SyncRpcCall(services.ToRpcCtx(ctx, ""), "qry_bind_devices", userId, &pbobjs.Nil{}, func() proto.Message {
		return &pbobjs.BindDevices{}
	})
	if err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_INTERNAL_TIMEOUT)
		return
	}
	if code != errs.IMErrorCode_SUCCESS {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode(code))
		return
	}
	devicesPb, ok := respObj.(*pbobjs.BindDevices)
	if !ok || devicesPb == nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_INTERNAL_RESP_FAIL)
		return
	}
	out := &models.BindDevicesResp{Items: []*models.BindDevice{}}
	for _, d := range devicesPb.GetDevices() {
		out.Items = append(out.Items, &models.BindDevice{
			UserId:        userId,
			DeviceId:      d.GetDeviceId(),
			Platform:      d.GetPlatform(),
			DeviceCompany: d.GetDeviceCompany(),
			DeviceModel:   d.GetDeviceModel(),
			CreatedTime:   d.GetCreatedTime(),
		})
	}
	tools.SuccessHttpResp(ctx, out)
}

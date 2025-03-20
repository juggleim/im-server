package apis

import (
	"bytes"
	"encoding/base64"
	"image/png"

	"github.com/juggleim/jugglechat-server/apimodels"
	"github.com/juggleim/jugglechat-server/errs"
	"github.com/juggleim/jugglechat-server/services"
	"github.com/juggleim/jugglechat-server/utils"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
)

func QryUserInfo(ctx *HttpContext) {
	userId := ctx.Query("user_id")
	rpcCtx := ctx.ToRpcCtx()
	code, user := services.QryUserInfo(rpcCtx, userId)
	if code != errs.IMErrorCode_SUCCESS {
		ctx.ResponseErr(code)
		return
	}
	ctx.ResponseSucc(user)
}

func UpdateUser(ctx *HttpContext) {
	req := &apimodels.UserObj{}
	if err := ctx.BindJSON(req); err != nil {
		ctx.ResponseErr(errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	rpcCtx := ctx.ToRpcCtx()
	services.UpdateUser(rpcCtx, req)
	ctx.ResponseSucc(nil)
}

func UpdateUserSettings(ctx *HttpContext) {
	req := &apimodels.UserSettings{}
	if err := ctx.BindJSON(req); err != nil {
		ctx.ResponseErr(errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code := services.UpdateUserSettings(ctx.ToRpcCtx(), req)
	if code != errs.IMErrorCode_SUCCESS {
		ctx.ResponseErr(code)
		return
	}
	ctx.ResponseSucc(nil)
}

func SearchByPhone(ctx *HttpContext) {
	req := &apimodels.UserObj{}
	if err := ctx.BindJSON(req); err != nil {
		ctx.ResponseErr(errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	rpcCtx := ctx.ToRpcCtx()
	code, users := services.SearchByPhone(rpcCtx, req.Phone)
	if code != errs.IMErrorCode_SUCCESS {
		ctx.ResponseErr(code)
		return
	}
	ctx.ResponseSucc(users)
}

func QryUserQrCode(ctx *HttpContext) {
	userId := ctx.CurrentUserId

	m := map[string]interface{}{
		"action":  "add_friend",
		"user_id": userId,
	}
	buf := bytes.NewBuffer([]byte{})
	qrCode, _ := qr.Encode(utils.ToJson(m), qr.M, qr.Auto)
	qrCode, _ = barcode.Scale(qrCode, 400, 400)
	err := png.Encode(buf, qrCode)
	if err != nil {
		ctx.ResponseErr(errs.IMErrorCode_APP_DEFAULT)
		return
	}
	ctx.ResponseSucc(map[string]string{
		"qr_code": base64.StdEncoding.EncodeToString(buf.Bytes()),
	})
}

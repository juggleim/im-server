package apis

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/appbusiness/httputils"
	"im-server/services/appbusiness/models"
	"im-server/services/appbusiness/services"
	"im-server/services/appbusiness/storages"
	storageModels "im-server/services/appbusiness/storages/models"
	"image/png"
	"time"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"google.golang.org/protobuf/proto"
)

func Login(ctx *httputils.HttpContext) {
	req := &models.LoginReq{}
	if err := ctx.BindJson(req); err != nil || req.Account == "" {
		ctx.ResponseErr(errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	userId := tools.ShortMd5(req.Account)
	nickname := fmt.Sprintf("user%05d", tools.RandInt(100000))

	code, resp, err := bases.SyncRpcCall(ctx.ToRpcCtx(userId), "reg_user", userId, &pbobjs.UserInfo{
		UserId:   userId,
		Nickname: nickname,
		NoCover:  true,
		Settings: []*pbobjs.KvItem{
			{
				Key:   models.UserExtKey_FriendVerifyType,
				Value: tools.Int642String(int64(pbobjs.FriendVerifyType_NeedFriendVerify)),
			},
		},
	}, func() proto.Message {
		return &pbobjs.UserRegResp{}
	})
	if err != nil {
		ctx.ResponseErr(errs.IMErrorCode_APP_INTERNAL_TIMEOUT)
		return
	}
	if code != errs.IMErrorCode_SUCCESS {
		ctx.ResponseErr(code)
		return
	}
	ctx.ResponseSucc(resp)
}

func SmsLogin(ctx *httputils.HttpContext) {
	req := &models.SmsLoginReq{}
	if err := ctx.BindJson(req); err != nil || req.Phone == "" {
		ctx.ResponseErr(errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	succ := services.CheckPhoneSmsCode(req.Phone, req.Code)
	if succ {
		userId := tools.ShortMd5(req.Phone)
		nickname := fmt.Sprintf("user%05d", tools.RandInt(100000))
		code, resp, err := bases.SyncRpcCall(ctx.ToRpcCtx(userId), "reg_user", userId, &pbobjs.UserInfo{
			UserId:   userId,
			Nickname: nickname,
			NoCover:  true,
		}, func() proto.Message {
			return &pbobjs.UserRegResp{}
		})
		if err != nil {
			ctx.ResponseErr(errs.IMErrorCode_APP_INTERNAL_TIMEOUT)
			return
		}
		if code != errs.IMErrorCode_SUCCESS {
			ctx.ResponseErr(code)
			return
		}
		regResp := resp.(*pbobjs.UserRegResp)
		ctx.ResponseSucc(&models.LoginUserResp{
			UserId:        userId,
			Authorization: regResp.Token,
			NickName:      regResp.Nickname,
			Avatar:        regResp.UserPortrait,
			Status:        0,
			ImToken:       regResp.Token,
		})
	} else {
		ctx.ResponseErr(errs.IMErrorCode_APP_NOT_LOGIN)
		return
	}
}

func GenerateQrCode(ctx *httputils.HttpContext) {
	uuidStr := tools.GenerateUUIDString()
	qrCode, _ := qr.Encode(uuidStr, qr.M, qr.Auto)
	qrCode, _ = barcode.Scale(qrCode, 400, 400)
	buf := bytes.NewBuffer([]byte{})
	err := png.Encode(buf, qrCode)
	if err != nil {
		ctx.ResponseErr(errs.IMErrorCode_APP_DEFAULT)
		return
	}
	storage := storages.NewQrCodeRecordStorage()
	storage.Create(storageModels.QrCodeRecord{
		CodeId:      uuidStr,
		AppKey:      ctx.AppKey,
		CreatedTime: time.Now().UnixMilli(),
	})

	ctx.ResponseSucc(map[string]string{
		"id":      uuidStr,
		"qr_code": base64.StdEncoding.EncodeToString(buf.Bytes()),
	})
}

func CheckQrCode(ctx *httputils.HttpContext) {
	req := &models.QrCode{}
	if err := ctx.BindJson(req); err != nil || req.Id == "" {
		ctx.ResponseErr(errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	storage := storages.NewQrCodeRecordStorage()
	record, err := storage.FindById(ctx.AppKey, req.Id)
	if err != nil {
		ctx.ResponseErr(errs.IMErrorCode_APP_DEFAULT)
		return
	}
	if time.Now().UnixMilli()-record.CreatedTime > 10*60*1000 {
		ctx.ResponseErr(errs.IMErrorCode_APP_QRCODE_EXPIRED)
		return
	}
	if record.Status == storageModels.QrCodeRecordStatus_OK {
		userId := record.UserId
		code, resp, err := bases.SyncRpcCall(ctx.ToRpcCtx(userId), "reg_user", userId, &pbobjs.UserInfo{
			UserId:  userId,
			NoCover: true,
		}, func() proto.Message {
			return &pbobjs.UserRegResp{}
		})
		if err != nil {
			ctx.ResponseErr(errs.IMErrorCode_APP_INTERNAL_TIMEOUT)
			return
		}
		if code != errs.IMErrorCode_SUCCESS {
			ctx.ResponseErr(code)
			return
		}
		regResp := resp.(*pbobjs.UserRegResp)
		ctx.ResponseSucc(&models.LoginUserResp{
			UserId:        userId,
			Authorization: regResp.Token,
			NickName:      regResp.Nickname,
			Avatar:        regResp.UserPortrait,
			Status:        0,
			ImToken:       regResp.Token,
		})
	} else if record.Status == storageModels.QrCodeRecordStatus_Default {
		ctx.ResponseErr(errs.IMErrorCode_APP_CONTINUE)
		return
	}
}

func ConfirmQrCode(ctx *httputils.HttpContext) {
	req := &models.QrCode{}
	if err := ctx.BindJson(req); err != nil || req.Id == "" {
		ctx.ResponseErr(errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	appkey := ctx.AppKey
	userId := ctx.CurrentUserId
	storage := storages.NewQrCodeRecordStorage()
	err := storage.UpdateStatus(appkey, req.Id, storageModels.QrCodeRecordStatus_OK, userId)
	if err != nil {
		ctx.ResponseErr(errs.IMErrorCode_APP_DEFAULT)
		return
	}
	ctx.ResponseSucc(nil)
}

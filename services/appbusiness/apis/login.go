package apis

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/appbusiness/apimodels"
	"im-server/services/appbusiness/httputils"
	"im-server/services/appbusiness/services"
	"im-server/services/appbusiness/storages"
	storageModels "im-server/services/appbusiness/storages/models"
	userStorage "im-server/services/usermanager/storages"
	userModels "im-server/services/usermanager/storages/models"
	"image/png"
	"time"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/jinzhu/gorm"
	"google.golang.org/protobuf/proto"
)

func Login(ctx *httputils.HttpContext) {
	req := &apimodels.LoginReq{}
	if err := ctx.BindJson(req); err != nil || req.Account == "" {
		ctx.ResponseErr(errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	userId := tools.ShortMd5(req.Account)
	nickname := fmt.Sprintf("user%05d", tools.RandInt(100000))

	code, resp, err := bases.SyncRpcCall(ctx.ToRpcCtx(), "reg_user", userId, &pbobjs.UserInfo{
		UserId:   userId,
		Nickname: nickname,
		NoCover:  true,
		Settings: []*pbobjs.KvItem{
			{
				Key:   apimodels.UserExtKey_FriendVerifyType,
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

func SmsSend(ctx *httputils.HttpContext) {
	req := &apimodels.SmsLoginReq{}
	if err := ctx.BindJson(req); err != nil || req.Phone == "" {
		ctx.ResponseErr(errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code := services.SmsSend(ctx.ToRpcCtx(), req.Phone)
	if code != errs.IMErrorCode_SUCCESS {
		ctx.ResponseErr(code)
		return
	}
	ctx.ResponseSucc(nil)
}

func SmsLogin(ctx *httputils.HttpContext) {
	req := &apimodels.SmsLoginReq{}
	if err := ctx.BindJson(req); err != nil || req.Phone == "" {
		ctx.ResponseErr(errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code := services.CheckPhoneSmsCode(ctx.ToRpcCtx(), req.Phone, req.Code)
	if code == errs.IMErrorCode_SUCCESS {
		appkey := ctx.AppKey
		userId := tools.ShortMd5(req.Phone)
		nickname := fmt.Sprintf("user%05d", tools.RandInt(100000))
		storage := userStorage.NewUserStorage()
		user, err := storage.FindByPhone(appkey, req.Phone)
		if err == nil && user != nil {
			userId = user.UserId
			nickname = user.Nickname
		} else {
			user, err = storage.FindByUserId(appkey, userId)
			if err == nil && user != nil {
				userId = user.UserId
				nickname = user.Nickname
			} else {
				if err != gorm.ErrRecordNotFound {
					ctx.ResponseErr(errs.IMErrorCode_APP_NOT_LOGIN)
					return
				}
				userId = tools.GenerateUUIDShort11()
				err = storage.Create(userModels.User{
					UserId:   userId,
					Nickname: nickname,
					Phone:    req.Phone,
					AppKey:   appkey,
				})
				if err != nil {
					ctx.ResponseErr(errs.IMErrorCode_APP_NOT_LOGIN)
					return
				}
			}
		}

		code, resp, err := bases.SyncRpcCall(ctx.ToRpcCtx(), "reg_user", userId, &pbobjs.UserInfo{
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
		ctx.ResponseSucc(&apimodels.LoginUserResp{
			UserId:        userId,
			Authorization: regResp.Token,
			NickName:      regResp.Nickname,
			Avatar:        regResp.UserPortrait,
			Status:        0,
			ImToken:       regResp.Token,
		})
	} else {
		ctx.ResponseErr(code)
		return
	}
}

func GenerateQrCode(ctx *httputils.HttpContext) {
	uuidStr := tools.GenerateUUIDString()
	m := map[string]interface{}{
		"action": "login",
		"code":   uuidStr,
	}
	qrCode, _ := qr.Encode(tools.ToJson(m), qr.M, qr.Auto)
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
	req := &apimodels.QrCode{}
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
		code, resp, err := bases.SyncRpcCall(ctx.ToRpcCtx(), "reg_user", userId, &pbobjs.UserInfo{
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
		ctx.ResponseSucc(&apimodels.LoginUserResp{
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
	req := &apimodels.QrCode{}
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

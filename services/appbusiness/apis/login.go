package apis

import (
	"fmt"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/appbusiness/httputils"
	"im-server/services/appbusiness/models"
	"im-server/services/appbusiness/services"

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

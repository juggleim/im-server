package services

import (
	"context"
	"errors"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/logs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/connectmanager/server/codec"
	"im-server/services/connectmanager/server/imcontext"
	"time"
)

func KickUser(ctx context.Context, req *pbobjs.KickUserReq, code errs.IMErrorCode) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	ctxs := []imcontext.WsHandleContext{}
	ctxMap := GetConnectCtxByUser(appkey, req.UserId)
	platformMap := map[string]bool{}
	for _, plat := range req.Platforms {
		if _, exist := supportPlatforms[plat]; exist {
			platformMap[plat] = true
		}
	}
	deviceIdMap := map[string]bool{}
	for _, deviceId := range req.DeviceIds {
		deviceIdMap[deviceId] = true
	}

	if len(req.Platforms) > 0 || len(req.DeviceIds) > 0 {
		for _, ctx := range ctxMap {
			if len(req.Platforms) > 0 {
				platform := imcontext.GetContextAttrString(ctx, imcontext.StateKey_Platform)
				if _, exist := platformMap[platform]; exist {
					ctxs = append(ctxs, ctx)
				}
			}
			if len(req.DeviceIds) > 0 {
				deviceId := imcontext.GetContextAttrString(ctx, imcontext.StateKey_DeviceID)
				if _, exist := deviceIdMap[deviceId]; exist {
					ctxs = append(ctxs, ctx)
				}
			}
		}
	} else {
		for _, ctx := range ctxMap {
			ctxs = append(ctxs, ctx)
		}
	}

	for _, ctx := range ctxs {
		tmpCtx := ctx
		msgAck := codec.NewDisconnectMessage(&codec.DisconnectMsgBody{
			Code:      int32(code),
			Timestamp: time.Now().UnixMilli(),
			Ext:       req.Ext,
		})
		tmpCtx.Write(msgAck)
		logs.Infof("session:%s\taction:%s\tcode:%d", imcontext.GetConnSession(tmpCtx), imcontext.Action_Disconnect, msgAck.MsgBody.Code)
		go func() {
			Offline(tmpCtx, code)
			time.Sleep(time.Millisecond * 50)
			tmpCtx.Close(errors.New("kick off"))
		}()
	}
}

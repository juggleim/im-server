package services

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/pbdefines/pbobjs"
)

func SyncSubChange(ctx context.Context, req *pbobjs.SubRelChangeReq) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	targetIds := bases.GetTargetIdsFromCtx(ctx)
	subscriberId := bases.GetRequesterIdFromCtx(ctx)
	deviceId := bases.GetDeviceIdFromCtx(ctx)
	if req.BusType == pbobjs.StatusSubBusType_UserStatus {
		for _, userId := range targetIds {
			userSub := GetUserSubscribers(appkey, userId)
			if req.IsAdd {
				userSub.AddSubscriber(subscriberId, deviceId)
			} else {
				userSub.RemoveSubscriber(subscriberId, deviceId)
			}
		}
	}
}

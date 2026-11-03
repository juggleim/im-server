package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices/logs"
	"im-server/services/connectmanager/services"

	"google.golang.org/protobuf/proto"
)

type UserOnlineStatusActor struct {
	bases.BaseActor
}

func (actor *UserOnlineStatusActor) OnReceive(ctx context.Context, input proto.Message) {
	onlineItems := []*pbobjs.UserOnlineItem{}
	if userOnlineReq, ok := input.(*pbobjs.UserOnlineStatusReq); ok {
		logs.WithContext(ctx).WithField("method", "qry_online_status").Infof("user_ids:%v", userOnlineReq.UserIds)
		appkey := bases.GetAppKeyFromCtx(ctx)
		if len(userOnlineReq.UserIds) > 0 {
			for _, userid := range userOnlineReq.UserIds {
				onlineDevices := services.GetOnlineDevicesByUser(appkey, userid)
				item := &pbobjs.UserOnlineItem{
					UserId:   userid,
					IsOnline: len(onlineDevices) > 0,
					Devices:  make([]*pbobjs.UserOnlineDevice, 0, len(onlineDevices)),
				}
				for _, device := range onlineDevices {
					item.Devices = append(item.Devices, &pbobjs.UserOnlineDevice{
						DeviceId: device.DeviceId,
						Platform: device.Platform,
					})
				}
				onlineItems = append(onlineItems, item)
			}
		}
	} else {
		logs.WithContext(ctx).Errorf("qry_online_status, input is illigal.")
	}
	ack := bases.CreateQueryAckWraper(ctx, errs.IMErrorCode_SUCCESS, &pbobjs.UserOnlineStatusResp{
		Items: onlineItems,
	})
	actor.Sender.Tell(ack, actorsystem.NoSender)
}

func (actor *UserOnlineStatusActor) CreateInputObj() proto.Message {
	return &pbobjs.UserOnlineStatusReq{}
}

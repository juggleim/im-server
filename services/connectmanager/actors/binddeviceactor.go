package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices/logs"
	"im-server/services/connectmanager/services"

	"google.golang.org/protobuf/proto"
)

type AddBindDeviceActor struct {
	bases.BaseActor
}

func (actor *AddBindDeviceActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.BindDevice); ok {
		logs.WithContext(ctx).Infof("user_id:%s\tdevice_id:%s\tplatform:%s", bases.GetRequesterIdFromCtx(ctx), req.DeviceId, req.Platform)
		errCode := services.AddBindDevice(ctx, req)
		ack := bases.CreateQueryAckWraper(ctx, errCode, nil)
		actor.Sender.Tell(ack, actorsystem.NoSender)
	} else {
		logs.WithContext(ctx).Error("input is illigal")
	}
}

func (actor *AddBindDeviceActor) CreateInputObj() proto.Message {
	return &pbobjs.BindDevice{}
}

type DelBindDeviceActor struct {
	bases.BaseActor
}

func (actor *DelBindDeviceActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.BindDevice); ok {
		logs.WithContext(ctx).Infof("user_id:%s\tdevice_id:%s\tplatform:%s", bases.GetRequesterIdFromCtx(ctx), req.DeviceId, req.Platform)
		errCode := services.DelBindDevice(ctx, req)
		ack := bases.CreateQueryAckWraper(ctx, errCode, nil)
		actor.Sender.Tell(ack, actorsystem.NoSender)
	} else {
		logs.WithContext(ctx).Error("input is illigal")
	}
}

func (actor *DelBindDeviceActor) CreateInputObj() proto.Message {
	return &pbobjs.BindDevice{}
}

type QryBindDevicesActor struct {
	bases.BaseActor
}

func (actor *QryBindDevicesActor) OnReceive(ctx context.Context, input proto.Message) {
	logs.WithContext(ctx).Infof("user_id:%s", bases.GetRequesterIdFromCtx(ctx))
	code, devices := services.QryBindDevices(ctx)
	ack := bases.CreateQueryAckWraper(ctx, code, devices)
	actor.Sender.Tell(ack, actorsystem.NoSender)
}

func (actor *QryBindDevicesActor) CreateInputObj() proto.Message {
	return &pbobjs.Nil{}
}

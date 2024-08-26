package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices/logs"
	"im-server/services/group/services"

	"google.golang.org/protobuf/proto"
)

type QryGrpSnapshotActor struct {
	bases.BaseActor
}

func (actor *QryGrpSnapshotActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.QryGrpSnapshotReq); ok {
		logs.WithContext(ctx).Infof("group_id:%s\tnearly_time:%d", req.GroupId, req.NearlyTime)
		code, info := services.QrySnapshot(ctx, req)
		ack := bases.CreateQueryAckWraper(ctx, code, info)
		actor.Sender.Tell(ack, actorsystem.NoSender)
	} else {
		ack := bases.CreateQueryAckWraper(ctx, errs.IMErrorCode_PBILLEGAL, nil)
		actor.Sender.Tell(ack, actorsystem.NoSender)
		logs.WithContext(ctx).Errorf("input is illigal")
	}
}

func (actor *QryGrpSnapshotActor) CreateInputObj() proto.Message {
	return &pbobjs.QryGrpSnapshotReq{}
}

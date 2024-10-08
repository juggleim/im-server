package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices/logs"
	"im-server/services/group/services"

	"google.golang.org/protobuf/proto"
)

type ImportGroupHisMsgActor struct {
	bases.BaseActor
}

func (actor *ImportGroupHisMsgActor) OnReceive(ctx context.Context, input proto.Message) {
	if msg, ok := input.(*pbobjs.UpMsg); ok {
		services.ImportGroupHisMsg(ctx, msg)
	} else {
		logs.WithContext(ctx).Errorf("input is illigal. val:%v", msg)
	}
}

func (actor *ImportGroupHisMsgActor) CreateInputObj() proto.Message {
	return &pbobjs.UpMsg{}
}

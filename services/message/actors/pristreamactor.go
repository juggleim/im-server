package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices/logs"

	"google.golang.org/protobuf/proto"
)

type PriStreamActor struct {
	bases.BaseActor
}

func (actor *PriStreamActor) OnReceive(ctx context.Context, input proto.Message) {
	if msg, ok := input.(*pbobjs.DownMsg); ok {
		logs.WithContext(ctx).Infof("%v", msg)
	} else {
		logs.WithContext(ctx).Error("input is illigal")
	}
}

func (actor *PriStreamActor) CreateInputObj() proto.Message {
	return &pbobjs.DownMsg{}
}

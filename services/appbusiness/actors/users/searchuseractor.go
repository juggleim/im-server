package users

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/commonservices"

	"google.golang.org/protobuf/proto"
)

type SearchUserActor struct {
	bases.BaseActor
}

func (actor *SearchUserActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.SearchUserReq); ok {
		targetUserId := tools.ShortMd5(req.Account)
		code, respObj, err := bases.SyncRpcCall(ctx, "qry_user_info", targetUserId, &pbobjs.UserIdReq{
			UserId:   targetUserId,
			AttTypes: []int32{int32(commonservices.AttItemType_Att)},
		}, func() proto.Message {
			return &pbobjs.UserInfo{}
		})
		var ret proto.Message
		if err == nil && respObj != nil {
			ret = respObj
		} else {
			ret = &pbobjs.UserInfo{
				UserId: targetUserId,
			}
		}
		ack := bases.CreateQueryAckWraper(ctx, code, ret)
		actor.Sender.Tell(ack, actorsystem.NoSender)
	}
}

func (actor *SearchUserActor) CreateInputObj() proto.Message {
	return &pbobjs.SearchUserReq{}
}

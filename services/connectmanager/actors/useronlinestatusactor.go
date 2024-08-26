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
				onlineCount := services.GetConnectCountByUser(appkey, userid)
				if onlineCount > 0 {
					onlineItems = append(onlineItems, &pbobjs.UserOnlineItem{
						UserId:   userid,
						IsOnline: true,
					})
				} else {
					onlineItems = append(onlineItems, &pbobjs.UserOnlineItem{
						UserId:   userid,
						IsOnline: false,
					})
				}
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

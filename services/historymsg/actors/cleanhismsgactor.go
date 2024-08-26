package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices/logs"
	"im-server/services/historymsg/services"
	"time"

	"google.golang.org/protobuf/proto"
)

type CleanHisMsgActor struct {
	bases.BaseActor
}

func (actor *CleanHisMsgActor) OnReceive(ctx context.Context, input proto.Message) {
	if cleanHisMsgReq, ok := input.(*pbobjs.CleanHisMsgReq); ok {
		userId := bases.GetRequesterIdFromCtx(ctx)
		targetId := cleanHisMsgReq.TargetId
		logs.WithContext(ctx).Infof("user_id:%s\ttargetId=%s\tchannel_type=%v\tclean_msg_time=%d\tclean_msg_time_offset:%d", userId, targetId, cleanHisMsgReq.ChannelType, cleanHisMsgReq.CleanMsgTime, cleanHisMsgReq.CleanTimeOffset)

		if cleanHisMsgReq.CleanMsgTime == 0 {
			cleanHisMsgReq.CleanMsgTime = time.Now().UnixMilli() - cleanHisMsgReq.CleanTimeOffset
		}

		code := services.CleanHisMsg(ctx, cleanHisMsgReq)
		qryAck := bases.CreateQueryAckWraper(ctx, code, nil)
		actor.Sender.Tell(qryAck, actorsystem.NoSender)
		logs.WithContext(ctx).Infof("result_code:%v", code)
	} else {
		logs.WithContext(ctx).Infof("user_id:%s\tinput is illegal", bases.GetRequesterIdFromCtx(ctx))
		qryAck := bases.CreateQueryAckWraper(ctx, errs.IMErrorCode_MSG_DEFAULT, nil)
		actor.Sender.Tell(qryAck, actorsystem.NoSender)
	}
}

func (actor *CleanHisMsgActor) CreateInputObj() proto.Message {
	return &pbobjs.CleanHisMsgReq{}
}

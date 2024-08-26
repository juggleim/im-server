package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices/logs"
	"im-server/services/sensitivemanager/dbs"
	"im-server/services/sensitivemanager/services"

	"google.golang.org/protobuf/proto"
)

type DelWordsActor struct {
	bases.BaseActor
}

func (actor *DelWordsActor) OnReceive(ctx context.Context, input proto.Message) {
	method := bases.GetMethodFromCtx(ctx)
	if req, ok := input.(*pbobjs.DelSensitiveWordsReq); ok {
		appkey := bases.GetAppKeyFromCtx(ctx)
		dao := dbs.SensitiveWordDao{}
		err := dao.DeleteWords(appkey, req.Words...)
		if err != nil {
			ack := bases.CreateQueryAckWraper(ctx, errs.IMErrorCode_API_INTERNAL_RESP_FAIL, nil)
			actor.Sender.Tell(ack, actorsystem.NoSender)
			logs.WithContext(ctx).WithField("method", method).Errorf("add words error:%v", err)
			return
		}
		filter := services.GetAppFilter(appkey)
		filter.DelWord(req.Words...)
		targetId := bases.GetTargetIdFromCtx(ctx)
		logs.WithContext(ctx).Infof("target_id:%s\treq:%v", targetId, req)
		ack := bases.CreateQueryAckWraper(ctx, errs.IMErrorCode_SUCCESS, nil)
		actor.Sender.Tell(ack, actorsystem.NoSender)
	} else {
		logs.WithContext(ctx).WithField("method", method).Infof("input is illegal")
	}
}

func (actor *DelWordsActor) CreateInputObj() proto.Message {
	return &pbobjs.DelSensitiveWordsReq{}
}

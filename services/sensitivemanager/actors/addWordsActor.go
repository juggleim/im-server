package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices/logs"
	"im-server/services/sensitivemanager/dbs"
	"im-server/services/sensitivemanager/sensitive"

	"github.com/samber/lo"
	"google.golang.org/protobuf/proto"
)

type AddWordsActor struct {
	bases.BaseActor
}

func (actor *AddWordsActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.AddSensitiveWordsReq); ok {
		appkey := bases.GetAppKeyFromCtx(ctx)
		dao := dbs.SensitiveWordDao{}
		err := dao.BatchUpsert(lo.Map(req.Words, func(item *pbobjs.SensitiveWord, index int) dbs.SensitiveWordDao {
			return dbs.SensitiveWordDao{
				AppKey:   appkey,
				Word:     item.Word,
				WordType: int(item.WordType),
			}
		}))
		if err != nil {
			ack := bases.CreateQueryAckWraper(ctx, errs.IMErrorCode_API_INTERNAL_RESP_FAIL, nil)
			actor.Sender.Tell(ack, actorsystem.NoSender)
			logs.WithContext(ctx).Errorf("add words error:%v", err)
			return
		}
		filter := sensitive.GetAppSensitiveFilter(bases.GetAppKeyFromCtx(ctx))
		filter.AddWord(req.Words...)
		targetId := bases.GetTargetIdFromCtx(ctx)
		logs.WithContext(ctx).Infof("target_id:%s\treq:%v", targetId, req)
		ack := bases.CreateQueryAckWraper(ctx, errs.IMErrorCode_SUCCESS, nil)
		actor.Sender.Tell(ack, actorsystem.NoSender)
	} else {
		logs.WithContext(ctx).Infof("input is illegal")
	}
}

func (actor *AddWordsActor) CreateInputObj() proto.Message {
	return &pbobjs.AddSensitiveWordsReq{}
}

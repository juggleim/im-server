package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices/logs"
	"im-server/services/sensitivemanager/services"

	"google.golang.org/protobuf/proto"
)

type QrySensitiveWordsActor struct {
	bases.BaseActor
}

func (actor *QrySensitiveWordsActor) OnReceive(ctx context.Context, input proto.Message) {
	method := bases.GetMethodFromCtx(ctx)
	if req, ok := input.(*pbobjs.QrySensitiveWordsReq); ok {
		//appkey := bases.GetAppKeyFromCtx(ctx)
		//var startId int64 = 0
		//if req.Offset != "" {
		//	intVal, err := tools.DecodeInt(req.Offset)
		//	if err == nil {
		//		startId = intVal
		//	}
		//}
		//var (
		//	code = errs.IMErrorCode_SUCCESS
		//	resp *pbobjs.QrySensitiveWordsResp
		//)
		//dao := dbs.SensitiveWordDao{}
		//list, err := dao.QrySensitiveWords(appkey, int64(req.Limit), startId)
		//if err != nil {
		//	code = errs.IMErrorCode_API_INTERNAL_RESP_FAIL
		//} else {
		//	resp = &pbobjs.QrySensitiveWordsResp{
		//		Words: lo.Map(list, func(item *dbs.SensitiveWordDao, index int) *pbobjs.SensitiveWord {
		//			idStr, _ := tools.EncodeInt(item.ID)
		//			return &pbobjs.SensitiveWord{
		//				Id:       idStr,
		//				Word:     item.Word,
		//				WordType: pbobjs.SensitiveWordType(item.WordType),
		//			}
		//		}),
		//	}
		//}
		code, resp := services.QrySensitiveWords(ctx, req)
		ack := bases.CreateQueryAckWraper(ctx, code, resp)
		actor.Sender.Tell(ack, actorsystem.NoSender)
	} else {
		logs.WithContext(ctx).WithField("method", method).Infof("input is illegal")
	}
}

func (actor *QrySensitiveWordsActor) CreateInputObj() proto.Message {
	return &pbobjs.QrySensitiveWordsReq{}
}

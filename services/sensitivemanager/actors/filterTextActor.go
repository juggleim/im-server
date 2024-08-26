package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices/logs"
	"im-server/services/sensitivemanager/services"

	"google.golang.org/protobuf/proto"
)

type FilterTextActor struct {
	bases.BaseActor
}

func (actor *FilterTextActor) OnReceive(ctx context.Context, input proto.Message) {
	method := bases.GetMethodFromCtx(ctx)
	if req, ok := input.(*pbobjs.SensitiveFilterReq); ok {
		appkey := bases.GetAppKeyFromCtx(ctx)
		filter := services.GetAppFilter(appkey)
		isDeny, replacedText := filter.ReplaceSensitiveWords(req.Text)
		targetId := bases.GetTargetIdFromCtx(ctx)
		logs.WithContext(ctx).Infof("target_id:%s\treq:%v", targetId, req)

		var handlerType pbobjs.SensitiveHandlerType
		if isDeny {
			handlerType = pbobjs.SensitiveHandlerType_deny
		} else if req.Text == replacedText {
			handlerType = pbobjs.SensitiveHandlerType_pass
		} else {
			handlerType = pbobjs.SensitiveHandlerType_replace
		}
		ack := bases.CreateQueryAckWraper(ctx, errs.IMErrorCode_SUCCESS, &pbobjs.SensitiveFilterResp{
			HandlerType:  handlerType,
			FilteredText: replacedText,
		})
		actor.Sender.Tell(ack, actorsystem.NoSender)
	} else {
		logs.WithContext(ctx).WithField("method", method).Infof("input is illegal")
	}
}

func (actor *FilterTextActor) CreateInputObj() proto.Message {
	return &pbobjs.SensitiveFilterReq{}
}

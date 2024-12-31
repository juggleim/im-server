package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices"
	"im-server/services/commonservices/logs"
	"im-server/services/commonservices/transengines"

	"google.golang.org/protobuf/proto"
)

type MultiTransActor struct {
	bases.BaseActor
}

func (actor *MultiTransActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.TransReq); ok {
		logs.WithContext(ctx).Infof("user_id:%s\tsource_lang:%s\ttarget_lang:%s\tlen:%d", bases.GetRequesterIdFromCtx(ctx), req.SourceLang, req.TargetLang, len(req.Items))
		transEngine := commonservices.GetTransEngine(bases.GetAppKeyFromCtx(ctx))
		resp := &pbobjs.TransReq{
			Items:      []*pbobjs.TransItem{},
			SourceLang: req.SourceLang,
			TargetLang: req.TargetLang,
		}
		code := errs.IMErrorCode_SUCCESS
		if transEngine != nil && transEngine != transengines.DefaultTransEngine {
			for _, item := range req.Items {
				result := transEngine.Translate(item.Content, []string{req.TargetLang})
				if len(result) > 0 {
					if afterTranslated, exist := result[req.TargetLang]; exist {
						resp.Items = append(resp.Items, &pbobjs.TransItem{
							Key:     item.Key,
							Content: afterTranslated,
						})
					}
				}
			}
		} else {
			code = errs.IMErrorCode_CONNECT_UNSUPPORTEDTOPIC
		}
		ack := bases.CreateQueryAckWraper(ctx, code, resp)
		actor.Sender.Tell(ack, actorsystem.NoSender)
	} else {
		logs.WithContext(ctx).Errorf("input is illigal. input:%v", input)
	}
}

func (actor *MultiTransActor) CreateInputObj() proto.Message {
	return &pbobjs.TransReq{}
}

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
	"sync"

	"google.golang.org/protobuf/proto"
)

type MultiTransActor struct {
	bases.BaseActor
}

func (actor *MultiTransActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.TransReq); ok {
		logs.WithContext(ctx).Infof("user_id:%s\tsource_lang:%s\ttarget_lang:%s\tlen:%d", bases.GetRequesterIdFromCtx(ctx), req.SourceLang, req.TargetLang, len(req.Items))
		resp := &pbobjs.TransReq{
			Items:      []*pbobjs.TransItem{},
			SourceLang: req.SourceLang,
			TargetLang: req.TargetLang,
		}
		code := errs.IMErrorCode_SUCCESS
		if req.TargetLang == "" || len(req.Items) <= 0 {
			code = errs.IMErrorCode_CONNECT_PARAM_REQUIRED
		} else {
			transEngine := commonservices.GetTransEngine(bases.GetAppKeyFromCtx(ctx))
			if transEngine != nil && transEngine != transengines.DefaultTransEngine {
				wg := &sync.WaitGroup{}
				for _, item := range req.Items {
					afterItem := &pbobjs.TransItem{
						Key:     item.Key,
						Content: item.Content,
					}
					resp.Items = append(resp.Items, afterItem)
					wg.Add(1)
					go func() {
						defer wg.Done()
						result := transEngine.Translate(afterItem.Content, []string{req.TargetLang})
						if len(result) > 0 {
							if afterTranslated, exist := result[req.TargetLang]; exist {
								afterItem.Content = afterTranslated
							}
						}
					}()
				}
				wg.Wait()
			} else {
				code = errs.IMErrorCode_CONNECT_UNSUPPORTEDTOPIC
			}
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

package services

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/appbusiness/storages"
	"im-server/services/appbusiness/storages/models"
)

func QryAiBots(ctx context.Context, req *pbobjs.QryAiBotsReq) (errs.IMErrorCode, *pbobjs.AiBotInfos) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	storage := storages.NewBotConfStorage()
	var startId int64 = 0
	if req.Offset != "" {
		intVal, err := tools.DecodeInt(req.Offset)
		if err == nil {
			startId = intVal
		}
	}
	ret := &pbobjs.AiBotInfos{
		Items: []*pbobjs.AiBotInfo{},
	}
	items, err := storage.QryBotConfsWithStatus(appkey, models.BotStatus_Enable, startId, req.Limit)
	if err == nil {
		for _, item := range items {
			ret.Offset, _ = tools.EncodeInt(item.ID)
			ret.Items = append(ret.Items, &pbobjs.AiBotInfo{
				BotId:    item.BotId,
				Nickname: item.Nickname,
				Avatar:   item.BotPortrait,
				BotType:  int32(item.BotType),
			})
		}
	}
	return errs.IMErrorCode_SUCCESS, ret
}

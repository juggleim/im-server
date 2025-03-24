package services

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/tools"
	"im-server/services/appbusiness/apimodels"
	"im-server/services/appbusiness/storages"
	"im-server/services/appbusiness/storages/models"
	"math"
)

func PromptAdd(ctx context.Context, req *apimodels.Prompt) (errs.IMErrorCode, *apimodels.Prompt) {
	storage := storages.NewPromptStorage()
	id, err := storage.Create(models.Prompt{
		UserId:  bases.GetRequesterIdFromCtx(ctx),
		Prompts: req.Prompts,
		AppKey:  bases.GetAppKeyFromCtx(ctx),
	})
	if err != nil {
		return errs.IMErrorCode_APP_ASSISTANT_PROMPT_DBERROR, nil
	}
	idStr, _ := tools.EncodeInt(id)
	return errs.IMErrorCode_SUCCESS, &apimodels.Prompt{
		Id: idStr,
	}
}

func PromptUpdate(ctx context.Context, req *apimodels.Prompt) errs.IMErrorCode {
	id, _ := tools.DecodeInt(req.Id)
	if id <= 0 {
		return errs.IMErrorCode_APP_REQ_BODY_ILLEGAL
	}
	storage := storages.NewPromptStorage()
	err := storage.UpdatePrompts(bases.GetAppKeyFromCtx(ctx), bases.GetRequesterIdFromCtx(ctx), id, req.Prompts)
	if err != nil {
		return errs.IMErrorCode_APP_ASSISTANT_PROMPT_DBERROR
	}
	return errs.IMErrorCode_SUCCESS
}

func PromptDel(ctx context.Context, req *apimodels.Prompt) errs.IMErrorCode {
	id, _ := tools.DecodeInt(req.Id)
	if id <= 0 {
		return errs.IMErrorCode_APP_REQ_BODY_ILLEGAL
	}
	storage := storages.NewPromptStorage()
	err := storage.DelPrompts(bases.GetAppKeyFromCtx(ctx), bases.GetRequesterIdFromCtx(ctx), id)
	if err != nil {
		return errs.IMErrorCode_APP_ASSISTANT_PROMPT_DBERROR
	}
	return errs.IMErrorCode_SUCCESS
}

func PromptBatchDel(ctx context.Context, req *apimodels.PromptIds) errs.IMErrorCode {
	ids := []int64{}
	for _, idStr := range req.Ids {
		id, _ := tools.DecodeInt(idStr)
		if id > 0 {
			ids = append(ids, id)
		}
	}
	storage := storages.NewPromptStorage()
	err := storage.BatchDelPrompts(bases.GetAppKeyFromCtx(ctx), bases.GetRequesterIdFromCtx(ctx), ids)
	if err != nil {
		return errs.IMErrorCode_APP_ASSISTANT_PROMPT_DBERROR
	}
	return errs.IMErrorCode_SUCCESS
}

func QryPrompts(ctx context.Context, count int64, offset string) (errs.IMErrorCode, *apimodels.Prompts) {
	var startId int64 = math.MaxInt64
	if offset != "" {
		id, _ := tools.DecodeInt(offset)
		if id > 0 {
			startId = id
		}
	}
	ret := &apimodels.Prompts{
		Items: []*apimodels.Prompt{},
	}
	storage := storages.NewPromptStorage()
	items, err := storage.QryPrompts(bases.GetAppKeyFromCtx(ctx), bases.GetRequesterIdFromCtx(ctx), count, startId)
	if err == nil {
		for _, item := range items {
			idStr, _ := tools.EncodeInt(item.ID)
			ret.Items = append(ret.Items, &apimodels.Prompt{
				Id:          idStr,
				Prompts:     item.Prompts,
				CreatedTime: item.CreatedTime,
			})
			ret.Offset = idStr
		}
	}
	return errs.IMErrorCode_SUCCESS, ret
}

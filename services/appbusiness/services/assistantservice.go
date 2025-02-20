package services

import (
	"bytes"
	"context"
	"fmt"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/tools"
	"im-server/services/appbusiness/apimodels"
	"im-server/services/appbusiness/storages"
	"im-server/services/appbusiness/storages/models"
	"im-server/services/botmsg/services"
	"math"
)

func AssistantAnswer(ctx context.Context, req *apimodels.AssistantAnswerReq) (errs.IMErrorCode, *apimodels.AssistantAnswerResp) {
	if req == nil || len(req.Msgs) <= 0 {
		return errs.IMErrorCode_APP_DEFAULT, nil
	}
	buf := bytes.NewBuffer([]byte{})
	userId := bases.GetRequesterIdFromCtx(ctx)
	for _, msg := range req.Msgs {
		if msg.SenderId != userId {
			buf.WriteString(fmt.Sprintf("对方:%s\n", msg.Content))
		} else {
			buf.WriteString(fmt.Sprintf("我:%s\n", msg.Content))
		}
	}
	buf.WriteString("帮我生成回复")
	answer := services.GenerateAnswer(ctx, buf.String())
	return errs.IMErrorCode_SUCCESS, &apimodels.AssistantAnswerResp{
		Answer: answer,
	}
}

func PromptAdd(ctx context.Context, req *apimodels.Prompt) errs.IMErrorCode {
	storage := storages.NewPromptStorage()
	err := storage.Create(models.Prompt{
		UserId:  bases.GetRequesterIdFromCtx(ctx),
		Prompts: req.Prompts,
		AppKey:  bases.GetAppKeyFromCtx(ctx),
	})
	if err != nil {
		return errs.IMErrorCode_APP_ASSISTANT_PROMPT_DBERROR
	}
	return errs.IMErrorCode_SUCCESS
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

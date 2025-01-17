package services

import (
	"bytes"
	"context"
	"fmt"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/services/appbusiness/models"
	"im-server/services/botmsg/services"
)

func AssistantAnswer(ctx context.Context, req *models.AssistantAnswerReq) (errs.IMErrorCode, *models.AssistantAnswerResp) {
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
	return errs.IMErrorCode_SUCCESS, &models.AssistantAnswerResp{
		Answer: answer,
	}
}

package apis

import (
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/apigateway/models"
	"im-server/services/apigateway/services"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	"google.golang.org/protobuf/proto"
)

func QryUserTags(ctx *gin.Context) {
	userIds := ctx.QueryArray("user_id")

	if len(userIds) == 0 {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}
	code, msg, err := services.SyncApiCall(ctx, "qry_user_tags", "", strings.Join(userIds, ","), &pbobjs.UserIds{
		UserIds: userIds,
	}, func() proto.Message {
		return &pbobjs.UserTagList{}
	})
	if err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_INTERNAL_TIMEOUT)
		return
	}
	if code != int32(errs.IMErrorCode_SUCCESS) {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode(code))
		return
	}
	ret := &models.UserTagsPayload{
		UserTags: []models.UserTag{},
	}
	utagList, ok := msg.(*pbobjs.UserTagList)
	if ok {
		for _, uTags := range utagList.UserTags {
			ret.UserTags = append(ret.UserTags, models.UserTag{
				UserID: uTags.UserId,
				Tags:   uTags.Tags,
			})
		}
	}
	tools.SuccessHttpResp(ctx, ret)
}

func AddUserTags(ctx *gin.Context) {
	var req models.UserTagsPayload
	if err := ctx.ShouldBindJSON(&req); err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}

	targetId := strings.Join(lo.Map(req.UserTags, func(item models.UserTag, index int) string {
		return item.UserID
	}), ",")
	code, _, err := services.SyncApiCall(ctx, "add_user_tags", "", targetId, &pbobjs.UserTagList{
		UserTags: lo.Map(req.UserTags, func(item models.UserTag, index int) *pbobjs.UserTag {
			return &pbobjs.UserTag{
				UserId: item.UserID,
				Tags:   item.Tags,
			}
		}),
	}, nil)
	if err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_INTERNAL_TIMEOUT)
		return
	}
	if code != int32(errs.IMErrorCode_SUCCESS) {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode(code))
		return
	}
	tools.SuccessHttpResp(ctx, nil)
}

func DelUserTags(ctx *gin.Context) {
	var req models.UserTagsPayload
	if err := ctx.ShouldBindJSON(&req); err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}
	targetId := strings.Join(lo.Map(req.UserTags, func(item models.UserTag, index int) string {
		return item.UserID
	}), ",")
	code, _, err := services.SyncApiCall(ctx, "del_user_tags", "", targetId, &pbobjs.UserTagList{
		UserTags: lo.Map(req.UserTags, func(item models.UserTag, index int) *pbobjs.UserTag {
			return &pbobjs.UserTag{
				UserId: item.UserID,
				Tags:   item.Tags,
			}
		}),
	}, nil)
	if err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_INTERNAL_TIMEOUT)
		return
	}
	if code != int32(errs.IMErrorCode_SUCCESS) {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode(code))
		return
	}
	tools.SuccessHttpResp(ctx, nil)
}

func ClearUserTags(ctx *gin.Context) {
	var req models.UserIds
	if err := ctx.ShouldBindJSON(&req); err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}
	code, _, err := services.SyncApiCall(ctx, "clear_user_tags", "", strings.Join(req.UserIds, ","), &pbobjs.UserIds{
		UserIds: req.UserIds,
	}, nil)
	if err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_INTERNAL_TIMEOUT)
		return
	}
	if code != int32(errs.IMErrorCode_SUCCESS) {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode(code))
		return
	}
	tools.SuccessHttpResp(ctx, nil)
}

func PushWithTags(ctx *gin.Context) {
	var req models.PushPayload
	if err := ctx.ShouldBindJSON(&req); err != nil || !req.Validate() {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}

	code, _, err := services.SyncApiCall(ctx, "push_with_tags", "", tools.ToJson(req.Condition), &pbobjs.PushNotificationWithTags{
		FromUserId: req.FromUserId,
		Condition: &pbobjs.PushNotificationWithTags_Condition{
			TagsAnd: req.Condition.TagsAnd,
			TagsOr:  req.Condition.TagsOr,
		},
		MsgBody: &pbobjs.PushNotificationWithTags_MsgBody{
			MsgType:    req.MsgBody.MsgType,
			MsgContent: req.MsgBody.MsgContent,
		},
		Notification: &pbobjs.PushNotificationWithTags_Notification{
			Title:    req.Notification.Title,
			PushText: req.Notification.PushText,
		},
	}, nil)
	if err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_INTERNAL_TIMEOUT)
		return
	}
	if code != int32(errs.IMErrorCode_SUCCESS) {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode(code))
		return
	}

	tools.SuccessHttpResp(ctx, map[string]string{
		"push_id": "todo",
	})
}

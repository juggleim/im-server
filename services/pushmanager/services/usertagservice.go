package services

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/pushmanager/storages"
	"im-server/services/pushmanager/storages/models"
	"time"

	"github.com/samber/lo"
)

func GetUserTags(ctx context.Context, userIds []string) (res *pbobjs.UserTagList, err error) {
	appKey := bases.GetAppKeyFromCtx(ctx)
	res = &pbobjs.UserTagList{
		UserTags: make([]*pbobjs.UserTag, 0, len(userIds)),
	}
	for _, userId := range userIds {
		tags, err := storages.NewTagStorage().GetUserTags(ctx, appKey, userId)
		if err != nil {
			return nil, err
		}
		res.UserTags = append(res.UserTags, &pbobjs.UserTag{
			UserId: userId,
			Tags:   tags,
		})
	}

	return
}

func AddUserTags(ctx context.Context, req *pbobjs.UserTagList) (err error) {
	appKey := bases.GetAppKeyFromCtx(ctx)
	for _, userTag := range req.UserTags {
		err = storages.NewTagStorage().AddUserTags(ctx, appKey, userTag.UserId, userTag.Tags...)
		if err != nil {
			return err
		}
	}
	return nil
}

func DelUserTags(ctx context.Context, req *pbobjs.UserTagList) (err error) {
	appKey := bases.GetAppKeyFromCtx(ctx)
	for _, userTag := range req.UserTags {
		err = storages.NewTagStorage().DeleteUserTags(ctx, appKey, userTag.UserId, userTag.Tags...)
		if err != nil {
			return err
		}
	}
	return nil
}

func CleanUserTags(ctx context.Context, userIds []string) (err error) {
	appKey := bases.GetAppKeyFromCtx(ctx)
	for _, userId := range userIds {
		err = storages.NewTagStorage().ClearUserTag(ctx, appKey, userId)
		if err != nil {
			return err
		}
	}
	return nil
}

func PushWithUserTags(ctx context.Context, req *pbobjs.PushNotificationWithTags) (err error) {
	appKey := bases.GetAppKeyFromCtx(ctx)

	var (
		userIDs []string
		perPage = 1000
	)

	for page := 1; ; page++ {
		list, err := storages.NewTagStorage().GetUserWithTags(ctx, appKey, models.Condition{
			TagsAnd: req.Condition.TagsAnd,
			TagsOr:  req.Condition.TagsOr,
		}, page, perPage)
		if err != nil {
			return err
		}
		userIDs = append(userIDs, list...)
		if len(list) < perPage {
			break
		}
	}

	go func() {
		idGroups := lo.Chunk(userIDs, 1000)
		for _, idGroup := range idGroups {
			for _, userId := range idGroup {
				if req.MsgBody != nil {
					pushRpc := bases.CreateServerPubWraper(ctx, req.FromUserId, userId, "s_msg", &pbobjs.UpMsg{
						MsgType:    req.MsgBody.MsgType,
						MsgContent: []byte(req.MsgBody.MsgContent),
					})
					bases.UnicastRouteWithNoSender(pushRpc)
				} else if req.Notification != nil {
					pushRpc := bases.CreateServerPubWraper(ctx, req.FromUserId, userId, "push", &pbobjs.PushData{
						Title:    req.Notification.Title,
						PushText: req.Notification.PushText,
					})
					bases.UnicastRouteWithNoSender(pushRpc)
				}
			}
			time.Sleep(time.Millisecond * 100)
		}
	}()
	return nil
}

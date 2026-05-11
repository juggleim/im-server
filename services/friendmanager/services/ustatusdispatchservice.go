package services

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/pbdefines/pbobjs"
	statusService "im-server/services/statussubscriptions/services"
)

func DispatchUserStatus(ctx context.Context, req *pbobjs.UserStatusFriDispatch) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	userId := bases.GetTargetIdFromCtx(ctx)
	container := GetFriendContainer(ctx, appkey, userId)
	excludeUserIds := map[string]bool{}
	for _, uId := range req.ExcludedUserIds {
		excludeUserIds[uId] = true
	}
	realUserIds := []string{}
	container.ForeachFriends(func(friendId string) {
		if _, exist := excludeUserIds[friendId]; !exist {
			realUserIds = append(realUserIds, friendId)
		}
	})
	statusService.Dispatch2Message(ctx, req.Msg, realUserIds)
}

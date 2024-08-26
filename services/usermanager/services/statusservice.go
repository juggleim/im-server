package services

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/pbdefines/pbobjs"
)

func QryUserStatus(ctx context.Context, targetIds []string) *pbobjs.UserInfos {
	appkey := bases.GetAppKeyFromCtx(ctx)
	resp := &pbobjs.UserInfos{
		UserInfos: []*pbobjs.UserInfo{},
	}
	for _, targetId := range targetIds {
		userInfo, exist := GetUserInfo(appkey, targetId)
		if exist {
			uInfo := &pbobjs.UserInfo{
				Statuses: []*pbobjs.KvItem{},
			}
			statuses := userInfo.GetStatus()
			for _, v := range statuses {
				uInfo.Statuses = append(uInfo.Statuses, &pbobjs.KvItem{
					Key:     v.ItemKey,
					Value:   v.ItemValue,
					UpdTime: v.UpdatedTime,
				})
			}
			resp.UserInfos = append(resp.UserInfos, uInfo)
		}
	}
	return resp
}

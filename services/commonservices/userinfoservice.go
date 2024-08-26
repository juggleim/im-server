package commonservices

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/pbdefines/pbobjs"

	"google.golang.org/protobuf/proto"
)

func GetSenderUserInfo(ctx context.Context) *pbobjs.UserInfo {
	userInfo := bases.GetSenderInfoFromCtx(ctx)
	if userInfo == nil {
		userId := bases.GetRequesterIdFromCtx(ctx)
		userInfo = GetUserInfoFromRpc(ctx, userId)
	}
	return userInfo
}

func GetUserInfoFromRpc(ctx context.Context, userId string) *pbobjs.UserInfo {
	return GetUserInfoFromRpcWithAttTypes(ctx, userId, []int32{int32(AttItemType_Att)})
}

func GetUserInfoFromRpcWithAttTypes(ctx context.Context, userId string, attTypes []int32) *pbobjs.UserInfo {
	_, respObj, err := bases.SyncRpcCall(ctx, "qry_user_info", userId, &pbobjs.UserIdReq{
		UserId:   userId,
		AttTypes: attTypes,
	}, func() proto.Message {
		return &pbobjs.UserInfo{}
	})
	if err == nil && respObj != nil {
		return respObj.(*pbobjs.UserInfo)
	}
	return &pbobjs.UserInfo{
		UserId: userId,
	}
}

func Map2KvItems(m map[string]string) []*pbobjs.KvItem {
	items := []*pbobjs.KvItem{}
	for k, v := range m {
		items = append(items, &pbobjs.KvItem{
			Key:   k,
			Value: v,
		})
	}
	return items
}

func Kvitems2Map(items []*pbobjs.KvItem) map[string]string {
	m := make(map[string]string)
	for _, item := range items {
		m[item.Key] = item.Value
	}
	return m
}

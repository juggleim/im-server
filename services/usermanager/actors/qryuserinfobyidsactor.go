package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices"
	"im-server/services/commonservices/logs"
	"im-server/services/usermanager/services"

	"google.golang.org/protobuf/proto"
)

type QryUserInfoByIdsActor struct {
	bases.BaseActor
}

func (actor *QryUserInfoByIdsActor) OnReceive(ctx context.Context, input proto.Message) {
	if userIdsReq, ok := input.(*pbobjs.UserIdsReq); ok {
		if len(userIdsReq.AttTypes) <= 0 {
			userIdsReq.AttTypes = append(userIdsReq.AttTypes, int32(commonservices.AttItemType_Att))
		}
		appkey := bases.GetAppKeyFromCtx(ctx)
		ret := &pbobjs.UserInfosResp{
			UserInfoMap: make(map[string]*pbobjs.UserInfo),
		}
		for _, userId := range userIdsReq.UserIds {
			userInfo, exist := services.GetUserInfo(appkey, userId)
			if exist && userInfo != nil {
				pbUserInfo := &pbobjs.UserInfo{
					UserId:   userId,
					UserType: pbobjs.UserType(userInfo.UserType),
				}
				for _, attType := range userIdsReq.AttTypes {
					if attType == int32(commonservices.AttItemType_Att) {
						pbUserInfo.Nickname = userInfo.Nickname
						pbUserInfo.UserPortrait = userInfo.UserPortrait
						pbUserInfo.UpdatedTime = userInfo.UpdatedTime
						pbUserInfo.ExtFields = commonservices.Map2KvItems(userInfo.ExtFields)
					} else if attType == int32(commonservices.AttItemType_Setting) {
						pbUserInfo.Settings = commonservices.Map2KvItems(userInfo.SettingFields)
					}
				}
				ret.UserInfoMap[userId] = pbUserInfo
			}
		}
		queryAck := bases.CreateQueryAckWraper(ctx, errs.IMErrorCode_SUCCESS, ret)
		actor.Sender.Tell(queryAck, actorsystem.NoSender)
	} else {
		logs.WithContext(ctx).Errorf("input is illigal")
	}
}

func (actor *QryUserInfoByIdsActor) CreateInputObj() proto.Message {
	return &pbobjs.UserIdsReq{}
}

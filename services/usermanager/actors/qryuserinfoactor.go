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

type QryUserInfoActor struct {
	bases.BaseActor
}

func (actor *QryUserInfoActor) OnReceive(ctx context.Context, input proto.Message) {
	if userIdReq, ok := input.(*pbobjs.UserIdReq); ok {
		appkey := bases.GetAppKeyFromCtx(ctx)
		if len(userIdReq.AttTypes) <= 0 {
			userIdReq.AttTypes = append(userIdReq.AttTypes, int32(commonservices.AttItemType_Att))
		}
		var code errs.IMErrorCode
		var userInfo *pbobjs.UserInfo
		user, exist := services.GetUserInfo(appkey, userIdReq.UserId)
		if exist && user != nil {
			userInfo = &pbobjs.UserInfo{
				UserId:   user.UserId,
				UserType: pbobjs.UserType(user.UserType),
			}
			for _, attType := range userIdReq.AttTypes {
				if attType == int32(commonservices.AttItemType_Att) {
					userInfo.Nickname = user.Nickname
					userInfo.UserPortrait = user.UserPortrait
					userInfo.ExtFields = commonservices.Map2KvItems(user.ExtFields)
					userInfo.UpdatedTime = user.UpdatedTime
				} else if attType == int32(commonservices.AttItemType_Setting) {
					userInfo.Settings = commonservices.Map2KvItems(user.SettingFields)
				} else if attType == int32(commonservices.AttItemType_Status) {
					status := user.GetStatus()
					userInfo.Statuses = []*pbobjs.KvItem{}
					for _, v := range status {
						userInfo.Statuses = append(userInfo.Statuses, &pbobjs.KvItem{
							Key:     v.ItemKey,
							Value:   v.ItemValue,
							UpdTime: v.UpdatedTime,
						})
					}
				}
			}
			code = errs.IMErrorCode_SUCCESS
		} else {
			code = errs.IMErrorCode_USER_NOT_EXIST
		}
		queryAck := bases.CreateQueryAckWraper(ctx, code, userInfo)
		actor.Sender.Tell(queryAck, actorsystem.NoSender)
	} else {
		logs.WithContext(ctx).Errorf("input is illigal")
	}
}

func (actor *QryUserInfoActor) CreateInputObj() proto.Message {
	return &pbobjs.UserIdReq{}
}

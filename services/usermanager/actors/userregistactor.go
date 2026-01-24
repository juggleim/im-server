package actors

import (
	"context"
	"time"

	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/commonservices"
	"im-server/services/commonservices/logs"
	"im-server/services/commonservices/tokens"
	"im-server/services/usermanager/services"

	"google.golang.org/protobuf/proto"
)

type UserRegistActor struct {
	bases.BaseActor
}

func (actor *UserRegistActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.UserInfo); ok {
		req.Nickname = tools.TruncateText(req.Nickname, 32)
		token := tokens.ImToken{
			AppKey:    bases.GetAppKeyFromCtx(ctx),
			UserId:    req.UserId,
			DeviceId:  "",
			TokenTime: time.Now().UnixMilli(),
			UserType:  req.UserType,
		}
		appInfo, exist := commonservices.GetAppInfo(token.AppKey)
		if exist && appInfo != nil {
			tokenStr, _ := token.ToTokenString([]byte(appInfo.AppSecureKey))
			if !req.NoCover || !services.CheckUserExist(ctx, req.UserId) {
				services.AddUser(ctx, req.UserId, req.Nickname, req.UserPortrait, req.ExtFields, req.Settings, req.UserType)
			}
			queryAck := bases.CreateQueryAckWraper(ctx, errs.IMErrorCode_SUCCESS, &pbobjs.UserRegResp{
				UserId: req.UserId,
				Token:  tokenStr,
			})
			actor.Sender.Tell(queryAck, actorsystem.NoSender)
		} else {
			queryAck := bases.CreateQueryAckWraper(ctx, errs.IMErrorCode_API_APP_NOT_EXISTED, &pbobjs.UserRegResp{
				UserId: req.UserId,
				Token:  "",
			})
			actor.Sender.Tell(queryAck, actorsystem.NoSender)
		}
	} else {
		logs.WithContext(ctx).Errorf("input is illigal")
	}
}

func (actor *UserRegistActor) CreateInputObj() proto.Message {
	return &pbobjs.UserInfo{}
}

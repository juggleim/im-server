package services

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/appbusiness/storages"
	"im-server/services/appbusiness/storages/models"
	"im-server/services/commonservices/logs"
	"math"
)

func AddFavoriteMsg(ctx context.Context, req *pbobjs.FavoriteMsg) errs.IMErrorCode {
	appkey := bases.GetAppKeyFromCtx(ctx)
	userId := bases.GetRequesterIdFromCtx(ctx)
	storage := storages.NewFavoriteMsgStorage()
	err := storage.Create(models.FavoriteMsg{
		UserId:      userId,
		SenderId:    req.SenderId,
		ReceiverId:  req.ReceiverId,
		ChannelType: req.ChannelType,
		MsgId:       req.MsgId,
		MsgTime:     req.MsgTime,
		MsgType:     req.MsgType,
		MsgContent:  req.MsgContent,
		AppKey:      appkey,
	})
	if err != nil {
		logs.WithContext(ctx).Errorf("save favorite msg fail:%s", err.Error())
		return errs.IMErrorCode_APP_DEFAULT
	}
	return errs.IMErrorCode_SUCCESS
}

func QryFavoriteMsgs(ctx context.Context, limit int64, offset string) (errs.IMErrorCode, *pbobjs.FavoriteMsgs) {
	var startId int64 = math.MaxInt64
	if offset != "" {
		id, err := tools.DecodeInt(offset)
		if err == nil && id > 0 {
			startId = id
		}
	}
	storage := storages.NewFavoriteMsgStorage()
	msgs, err := storage.QueryFavoriteMsgs(bases.GetAppKeyFromCtx(ctx), bases.GetRequesterIdFromCtx(ctx), startId, limit)
	if err != nil {
		logs.WithContext(ctx).Errorf("failed to query favorite msgs. err:%s", err.Error())
		return errs.IMErrorCode_APP_DEFAULT, nil
	}
	ret := &pbobjs.FavoriteMsgs{
		Items: []*pbobjs.FavoriteMsg{},
	}
	for _, msg := range msgs {
		ret.Offset, _ = tools.EncodeInt(msg.ID)
		ret.Items = append(ret.Items, &pbobjs.FavoriteMsg{
			SenderId:    msg.SenderId,
			ReceiverId:  msg.ReceiverId,
			ChannelType: int32(msg.ChannelType),
			MsgId:       msg.MsgId,
			MsgTime:     msg.MsgTime,
			MsgType:     msg.MsgType,
			MsgContent:  msg.MsgContent,
			CreatedTime: msg.CreatedTime.UnixMilli(),
		})
	}
	return errs.IMErrorCode_SUCCESS, ret
}

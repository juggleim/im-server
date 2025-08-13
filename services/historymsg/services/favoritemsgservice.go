package services

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/commonservices"
	"im-server/services/commonservices/logs"
	"im-server/services/historymsg/storages"
	"im-server/services/historymsg/storages/models"
	"math"
	"time"
)

func AddFavoriteMsgs(ctx context.Context, req *pbobjs.FavoriteMsgIds) errs.IMErrorCode {
	appkey := bases.GetAppKeyFromCtx(ctx)
	userId := bases.GetRequesterIdFromCtx(ctx)

	for _, msg := range req.Items {
		converId := commonservices.GetConversationId(msg.SenderId, msg.ReceiverId, msg.ChannelType)
		var downMsg *pbobjs.DownMsg
		//qry msg from history
		if msg.ChannelType == pbobjs.ChannelType_Private {
			hisStorage := storages.NewPrivateHisMsgStorage()
			hisMsg, err := hisStorage.FindById(appkey, converId, msg.SubChannel, msg.MsgId)
			if err != nil || hisMsg == nil {
				return errs.IMErrorCode_MSG_MSGNOTFOUND
			}
			msg := &pbobjs.DownMsg{}
			err = tools.PbUnMarshal(hisMsg.MsgBody, msg)
			if err != nil {
				return errs.IMErrorCode_MSG_MSGNOTFOUND
			}
			downMsg = msg
		} else if msg.ChannelType == pbobjs.ChannelType_Group {
			hisStorage := storages.NewGroupHisMsgStorage()
			hisMsg, err := hisStorage.FindById(appkey, converId, msg.SubChannel, msg.MsgId)
			if err != nil || hisMsg == nil {
				return errs.IMErrorCode_MSG_MSGNOTFOUND
			}
			msg := &pbobjs.DownMsg{}
			err = tools.PbUnMarshal(hisMsg.MsgBody, msg)
			if err != nil {
				return errs.IMErrorCode_MSG_MSGNOTFOUND
			}
			downMsg = msg
		}
		if downMsg == nil {
			return errs.IMErrorCode_MSG_MSGNOTFOUND
		}
		storage := storages.NewFavoriteMsgStorage()
		msgBs, _ := tools.PbMarshal(downMsg)
		err := storage.Create(models.FavoriteMsg{
			UserId:      userId,
			SenderId:    msg.SenderId,
			ReceiverId:  msg.ReceiverId,
			ChannelType: msg.ChannelType,
			SubChannel:  msg.SubChannel,
			MsgId:       msg.MsgId,
			MsgTime:     downMsg.MsgTime,
			MsgType:     downMsg.MsgType,
			MsgBody:     msgBs,
			CreatedTime: time.Now(),
			AppKey:      appkey,
		})
		if err != nil {
			logs.WithContext(ctx).Errorf("save favorite msgs fail:%s", err.Error())
			return errs.IMErrorCode_MSG_FAVORITEDUPLICATE
		}
	}
	return errs.IMErrorCode_SUCCESS
}

func DelFavoriteMsgs(ctx context.Context, req *pbobjs.FavoriteMsgIds) errs.IMErrorCode {
	msgIds := []string{}
	for _, msg := range req.Items {
		msgIds = append(msgIds, msg.MsgId)
	}
	if len(msgIds) > 0 {
		appkey := bases.GetAppKeyFromCtx(ctx)
		userId := bases.GetRequesterIdFromCtx(ctx)
		storage := storages.NewFavoriteMsgStorage()
		err := storage.BatchDelete(appkey, userId, msgIds)
		if err != nil {
			logs.WithContext(ctx).Errorf("del favorite msgs fail:%s", err.Error())
		}
	}
	return errs.IMErrorCode_SUCCESS
}

func QryFavoriteMsgs(ctx context.Context, req *pbobjs.QryFavoriteMsgsReq) (errs.IMErrorCode, *pbobjs.FavoriteMsgs) {
	userId := bases.GetRequesterIdFromCtx(ctx)
	var startId int64 = math.MaxInt64
	if req.Offset != "" {
		id, err := tools.DecodeInt(req.Offset)
		if err == nil && id > 0 {
			startId = id
		}
	}
	limit := req.Limit
	if limit <= 0 || limit > 100 {
		limit = 100
	}
	storage := storages.NewFavoriteMsgStorage()
	msgs, err := storage.QueryFavoriteMsgs(bases.GetAppKeyFromCtx(ctx), bases.GetRequesterIdFromCtx(ctx), startId, limit)
	if err != nil {
		logs.WithContext(ctx).Errorf("failed to query favorite msgs. err:%s", err.Error())
		return errs.IMErrorCode_MSG_MSGNOTFOUND, nil
	}
	ret := &pbobjs.FavoriteMsgs{
		Items: []*pbobjs.FavoriteMsg{},
	}
	for _, msg := range msgs {
		ret.Offset, _ = tools.EncodeInt(msg.ID)
		downMsg := &pbobjs.DownMsg{}
		err = tools.PbUnMarshal(msg.MsgBody, downMsg)
		if err == nil {
			if userId == msg.SenderId {
				downMsg.IsSend = true
				downMsg.TargetId = msg.ReceiverId
			}
			ret.Items = append(ret.Items, &pbobjs.FavoriteMsg{
				Msg:         downMsg,
				CreatedTime: msg.CreatedTime.UnixMilli(),
			})
		}
	}
	return errs.IMErrorCode_SUCCESS, ret
}

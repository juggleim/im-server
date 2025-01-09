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

func AddFavoriteMsg(ctx context.Context, req *pbobjs.AddFavoriteMsgReq) errs.IMErrorCode {
	appkey := bases.GetAppKeyFromCtx(ctx)
	userId := bases.GetRequesterIdFromCtx(ctx)
	converId := commonservices.GetConversationId(req.SenderId, req.ReceiverId, req.ChannelType)

	var downMsg *pbobjs.DownMsg
	//qry msg from history
	if req.ChannelType == pbobjs.ChannelType_Private {
		hisStorage := storages.NewPrivateHisMsgStorage()
		hisMsg, err := hisStorage.FindById(appkey, converId, req.MsgId)
		if err != nil || hisMsg == nil {
			return errs.IMErrorCode_MSG_DEFAULT
		}
		msg := &pbobjs.DownMsg{}
		err = tools.PbUnMarshal(hisMsg.MsgBody, msg)
		if err != nil {
			return errs.IMErrorCode_MSG_DEFAULT
		}
		downMsg = msg
	} else if req.ChannelType == pbobjs.ChannelType_Group {
		hisStorage := storages.NewGroupHisMsgStorage()
		hisMsg, err := hisStorage.FindById(appkey, converId, req.MsgId)
		if err != nil || hisMsg == nil {
			return errs.IMErrorCode_MSG_DEFAULT
		}
		msg := &pbobjs.DownMsg{}
		err = tools.PbUnMarshal(hisMsg.MsgBody, msg)
		if err != nil {
			return errs.IMErrorCode_MSG_DEFAULT
		}
		downMsg = msg
	}
	if downMsg == nil {
		return errs.IMErrorCode_MSG_DEFAULT
	}
	storage := storages.NewFavoriteMsgStorage()
	msgBs, _ := tools.PbMarshal(downMsg)
	err := storage.Create(models.FavoriteMsg{
		UserId:      userId,
		SenderId:    req.SenderId,
		ReceiverId:  req.ReceiverId,
		ChannelType: req.ChannelType,
		MsgId:       req.MsgId,
		MsgTime:     downMsg.MsgTime,
		MsgType:     downMsg.MsgType,
		MsgBody:     msgBs,
		CreatedTime: time.Now(),
		AppKey:      appkey,
	})
	if err != nil {
		logs.WithContext(ctx).Errorf("save favorite msg fail:%s", err.Error())
		return errs.IMErrorCode_MSG_DEFAULT
	}
	return errs.IMErrorCode_SUCCESS
}

func QryFavoriteMsgs(ctx context.Context, req *pbobjs.QryFavoriteMsgsReq) (errs.IMErrorCode, *pbobjs.FavoriteMsgs) {
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
		return errs.IMErrorCode_MSG_DEFAULT, nil
	}
	ret := &pbobjs.FavoriteMsgs{
		Items: []*pbobjs.FavoriteMsg{},
	}
	for _, msg := range msgs {
		ret.Offset, _ = tools.EncodeInt(msg.ID)
		downMsg := &pbobjs.DownMsg{}
		tools.PbUnMarshal(msg.MsgBody, downMsg)
		ret.Items = append(ret.Items, &pbobjs.FavoriteMsg{
			Msg:         downMsg,
			CreatedTime: msg.CreatedTime.UnixMilli(),
		})
	}
	return errs.IMErrorCode_SUCCESS, ret
}

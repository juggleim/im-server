package services

import (
	"context"
	"encoding/json"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices"
	"im-server/services/historymsg/storages"
	"im-server/services/historymsg/storages/models"
)

func SetMsgExt(ctx context.Context, req *pbobjs.MsgExtItem) errs.IMErrorCode {
	appkey := bases.GetAppKeyFromCtx(ctx)
	userId := bases.GetRequesterIdFromCtx(ctx)
	targetId := req.TargetId
	extStorage := storages.NewMsgExtStorage()
	err := extStorage.Upsert(models.MsgExt{
		AppKey: appkey,
		MsgId:  req.MsgId,
		Key:    req.Key,
		Value:  req.Value,
	})
	if err == nil {
		//set msg's ext state   TODO: need to add cache
		converId := commonservices.GetConversationId(userId, targetId, req.ChannelType)
		if req.ChannelType == pbobjs.ChannelType_Private {
			storage := storages.NewPrivateHisMsgStorage()
			storage.UpdateMsgExtState(appkey, converId, req.MsgId, 1)
		} else if req.ChannelType == pbobjs.ChannelType_Group {
			storage := storages.NewGroupHisMsgStorage()
			storage.UpdateMsgExtState(appkey, converId, req.MsgId, 1)
		}
		msgExt := &MsgExt{
			MsgId:       req.MsgId,
			TargetId:    req.TargetId,
			ChannelType: int(req.ChannelType),
			Exts:        []*ExtItem{},
		}
		msgExt.Exts = append(msgExt.Exts, &ExtItem{
			Key:   req.Key,
			Value: req.Value,
		})
		bs, _ := json.Marshal(msgExt)
		upMsg := &pbobjs.UpMsg{
			MsgType:    MsgExtCmdType,
			MsgContent: bs,
			Flags:      commonservices.SetCmdMsg(0),
		}
		if req.ChannelType == pbobjs.ChannelType_Private {
			commonservices.AsyncPrivateMsg(ctx, userId, targetId, upMsg)
		} else if req.ChannelType == pbobjs.ChannelType_Group {
			commonservices.AsyncGroupMsg(ctx, userId, req.TargetId, upMsg)
		}
	} else {
		return errs.IMErrorCode_MSG_MSGEXTDUPLICATE
	}
	return errs.IMErrorCode_SUCCESS
}

var MsgExtCmdType string = "jg:msgext"

type MsgExt struct {
	MsgId       string     `json:"msg_id"`
	TargetId    string     `json:"target_id"`
	ChannelType int        `json:"channel_type"`
	Exts        []*ExtItem `json:"exts"`
}
type ExtItem struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

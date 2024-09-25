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
	"time"
)

func SetMsgExt(ctx context.Context, req *pbobjs.MsgExt) errs.IMErrorCode {
	if req.Ext == nil || req.Ext.Key == "" {
		return errs.IMErrorCode_SUCCESS
	}
	appkey := bases.GetAppKeyFromCtx(ctx)
	userId := bases.GetRequesterIdFromCtx(ctx)
	targetId := req.TargetId
	extStorage := storages.NewMsgExtStorage()
	optTime := time.Now().UnixMilli()
	err := extStorage.Upsert(models.MsgExt{
		AppKey:      appkey,
		MsgId:       req.MsgId,
		Key:         req.Ext.Key,
		Value:       req.Ext.Value,
		CreatedTime: optTime,
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
			Key:       req.Ext.Key,
			Value:     req.Ext.Value,
			Timestamp: optTime,
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

func DelMsgExt(ctx context.Context, req *pbobjs.MsgExt) errs.IMErrorCode {
	if req.Ext == nil || req.Ext.Key == "" {
		return errs.IMErrorCode_SUCCESS
	}
	return errs.IMErrorCode_SUCCESS
}

func AddMsgExSet(ctx context.Context, req *pbobjs.MsgExt) errs.IMErrorCode {
	if req.Ext == nil || req.Ext.Key == "" {
		return errs.IMErrorCode_SUCCESS
	}
	//TODO: need to add cache
	appkey := bases.GetAppKeyFromCtx(ctx)
	msgId := req.MsgId
	userId := bases.GetRequesterIdFromCtx(ctx)
	exsetStorage := storages.NewMsgExSetStorage()
	optTime := time.Now().UnixMilli()
	err := exsetStorage.Create(models.MsgExSet{
		AppKey:      appkey,
		MsgId:       msgId,
		Key:         req.Ext.Key,
		Item:        req.Ext.Value,
		CreatedTime: optTime,
	})
	if err == nil {
		msgExSet := &MsgExt{
			MsgId:       msgId,
			TargetId:    req.TargetId,
			ChannelType: int(req.ChannelType),
			Exts:        []*ExtItem{},
		}
		msgExSet.Exts = append(msgExSet.Exts, &ExtItem{
			Key:       req.Ext.Key,
			Value:     req.Ext.Value,
			Timestamp: optTime,
		})
		bs, _ := json.Marshal(msgExSet)
		upMsg := &pbobjs.UpMsg{
			MsgType:    MsgExSetCmdType,
			MsgContent: bs,
			Flags:      commonservices.SetCmdMsg(0),
		}
		if req.ChannelType == pbobjs.ChannelType_Private {
			commonservices.AsyncPrivateMsg(ctx, userId, req.TargetId, upMsg)
		} else if req.ChannelType == pbobjs.ChannelType_Group {
			commonservices.AsyncGroupMsg(ctx, userId, req.TargetId, upMsg)
		}
	}
	return errs.IMErrorCode_SUCCESS
}

func DelMsgExSet(ctx context.Context, req *pbobjs.MsgExt) errs.IMErrorCode {
	if req.Ext == nil || req.Ext.Key == "" || req.Ext.Value == "" {
		return errs.IMErrorCode_SUCCESS
	}
	appkey := bases.GetAppKeyFromCtx(ctx)
	msgId := req.MsgId
	userId := bases.GetRequesterIdFromCtx(ctx)
	exsetStorage := storages.NewMsgExSetStorage()
	err := exsetStorage.Delete(appkey, msgId, req.Ext.Key, req.Ext.Value)
	if err == nil {
		msgExSet := &MsgExt{
			MsgId:       msgId,
			TargetId:    req.TargetId,
			ChannelType: int(req.ChannelType),
			Exts:        []*ExtItem{},
		}
		msgExSet.Exts = append(msgExSet.Exts, &ExtItem{
			IsDel:     1,
			Key:       req.Ext.Key,
			Value:     req.Ext.Value,
			Timestamp: time.Now().UnixMilli(),
		})
		bs, _ := json.Marshal(msgExSet)
		upMsg := &pbobjs.UpMsg{
			MsgType:    MsgExSetCmdType,
			MsgContent: bs,
			Flags:      commonservices.SetCmdMsg(0),
		}
		if req.ChannelType == pbobjs.ChannelType_Private {
			commonservices.AsyncPrivateMsg(ctx, userId, req.TargetId, upMsg)
		} else if req.ChannelType == pbobjs.ChannelType_Group {
			commonservices.AsyncGroupMsg(ctx, userId, req.TargetId, upMsg)
		}
	}
	return errs.IMErrorCode_SUCCESS
}

var MsgExtCmdType string = "jg:msgext"
var MsgExSetCmdType string = "jg:msgexset"

type MsgExt struct {
	MsgId       string     `json:"msg_id"`
	TargetId    string     `json:"target_id"`
	ChannelType int        `json:"channel_type"`
	Exts        []*ExtItem `json:"exts"`
}
type ExtItem struct {
	IsDel     int    `json:"is_del"`
	Key       string `json:"key"`
	Value     string `json:"value"`
	Timestamp int64  `json:"timestamp"`
}

package services

import (
	"context"
	"encoding/json"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/commonservices"
	"im-server/services/commonservices/msgdefines"
	"im-server/services/historymsg/storages"
	"time"
)

func SetMsgExt(ctx context.Context, req *pbobjs.MsgExt) errs.IMErrorCode {
	if req.Ext == nil || req.Ext.Key == "" {
		return errs.IMErrorCode_SUCCESS
	}
	appkey := bases.GetAppKeyFromCtx(ctx)
	userId := bases.GetRequesterIdFromCtx(ctx)
	targetId := req.TargetId
	converId := commonservices.GetConversationId(userId, req.TargetId, req.ChannelType)
	msgInfo := GetMsgInfo(appkey, converId, req.MsgId, req.ChannelType)
	optTime := time.Now().UnixMilli()
	code := msgInfo.SetMsgExt(&pbobjs.MsgExtItem{
		Key:       req.Ext.Key,
		Value:     req.Ext.Value,
		Timestamp: optTime,
		UserInfo: &pbobjs.UserInfo{
			UserId: userId,
		},
	})
	if code != errs.IMErrorCode_SUCCESS {
		return code
	}
	extItems := &pbobjs.MsgExtItems{
		MsgId: req.MsgId,
		Exts:  []*pbobjs.MsgExtItem{},
	}
	msgInfo.ForeachMsgExt(func(key string, ext *pbobjs.MsgExtItem) {
		extItem := &pbobjs.MsgExtItem{
			Key:       ext.Key,
			Value:     ext.Value,
			Timestamp: ext.Timestamp,
		}
		if ext.UserInfo != nil {
			extItem.UserInfo = &pbobjs.UserInfo{
				UserId: ext.UserInfo.UserId,
			}
		}
		extItems.Exts = append(extItems.Exts, extItem)
	})
	extItemsBs, _ := tools.PbMarshal(extItems)
	if req.ChannelType == pbobjs.ChannelType_Private {
		storage := storages.NewPrivateHisMsgStorage()
		storage.UpdateMsgExt(appkey, converId, req.MsgId, extItemsBs)
	} else if req.ChannelType == pbobjs.ChannelType_Group {
		storage := storages.NewGroupHisMsgStorage()
		storage.UpdateMsgExt(appkey, converId, req.MsgId, extItemsBs)
	}
	msgExt := &MsgExt{
		MsgId: req.MsgId,
		Exts:  []*ExtItem{},
	}
	msgExt.Exts = append(msgExt.Exts, &ExtItem{
		Key:       req.Ext.Key,
		Value:     req.Ext.Value,
		Timestamp: optTime,
		User: &UserInfo{
			UserId: userId,
		},
	})
	bs, _ := json.Marshal(msgExt)
	upMsg := &pbobjs.UpMsg{
		MsgType:    MsgExtCmdType,
		MsgContent: bs,
		Flags:      msgdefines.SetStateMsg(0),
	}
	if req.ChannelType == pbobjs.ChannelType_Private {
		commonservices.AsyncPrivateMsg(ctx, userId, targetId, upMsg)
	} else if req.ChannelType == pbobjs.ChannelType_Group {
		commonservices.AsyncGroupMsg(ctx, userId, req.TargetId, upMsg)
	}
	return errs.IMErrorCode_SUCCESS
}

func DelMsgExt(ctx context.Context, req *pbobjs.MsgExt) errs.IMErrorCode {
	if req.Ext == nil || req.Ext.Key == "" {
		return errs.IMErrorCode_SUCCESS
	}
	appkey := bases.GetAppKeyFromCtx(ctx)
	userId := bases.GetRequesterIdFromCtx(ctx)
	converId := commonservices.GetConversationId(userId, req.TargetId, req.ChannelType)
	msgInfo := GetMsgInfo(appkey, converId, req.MsgId, req.ChannelType)
	succ := msgInfo.DelMsgExt(req.Ext.Key)
	if succ {
		extItems := &pbobjs.MsgExtItems{
			MsgId: req.MsgId,
			Exts:  []*pbobjs.MsgExtItem{},
		}
		msgInfo.ForeachMsgExt(func(key string, ext *pbobjs.MsgExtItem) {
			extItem := &pbobjs.MsgExtItem{
				Key:       ext.Key,
				Value:     ext.Value,
				Timestamp: ext.Timestamp,
			}
			if ext.UserInfo != nil {
				extItem.UserInfo = &pbobjs.UserInfo{
					UserId: ext.UserInfo.UserId,
				}
			}
			extItems.Exts = append(extItems.Exts, extItem)
		})
		extItemsBs, _ := tools.PbMarshal(extItems)
		if req.ChannelType == pbobjs.ChannelType_Private {
			storage := storages.NewPrivateHisMsgStorage()
			storage.UpdateMsgExt(appkey, converId, req.MsgId, extItemsBs)
		} else if req.ChannelType == pbobjs.ChannelType_Group {
			storage := storages.NewGroupHisMsgStorage()
			storage.UpdateMsgExt(appkey, converId, req.MsgId, extItemsBs)
		}
		msgExSet := &MsgExt{
			MsgId: req.MsgId,
			Exts:  []*ExtItem{},
		}
		msgExSet.Exts = append(msgExSet.Exts, &ExtItem{
			IsDel:     1,
			Key:       req.Ext.Key,
			Value:     req.Ext.Value,
			Timestamp: time.Now().UnixMilli(),
			User: &UserInfo{
				UserId: userId,
			},
		})
		bs, _ := json.Marshal(msgExSet)
		upMsg := &pbobjs.UpMsg{
			MsgType:    MsgExtCmdType,
			MsgContent: bs,
			Flags:      msgdefines.SetStateMsg(0),
		}
		if req.ChannelType == pbobjs.ChannelType_Private {
			commonservices.AsyncPrivateMsg(ctx, userId, req.TargetId, upMsg)
		} else if req.ChannelType == pbobjs.ChannelType_Group {
			commonservices.AsyncGroupMsg(ctx, userId, req.TargetId, upMsg)
		}
	}
	return errs.IMErrorCode_SUCCESS
}

func QryMsgExts(ctx context.Context, req *pbobjs.QryMsgExtReq) (errs.IMErrorCode, *pbobjs.MsgExtItemsList) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	userId := bases.GetRequesterIdFromCtx(ctx)
	converId := commonservices.GetConversationId(userId, req.TargetId, req.ChannelType)
	msgMap := map[string][]byte{}
	if req.ChannelType == pbobjs.ChannelType_Private {
		storage := storages.NewPrivateHisMsgStorage()
		msgs, err := storage.FindByIds(appkey, converId, req.MsgIds, 0)
		if err == nil {
			for _, msg := range msgs {
				msgMap[msg.MsgId] = msg.MsgExt
			}
		}
	} else if req.ChannelType == pbobjs.ChannelType_Group {
		storage := storages.NewGroupHisMsgStorage()
		msgs, err := storage.FindByIds(appkey, converId, req.MsgIds, 0)
		if err == nil {
			for _, msg := range msgs {
				msgMap[msg.MsgId] = msg.MsgExt
			}
		}
	}
	ret := &pbobjs.MsgExtItemsList{
		Items: []*pbobjs.MsgExtItems{},
	}
	for msgId, extBs := range msgMap {
		msgExtItems := &pbobjs.MsgExtItems{
			MsgId: msgId,
			Exts:  []*pbobjs.MsgExtItem{},
		}
		if len(extBs) > 0 {
			tools.PbUnMarshal(extBs, msgExtItems)
		}
		ret.Items = append(ret.Items, msgExtItems)
	}
	return errs.IMErrorCode_SUCCESS, ret
}

func AddMsgExSet(ctx context.Context, req *pbobjs.MsgExt) errs.IMErrorCode {
	if req.Ext == nil || req.Ext.Key == "" {
		return errs.IMErrorCode_SUCCESS
	}
	appkey := bases.GetAppKeyFromCtx(ctx)
	msgId := req.MsgId
	userId := bases.GetRequesterIdFromCtx(ctx)

	converId := commonservices.GetConversationId(userId, req.TargetId, req.ChannelType)
	msgInfo := GetMsgInfo(appkey, converId, req.MsgId, req.ChannelType)
	optTime := time.Now().UnixMilli()
	code := msgInfo.AddMsgExset(&pbobjs.MsgExtItem{
		Key:       req.Ext.Key,
		Value:     req.Ext.Value,
		Timestamp: optTime,
		UserInfo: &pbobjs.UserInfo{
			UserId: userId,
		},
	})
	if code != errs.IMErrorCode_SUCCESS {
		return code
	}
	extItems := &pbobjs.MsgExtItems{
		MsgId: req.MsgId,
		Exts:  []*pbobjs.MsgExtItem{},
	}
	msgInfo.ForeachMsgExset(func(key string, exts []*pbobjs.MsgExtItem) {
		for _, ext := range exts {
			extItem := &pbobjs.MsgExtItem{
				Key:       ext.Key,
				Value:     ext.Value,
				Timestamp: ext.Timestamp,
			}
			if ext.UserInfo != nil {
				extItem.UserInfo = &pbobjs.UserInfo{
					UserId: ext.UserInfo.UserId,
				}
			}
			extItems.Exts = append(extItems.Exts, extItem)
		}
	})
	extItemsBs, _ := tools.PbMarshal(extItems)
	if req.ChannelType == pbobjs.ChannelType_Private {
		storage := storages.NewPrivateHisMsgStorage()
		storage.UpdateMsgExset(appkey, converId, msgId, extItemsBs)
	} else if req.ChannelType == pbobjs.ChannelType_Group {
		storage := storages.NewGroupHisMsgStorage()
		storage.UpdateMsgExset(appkey, converId, msgId, extItemsBs)
	}
	msgExSet := &MsgExt{
		MsgId: msgId,
		Exts:  []*ExtItem{},
	}
	uInfo := commonservices.GetTargetDisplayUserInfo(ctx, userId)
	msgUserInfo := &UserInfo{
		UserId: userId,
	}
	if uInfo != nil {
		msgUserInfo.Nickname = uInfo.Nickname
		msgUserInfo.UserPortrait = uInfo.UserPortrait
	}
	msgExSet.Exts = append(msgExSet.Exts, &ExtItem{
		Key:       req.Ext.Key,
		Value:     req.Ext.Value,
		Timestamp: optTime,
		User:      msgUserInfo,
	})
	bs, _ := json.Marshal(msgExSet)
	upMsg := &pbobjs.UpMsg{
		MsgType:    MsgExSetCmdType,
		MsgContent: bs,
		Flags:      msgdefines.SetStateMsg(0),
	}
	if req.ChannelType == pbobjs.ChannelType_Private {
		commonservices.AsyncPrivateMsg(ctx, userId, req.TargetId, upMsg)
	} else if req.ChannelType == pbobjs.ChannelType_Group {
		commonservices.AsyncGroupMsg(ctx, userId, req.TargetId, upMsg)
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
	converId := commonservices.GetConversationId(userId, req.TargetId, req.ChannelType)

	msgInfo := GetMsgInfo(appkey, converId, msgId, req.ChannelType)
	succ := msgInfo.DelMsgExset(req.Ext.Key, req.Ext.Value)
	if succ {
		extItems := &pbobjs.MsgExtItems{
			MsgId: req.MsgId,
			Exts:  []*pbobjs.MsgExtItem{},
		}
		msgInfo.ForeachMsgExset(func(key string, exts []*pbobjs.MsgExtItem) {
			for _, ext := range exts {
				extItem := &pbobjs.MsgExtItem{
					Key:       ext.Key,
					Value:     ext.Value,
					Timestamp: ext.Timestamp,
				}
				if ext.UserInfo != nil {
					extItem.UserInfo = &pbobjs.UserInfo{
						UserId: ext.UserInfo.UserId,
					}
				}
				extItems.Exts = append(extItems.Exts, extItem)
			}
		})
		extItemsBs, _ := tools.PbMarshal(extItems)
		if req.ChannelType == pbobjs.ChannelType_Private {
			storage := storages.NewPrivateHisMsgStorage()
			storage.UpdateMsgExset(appkey, converId, req.MsgId, extItemsBs)
		} else if req.ChannelType == pbobjs.ChannelType_Group {
			storage := storages.NewGroupHisMsgStorage()
			storage.UpdateMsgExset(appkey, converId, req.MsgId, extItemsBs)
		}
		msgExSet := &MsgExt{
			MsgId: msgId,
			Exts:  []*ExtItem{},
		}
		uInfo := commonservices.GetTargetDisplayUserInfo(ctx, userId)
		msgUserInfo := &UserInfo{
			UserId: userId,
		}
		if uInfo != nil {
			msgUserInfo.Nickname = uInfo.Nickname
			msgUserInfo.UserPortrait = uInfo.UserPortrait
		}
		msgExSet.Exts = append(msgExSet.Exts, &ExtItem{
			IsDel:     1,
			Key:       req.Ext.Key,
			Value:     req.Ext.Value,
			Timestamp: time.Now().UnixMilli(),
			User:      msgUserInfo,
		})
		bs, _ := json.Marshal(msgExSet)
		upMsg := &pbobjs.UpMsg{
			MsgType:    MsgExSetCmdType,
			MsgContent: bs,
			Flags:      msgdefines.SetStateMsg(0),
		}
		if req.ChannelType == pbobjs.ChannelType_Private {
			commonservices.AsyncPrivateMsg(ctx, userId, req.TargetId, upMsg)
		} else if req.ChannelType == pbobjs.ChannelType_Group {
			commonservices.AsyncGroupMsg(ctx, userId, req.TargetId, upMsg)
		}
	} else {
		return errs.IMErrorCode_MSG_MSGEXTDUPLICATE
	}
	return errs.IMErrorCode_SUCCESS
}

func QryMsgExSets(ctx context.Context, req *pbobjs.QryMsgExtReq) (errs.IMErrorCode, *pbobjs.MsgExtItemsList) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	userId := bases.GetRequesterIdFromCtx(ctx)
	converId := commonservices.GetConversationId(userId, req.TargetId, req.ChannelType)
	msgMap := map[string][]byte{}
	if req.ChannelType == pbobjs.ChannelType_Private {
		storage := storages.NewPrivateHisMsgStorage()
		msgs, err := storage.FindByIds(appkey, converId, req.MsgIds, 0)
		if err == nil {
			for _, msg := range msgs {
				msgMap[msg.MsgId] = msg.MsgExset
			}
		}
	} else if req.ChannelType == pbobjs.ChannelType_Group {
		storage := storages.NewGroupHisMsgStorage()
		msgs, err := storage.FindByIds(appkey, converId, req.MsgIds, 0)
		if err == nil {
			for _, msg := range msgs {
				msgMap[msg.MsgId] = msg.MsgExset
			}
		}
	}
	ret := &pbobjs.MsgExtItemsList{
		Items: []*pbobjs.MsgExtItems{},
	}
	for msgId, extBs := range msgMap {
		msgExtItems := &pbobjs.MsgExtItems{
			MsgId: msgId,
			Exts:  []*pbobjs.MsgExtItem{},
		}
		if len(extBs) > 0 {
			tools.PbUnMarshal(extBs, msgExtItems)
		}
		ret.Items = append(ret.Items, msgExtItems)
	}
	//fill userinfo
	fillUserInfos(ctx, ret.Items)
	return errs.IMErrorCode_SUCCESS, ret
}

var MsgExtCmdType string = msgdefines.CmdMsgType_MsgExt
var MsgExSetCmdType string = msgdefines.CmdMsgType_MsgExSet

type MsgExt struct {
	MsgId string     `json:"msg_id"`
	Exts  []*ExtItem `json:"exts"`
}

type ExtItem struct {
	IsDel     int       `json:"is_del"`
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	Timestamp int64     `json:"timestamp"`
	User      *UserInfo `json:"user,omitempty"`
}

type UserInfo struct {
	UserId       string `json:"user_id"`
	Nickname     string `json:"nickname"`
	UserPortrait string `json:"user_portrait"`
}

func fillUserInfos(ctx context.Context, extItemsList []*pbobjs.MsgExtItems) {
	userIdsMap := map[string]bool{}
	userIds := []string{}
	for _, msgExtItems := range extItemsList {
		for _, msgExtItem := range msgExtItems.Exts {
			if msgExtItem.UserInfo != nil {
				uId := msgExtItem.UserInfo.UserId
				if uId != "" {
					if _, exist := userIdsMap[uId]; !exist {
						userIdsMap[uId] = true
						userIds = append(userIds, uId)
					}
				}
			}
		}
	}
	if len(userIds) > 0 {
		userInfoMap := commonservices.GetTargetDisplayUserInfosMap(ctx, userIds)
		for _, msgExtItems := range extItemsList {
			for _, msgExtItem := range msgExtItems.Exts {
				if msgExtItem.UserInfo != nil {
					uId := msgExtItem.UserInfo.UserId
					if uId != "" {
						if userInfo, exist := userInfoMap[uId]; exist {
							msgExtItem.UserInfo.Nickname = userInfo.Nickname
							msgExtItem.UserInfo.UserPortrait = userInfo.UserPortrait
						}
					}
				}
			}
		}
	}
}

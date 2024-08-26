package services

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/chatroom/storages"
	"im-server/services/chatroom/storages/models"
	"im-server/services/commonservices"
)

func QryChatAtts(ctx context.Context, req *pbobjs.ChatroomInfo) (errs.IMErrorCode, *pbobjs.ChatAtts) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	chatId := req.ChatId
	container, exist := getChatroomContainer(appkey, chatId)
	if !exist {
		return errs.IMErrorCode_CHATROOM_NOTEXIST, nil
	}
	lock := chatroomLocks.GetLocks(appkey, chatId)
	lock.RLock()
	defer lock.RUnlock()
	resp := &pbobjs.ChatAtts{
		Atts:       []*pbobjs.ChatAttItem{},
		IsComplete: true,
		IsFinished: true,
	}
	atts := container.Atts
	for k, v := range atts {
		resp.Atts = append(resp.Atts, &pbobjs.ChatAttItem{
			Key:     k,
			Value:   v.Value,
			UserId:  v.UserId,
			AttTime: v.AttTime,
			OptType: v.OptType,
		})
	}
	return errs.IMErrorCode_SUCCESS, resp
}

func getMaxAttCount(appkey string) int {
	count := 100
	if appinfo, exist := commonservices.GetAppInfo(appkey); exist {
		count = appinfo.ChrmAttMaxCount
	}
	return count
}

func BatchAddChatAtt(ctx context.Context, req *pbobjs.ChatAttBatchReq) (errs.IMErrorCode, *pbobjs.ChatAttBatchResp) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	chatId := bases.GetTargetIdFromCtx(ctx)
	userId := bases.GetRequesterIdFromCtx(ctx)
	//add to cache
	container, exist := getChatroomContainer(appkey, chatId)
	if !exist {
		return errs.IMErrorCode_CHATROOM_NOTEXIST, nil
	}
	isFromApi := bases.GetIsFromApiFromCtx(ctx)
	if !isFromApi && !container.CheckMemberExist(userId) {
		return errs.IMErrorCode_CHATROOM_NOTMEMBER, nil
	}
	resp := &pbobjs.ChatAttBatchResp{
		AttResps: []*pbobjs.ChatAttResp{},
	}
	broadReq := &pbobjs.ChatAtts{
		ChatId: chatId,
		Atts:   []*pbobjs.ChatAttItem{},
	}
	retCode := errs.IMErrorCode_SUCCESS
	storage := storages.NewChatroomExtStorage()
	for _, attReq := range req.Atts {
		code, attTime := container.AddAtt(userId, attReq.Key, attReq.Value, attReq.IsForce, attReq.IsAutoDel)
		resp.AttResps = append(resp.AttResps, &pbobjs.ChatAttResp{
			Key:     attReq.Key,
			Code:    int32(code),
			AttTime: attTime,
		})
		if code == errs.IMErrorCode_SUCCESS {
			//save to db
			storage.Upsert(models.ChatroomExt{
				ChatId:    chatId,
				ItemKey:   attReq.Key,
				ItemValue: attReq.Value,
				ItemType:  0,
				ItemTime:  attTime,
				AppKey:    appkey,
				MemberId:  userId,
			})
			//add to broadcast msg
			broadReq.Atts = append(broadReq.Atts, &pbobjs.ChatAttItem{
				Key:     attReq.Key,
				Value:   attReq.Value,
				UserId:  userId,
				AttTime: attTime,
				OptType: pbobjs.ChatAttOptType_ChatAttOpt_Add,
			})
		}
	}
	//broadcast to all of nodes
	if len(broadReq.Atts) > 0 {
		bases.Broadcast(ctx, "c_atts_dispatch", broadReq)
	}
	return retCode, resp
}

func AddChatAtt(ctx context.Context, req *pbobjs.ChatAttReq) (errs.IMErrorCode, *pbobjs.ChatAttResp) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	chatId := bases.GetTargetIdFromCtx(ctx)
	userId := bases.GetRequesterIdFromCtx(ctx)
	//add to cache
	container, exist := getChatroomContainer(appkey, chatId)
	if !exist {
		return errs.IMErrorCode_CHATROOM_NOTEXIST, &pbobjs.ChatAttResp{
			Key:  req.Key,
			Code: int32(errs.IMErrorCode_CHATROOM_NOTEXIST),
		}
	}
	isFromApi := bases.GetIsFromApiFromCtx(ctx)
	if !isFromApi && !container.CheckMemberExist(userId) {
		return errs.IMErrorCode_CHATROOM_NOTMEMBER, &pbobjs.ChatAttResp{
			Key:  req.Key,
			Code: int32(errs.IMErrorCode_CHATROOM_NOTMEMBER),
		}
	}
	code, attTime := container.AddAtt(userId, req.Key, req.Value, req.IsForce, req.IsAutoDel)
	if code != errs.IMErrorCode_SUCCESS {
		return code, &pbobjs.ChatAttResp{
			Key:  req.Key,
			Code: int32(code),
		}
	}
	//save to db
	storage := storages.NewChatroomExtStorage()
	storage.Upsert(models.ChatroomExt{
		ChatId:    chatId,
		ItemKey:   req.Key,
		ItemValue: req.Value,
		ItemType:  0,
		ItemTime:  attTime,
		AppKey:    appkey,
		MemberId:  userId,
	})
	//broadcast to all of nodes
	bases.Broadcast(ctx, "c_atts_dispatch", &pbobjs.ChatAtts{
		ChatId: chatId,
		Atts: []*pbobjs.ChatAttItem{
			{
				Key:     req.Key,
				Value:   req.Value,
				UserId:  userId,
				AttTime: attTime,
				OptType: pbobjs.ChatAttOptType_ChatAttOpt_Add,
			},
		},
	})
	resp := &pbobjs.ChatAttResp{
		Key:     req.Key,
		Code:    int32(errs.IMErrorCode_SUCCESS),
		AttTime: attTime,
	}
	//send msg
	if req.Msg != nil {
		msgCode, msgId, msgTime, msgSeq := SendChatroomMsg(ctx, req.Msg)
		resp.MsgCode = int32(msgCode)
		resp.MsgId = msgId
		resp.MsgTime = msgTime
		resp.MsgSeq = msgSeq

	}
	return errs.IMErrorCode_SUCCESS, resp
}

func BatchDelChatAtt(ctx context.Context, req *pbobjs.ChatAttBatchReq) (errs.IMErrorCode, *pbobjs.ChatAttBatchResp) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	chatId := bases.GetTargetIdFromCtx(ctx)
	userId := bases.GetRequesterIdFromCtx(ctx)
	//del from cache
	container, exist := getChatroomContainer(appkey, chatId)
	if !exist {
		return errs.IMErrorCode_CHATROOM_NOTEXIST, nil
	}
	//check member
	isFromApi := bases.GetIsFromApiFromCtx(ctx)
	if !isFromApi && !container.CheckMemberExist(userId) {
		return errs.IMErrorCode_CHATROOM_NOTMEMBER, nil
	}
	resp := &pbobjs.ChatAttBatchResp{
		AttResps: []*pbobjs.ChatAttResp{},
	}
	broadReq := &pbobjs.ChatAtts{
		ChatId: chatId,
		Atts:   []*pbobjs.ChatAttItem{},
	}
	retCode := errs.IMErrorCode_SUCCESS
	storage := storages.NewChatroomExtStorage()
	for _, attReq := range req.Atts {
		code, attTime := container.DelAtt(userId, attReq.Key, attReq.IsForce)
		resp.AttResps = append(resp.AttResps, &pbobjs.ChatAttResp{
			Key:     attReq.Key,
			Code:    int32(code),
			AttTime: attTime,
		})
		if code == errs.IMErrorCode_SUCCESS {
			//del from db
			storage.DeleteExt(appkey, chatId, attReq.Key)
			//add to broadcast msg
			broadReq.Atts = append(broadReq.Atts, &pbobjs.ChatAttItem{
				Key:     attReq.Key,
				UserId:  userId,
				AttTime: attTime,
				OptType: pbobjs.ChatAttOptType_ChatAttOpt_Del,
			})
		}
	}
	//broadcast to all of nodes
	if len(broadReq.Atts) > 0 {
		bases.Broadcast(ctx, "c_atts_dispatch", broadReq)
	}
	return retCode, resp
}

func DelChatAtt(ctx context.Context, req *pbobjs.ChatAttReq) (errs.IMErrorCode, *pbobjs.ChatAttResp) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	chatId := bases.GetTargetIdFromCtx(ctx)
	userId := bases.GetRequesterIdFromCtx(ctx)
	//del from cache
	container, exist := getChatroomContainer(appkey, chatId)
	if !exist {
		return errs.IMErrorCode_CHATROOM_NOTEXIST, nil
	}
	//check member
	isFromApi := bases.GetIsFromApiFromCtx(ctx)
	if !isFromApi && !container.CheckMemberExist(userId) {
		return errs.IMErrorCode_CHATROOM_NOTMEMBER, &pbobjs.ChatAttResp{
			Key:  req.Key,
			Code: int32(errs.IMErrorCode_CHATROOM_NOTMEMBER),
		}
	}
	code, attTime := container.DelAtt(userId, req.Key, req.IsForce)
	if code != errs.IMErrorCode_SUCCESS {
		return code, nil
	}
	//del from db
	storage := storages.NewChatroomExtStorage()
	storage.DeleteExt(appkey, chatId, req.Key)
	//broadcast to all of nodes
	bases.Broadcast(ctx, "c_atts_dispatch", &pbobjs.ChatAtts{
		ChatId: chatId,
		Atts: []*pbobjs.ChatAttItem{
			{
				Key:     req.Key,
				UserId:  userId,
				AttTime: attTime,
				OptType: pbobjs.ChatAttOptType_ChatAttOpt_Del,
			},
		},
	})
	resp := &pbobjs.ChatAttResp{
		Key:     req.Key,
		Code:    int32(errs.IMErrorCode_SUCCESS),
		AttTime: attTime,
	}
	//send msg
	if req.Msg != nil {
		msgCode, msgId, msgTime, msgSeq := SendChatroomMsg(ctx, req.Msg)
		resp.MsgCode = int32(msgCode)
		resp.MsgId = msgId
		resp.MsgTime = msgTime
		resp.MsgSeq = msgSeq

	}
	return errs.IMErrorCode_SUCCESS, resp
}

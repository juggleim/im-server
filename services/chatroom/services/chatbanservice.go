package services

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/chatroom/storages"
	"im-server/services/chatroom/storages/models"
	"time"
)

func foreachBanUsersFromDb(appkey, chatId string, banType pbobjs.ChrmBanType, f func(banUser *models.ChatroomBanUser)) {
	banUserStorage := storages.NewChatroomBanUserStorage()
	var startId int64 = 0
	for {
		banUsers, err := banUserStorage.QryBanUsers(appkey, chatId, banType, startId, 1000)
		if err != nil {
			break
		}
		if len(banUsers) > 0 {
			for _, banUser := range banUsers {
				f(banUser)
			}
		}
		if len(banUsers) < 1000 {
			break
		}
	}
}

func HandleBanUsers(ctx context.Context, req *pbobjs.BatchBanUserReq) errs.IMErrorCode {
	appkey := bases.GetAppKeyFromCtx(ctx)
	chatId := req.ChatId
	container, exist := getChatroomContainer(appkey, chatId)
	if !exist {
		return errs.IMErrorCode_CHATROOM_NOTEXIST
	}
	storage := storages.NewChatroomBanUserStorage()
	if req.BanType == pbobjs.ChrmBanType_Ban {
		for _, memberId := range req.MemberIds {
			if req.IsDelete { //delete
				if _, exist := container.BanUsers.Load(memberId); exist {
					//delete from cache
					container.BanUsers.Delete(memberId)
					//delete from db
					storage.DelBanUser(appkey, chatId, memberId, pbobjs.ChrmBanType_Ban)
				}
			} else {
				if _, exist := container.BanUsers.Load(memberId); !exist {
					cur := time.Now().UnixMilli()
					//add to cache
					container.BanUsers.Store(memberId, cur)
					//add to db
					storage.Create(models.ChatroomBanUser{
						ChatId:      chatId,
						BanType:     pbobjs.ChrmBanType_Ban,
						MemberId:    memberId,
						AppKey:      appkey,
						CreatedTime: cur,
					})
					//kick from chatroom
					QuitChatroom(ctx, memberId, &pbobjs.ChatroomInfo{
						ChatId: chatId,
					})
				}
			}
		}
	} else if req.BanType == pbobjs.ChrmBanType_Mute {
		for _, memberId := range req.MemberIds {
			if req.IsDelete { //delete
				if _, exist := container.MuteUsers.Load(memberId); exist {
					//delete from cache
					container.MuteUsers.Delete(memberId)
					//delete from db
					storage.DelBanUser(appkey, chatId, memberId, pbobjs.ChrmBanType_Mute)
				}
			} else {
				if _, exist := container.MuteUsers.Load(memberId); !exist {
					cur := time.Now().UnixMilli()
					//add to cache
					container.MuteUsers.Store(memberId, cur)
					//add to db
					storage.Create(models.ChatroomBanUser{
						ChatId:      chatId,
						BanType:     pbobjs.ChrmBanType_Mute,
						MemberId:    memberId,
						AppKey:      appkey,
						CreatedTime: cur,
					})
				}
			}
		}
	} else if req.BanType == pbobjs.ChrmBanType_Allow {
		for _, memberId := range req.MemberIds {
			if req.IsDelete {
				if _, exist := container.AllowUsers.Load(memberId); exist {
					//delete from cache
					container.AllowUsers.Delete(memberId)
					//delete from db
					storage.DelBanUser(appkey, chatId, memberId, pbobjs.ChrmBanType_Allow)
				}
			} else {
				if _, exist := container.AllowUsers.Load(memberId); !exist {
					cur := time.Now().UnixMilli()
					//add to cache
					container.AllowUsers.Store(memberId, cur)
					//add to db
					storage.Create(models.ChatroomBanUser{
						ChatId:      chatId,
						BanType:     pbobjs.ChrmBanType_Allow,
						MemberId:    memberId,
						AppKey:      appkey,
						CreatedTime: cur,
					})
				}
			}
		}
	}
	return errs.IMErrorCode_SUCCESS
}

func QryBanUsers(ctx context.Context, req *pbobjs.QryChrmBanUsersReq) (errs.IMErrorCode, *pbobjs.QryChrmBanUsersResp) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	storage := storages.NewChatroomBanUserStorage()
	var startId int64 = 0
	if req.Offset != "" {
		offsetInt, err := tools.DecodeInt(req.Offset)
		if err == nil {
			startId = offsetInt
		}
	}
	ret := &pbobjs.QryChrmBanUsersResp{
		ChatId:  req.ChatId,
		BanType: req.BanType,
		Members: []*pbobjs.ChrmBanMember{},
	}
	var maxId int64 = 0
	banUsers, err := storage.QryBanUsers(appkey, req.ChatId, req.BanType, startId, req.Limit)
	if err == nil {
		for _, user := range banUsers {
			if user.ID > maxId {
				maxId = user.ID
			}
			ret.Members = append(ret.Members, &pbobjs.ChrmBanMember{
				MemberId:    user.MemberId,
				CreatedTime: user.CreatedTime,
			})
		}
		if maxId > 0 {
			ret.Offset, _ = tools.EncodeInt(maxId)
		}
	}
	return errs.IMErrorCode_SUCCESS, ret
}

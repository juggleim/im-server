package services

import (
	"context"
	"encoding/json"
	"fmt"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/commonservices"
	"im-server/services/commonservices/msgdefines"
	"im-server/services/conversation/storages"
	"im-server/services/conversation/storages/models"
	hisStorages "im-server/services/historymsg/storages"
	hisModels "im-server/services/historymsg/storages/models"
	"sort"
	"time"
)

func SaveConversationV2(ctx context.Context, appkey string, userId string, msg *pbobjs.DownMsg, isOffline bool) {
	if msg == nil {
		return
	}
	if msgdefines.IsStoreMsg(msg.Flags) {
		var sortTime int64 = msg.MsgTime
		var unreadIndex int64 = 0
		if msg.IsSend {
			if msgdefines.IsNoAffectSenderConver(msg.Flags) {
				sortTime = 0
			}
		} else {
			unreadIndex = msg.UnreadIndex
		}
		if !isOffline || UserConversContains(appkey, userId) {
			userConvers := getUserConvers(appkey, userId)
			userConvers.UpsertCovner(models.Conversation{
				AppKey:      appkey,
				UserId:      userId,
				TargetId:    msg.TargetId,
				ChannelType: msg.ChannelType,
				LatestMsgId: msg.MsgId,
				// LatestMsg: msg,
				SortTime:             sortTime,
				SyncTime:             msg.MsgTime,
				LatestUnreadMsgIndex: unreadIndex,
			})
			if !msg.IsSend {
				HandleMentionedMsg(appkey, userId, msg, userConvers)
			}
			//save to db
			userConvers.PersistConver(msg.TargetId, msg.ChannelType)
		} else {
			UpsertOfflineConversation(&ConversationCacheItem{
				Appkey:      appkey,
				UserId:      userId,
				TargetId:    msg.TargetId,
				ChannelType: msg.ChannelType,
				LatestMsgId: msg.MsgId,
				SortTime:    sortTime,
				SyncTime:    msg.MsgTime,
				UnReadIndex: unreadIndex,
			})
			if !msg.IsSend {
				HandleMentionedMsg(appkey, userId, msg, nil)
			}
		}
	}
}

func SaveNilConversationV2(ctx context.Context, appkey string, userId string, targetId string, channelType pbobjs.ChannelType) (errs.IMErrorCode, *pbobjs.Conversation) {
	sortTime := time.Now().UnixMilli()
	resp := &pbobjs.Conversation{
		UserId:      userId,
		TargetId:    targetId,
		ChannelType: channelType,
		SortTime:    sortTime,
		SyncTime:    sortTime,
	}
	var grpInfo *GroupInfo
	var targetUserInfo *UserInfo
	if channelType == pbobjs.ChannelType_Private {
		resp.TargetUserInfo = commonservices.GetTargetDisplayUserInfo(ctx, targetId)
		targetUserInfo = &UserInfo{
			UserId:       resp.TargetUserInfo.UserId,
			Nickname:     resp.TargetUserInfo.Nickname,
			UserPortrait: resp.TargetUserInfo.UserPortrait,
			ExtFields:    commonservices.Kvitems2Map(resp.TargetUserInfo.ExtFields),
			UpdatedTime:  resp.TargetUserInfo.UpdatedTime,
		}
	} else if channelType == pbobjs.ChannelType_Group {
		resp.GroupInfo = commonservices.GetGroupInfoFromCache(ctx, targetId)
		grpInfo = &GroupInfo{
			GroupId:       resp.GroupInfo.GroupId,
			GroupName:     resp.GroupInfo.GroupName,
			GroupPortrait: resp.GroupInfo.GroupPortrait,
			UpdatedTime:   resp.GroupInfo.UpdatedTime,
			ExtFields:     commonservices.Kvitems2Map(resp.GroupInfo.ExtFields),
		}
	}
	userConvers := getUserConvers(appkey, userId)
	userConvers.UpsertCovner(models.Conversation{
		AppKey:      appkey,
		UserId:      userId,
		TargetId:    targetId,
		ChannelType: channelType,
		SortTime:    sortTime,
		SyncTime:    sortTime,
	})
	//notify other device
	addConver := &AddNilConver{
		Conversation: &Conversation{
			TargetId:       targetId,
			ChannelType:    int32(channelType),
			SortTime:       &sortTime,
			SyncTime:       &sortTime,
			TargetUserInfo: targetUserInfo,
			GroupInfo:      grpInfo,
		},
	}
	flag := msgdefines.SetCmdMsg(0)
	bs, _ := json.Marshal(addConver)
	commonservices.AsyncPrivateMsg(ctx, userId, userId, &pbobjs.UpMsg{
		MsgType:    addConverMsgType,
		MsgContent: bs,
		Flags:      flag,
	})
	return errs.IMErrorCode_SUCCESS, resp
}

func SyncConversationsV2(ctx context.Context, appkey, userId string, startTime int64, count int32) *pbobjs.QryConversationsResp {
	resp := &pbobjs.QryConversationsResp{
		Conversations: []*pbobjs.Conversation{},
	}
	userConvers := getUserConvers(appkey, userId)
	convers := userConvers.SyncConvers(startTime, count+1)
	if len(convers) > int(count) {
		convers = convers[:count]
	} else {
		resp.IsFinished = true
	}
	conversations := fillConvers(ctx, userId, convers, userConvers)
	resp.Conversations = append(resp.Conversations, conversations...)
	return resp
}

func QryConverV2(ctx context.Context, userId string, req *pbobjs.QryConverReq) (errs.IMErrorCode, *pbobjs.Conversation) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	userConvers := getUserConvers(appkey, userId)
	conver := userConvers.QryConver(req.TargetId, req.ChannelType)
	if conver == nil {
		return errs.IMErrorCode_SUCCESS, &pbobjs.Conversation{
			TargetId:    req.TargetId,
			ChannelType: req.ChannelType,
		}
	}
	if req.IsInner {
		converTags := []*pbobjs.ConverTag{}
		//conver tag
		tagStorage := storages.NewUserConverTagStorage()
		tags, err := tagStorage.QryTagsByConver(appkey, userId, req.TargetId, req.ChannelType)
		if err == nil {
			for _, tag := range tags {
				converTags = append(converTags, &pbobjs.ConverTag{
					Tag:     tag.Tag,
					TagName: tag.TagName,
					TagType: pbobjs.ConverTagType_UserConverTag,
				})
			}
		}
		return errs.IMErrorCode_SUCCESS, &pbobjs.Conversation{
			TargetId:          req.TargetId,
			ChannelType:       req.ChannelType,
			UndisturbType:     conver.UndisturbType,
			LatestUnreadIndex: conver.LatestUnreadMsgIndex,
			LatestReadIndex:   conver.LatestReadMsgIndex,
			LatestReadMsgId:   conver.LatestReadMsgId,
			LatestReadMsgTime: conver.LatestReadMsgTime,
			ConverTags:        converTags,
		}
	} else {
		conversations := fillConvers(ctx, userId, []*models.Conversation{conver}, userConvers)
		if len(conversations) <= 0 {
			return errs.IMErrorCode_SUCCESS, &pbobjs.Conversation{
				TargetId:    req.TargetId,
				ChannelType: req.ChannelType,
			}
		}
		return errs.IMErrorCode_SUCCESS, conversations[0]
	}
}

func BatchQryConversV2(ctx context.Context, req *pbobjs.QryConverReq) (errs.IMErrorCode, *pbobjs.QryConversationsResp) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	ret := &pbobjs.QryConversationsResp{
		Conversations: []*pbobjs.Conversation{},
	}
	for _, uId := range req.UserIds {
		userConvers := getUserConvers(appkey, uId)
		conver := userConvers.QryConver(req.TargetId, req.ChannelType)
		if conver != nil {
			ret.Conversations = append(ret.Conversations, &pbobjs.Conversation{
				UserId:            conver.UserId,
				TargetId:          conver.TargetId,
				ChannelType:       conver.ChannelType,
				LatestUnreadIndex: conver.LatestUnreadMsgIndex,
				LatestReadMsgTime: conver.LatestReadMsgTime,
				UndisturbType:     conver.UndisturbType,
			})
		}
	}
	return errs.IMErrorCode_SUCCESS, ret
}

func QryConversationsV2(ctx context.Context, req *pbobjs.QryConversationsReq) *pbobjs.QryConversationsResp {
	appkey := bases.GetAppKeyFromCtx(ctx)
	userId := bases.GetRequesterIdFromCtx(ctx)
	startTime := req.StartTime
	isPositiveOrder := false
	if req.Order == 0 {
		if startTime <= 0 {
			startTime = time.Now().UnixMilli()
		}
	} else {
		isPositiveOrder = true
	}
	count := req.Count
	resp := &pbobjs.QryConversationsResp{
		Conversations: []*pbobjs.Conversation{},
	}
	userConvers := getUserConvers(appkey, userId)
	convers := userConvers.QryConvers(startTime, count+1, isPositiveOrder, req.TargetId, req.ChannelType, req.Tag)
	if len(convers) > int(count) {
		convers = convers[:count]
	} else {
		resp.IsFinished = true
	}
	conversations := fillConvers(ctx, userId, convers, userConvers)
	resp.Conversations = append(resp.Conversations, conversations...)
	if isPositiveOrder {
		sort.Slice(resp.Conversations, func(i, j int) bool {
			return resp.Conversations[i].SortTime > resp.Conversations[j].SortTime
		})
	}
	return resp
}

func ClearTotalUnreadV2(ctx context.Context, userId string) errs.IMErrorCode {
	appkey := bases.GetAppKeyFromCtx(ctx)
	userConvers := getUserConvers(appkey, userId)
	userConvers.ClearTotalUnread()
	storage := storages.NewConversationStorage()
	storage.ClearTotalUnreadCount(appkey, userId)
	//clear mention msgs
	mentionStorage := storages.NewMentionMsgStorage()
	mentionStorage.CleanMentionMsgsBaseUserId(appkey, userId)
	commonservices.AsyncPrivateMsg(ctx, userId, userId, &pbobjs.UpMsg{
		MsgType:    msgdefines.CmdMsgType_ClearTotalUnread,
		MsgContent: []byte(fmt.Sprintf(`{"clear_time":%d}`, time.Now().UnixMilli())),
		Flags:      msgdefines.SetCmdMsg(0),
	})
	return errs.IMErrorCode_SUCCESS
}

func ClearUnreadV2(ctx context.Context, userId string, convers []*pbobjs.Conversation, noCmdMsg bool) errs.IMErrorCode {
	appkey := bases.GetAppKeyFromCtx(ctx)
	if len(convers) > 0 {
		clearUnreadConvers := &ClearUnreadConvers{
			Conversations: []*ClearUnreadConver{},
		}
		userConvers := getUserConvers(appkey, userId)
		affected := false
		for _, conver := range convers {
			var msgIndex int64 = conver.LatestReadIndex
			clearUnreadConvers.Conversations = append(clearUnreadConvers.Conversations, &ClearUnreadConver{
				TargetId:           conver.TargetId,
				ChannelType:        int32(conver.ChannelType),
				LatestReadMsgIndex: msgIndex,
			})
			if msgIndex > 0 {
				readMsgId := conver.LatestReadMsgId
				readMsgTime := conver.LatestReadMsgTime
				if conver.LatestReadMsgTime <= 0 {
					readMsgTime = time.Now().UnixMilli()
				}
				rowsAffected := userConvers.ClearUnread(conver.TargetId, conver.ChannelType, msgIndex, readMsgId, readMsgTime)
				if rowsAffected {
					affected = true
					userConvers.PersistConver(conver.TargetId, conver.ChannelType)
				}
			} else {
				//TODO
				rowsAffected := userConvers.DefaultClearUnread(conver.TargetId, conver.ChannelType)
				if rowsAffected {
					affected = true
					userConvers.PersistConver(conver.TargetId, conver.ChannelType)
				}
			}
		}
		if affected && !noCmdMsg {
			//Notify other device to clear unread
			flag := msgdefines.SetCmdMsg(0)
			bs, _ := json.Marshal(clearUnreadConvers)
			exts := bases.GetExtsFromCtx(ctx)
			if len(convers) == 1 {
				exts[commonservices.RpcExtKey_UniqTag] = fmt.Sprintf("%s_%d", convers[0].TargetId, convers[0].ChannelType)
			}
			commonservices.AsyncPrivateMsg(ctx, userId, userId, &pbobjs.UpMsg{
				MsgType:    msgdefines.CmdMsgType_ClearUnread,
				MsgContent: bs,
				Flags:      flag,
			}, &bases.ExtsOption{Exts: exts})
		}
	}
	return errs.IMErrorCode_SUCCESS
}

func DelConversationV2(ctx context.Context, userId string, convers []*pbobjs.Conversation) errs.IMErrorCode {
	appkey := bases.GetAppKeyFromCtx(ctx)
	if len(convers) > 0 {
		delConvers := &DelConvers{
			Conversations: []*Conversation{},
		}
		userConvers := getUserConvers(appkey, userId)
		for _, conver := range convers {
			delConvers.Conversations = append(delConvers.Conversations, &Conversation{
				TargetId:    conver.TargetId,
				ChannelType: int32(conver.ChannelType),
			})
			affected := userConvers.DelConversation(conver.TargetId, conver.ChannelType)
			if affected {
				userConvers.PersistConver(conver.TargetId, conver.ChannelType)
			}
		}
		// bases.AsyncRpcCall(ctx, "del_conver_cache", userId, &pbobjs.ConversationsReq{
		// 	Conversations: convers,
		// })
		//notify other device
		flag := msgdefines.SetCmdMsg(0)
		bs, _ := json.Marshal(delConvers)
		commonservices.AsyncPrivateMsg(ctx, userId, userId, &pbobjs.UpMsg{
			MsgType:    delConversMsgType,
			MsgContent: bs,
			Flags:      flag,
		})
	}
	return errs.IMErrorCode_SUCCESS
}

func MarkUnreadV2(ctx context.Context, userId string, req *pbobjs.ConversationsReq) errs.IMErrorCode {
	appkey := bases.GetAppKeyFromCtx(ctx)
	if len(req.Conversations) > 0 {
		markUnreadConvers := &MarkUnreadConvers{
			Conversations: []*MarkUnreadConver{},
		}
		userConvers := getUserConvers(appkey, userId)
		affected := false
		for _, conver := range req.Conversations {
			markUnreadConvers.Conversations = append(markUnreadConvers.Conversations, &MarkUnreadConver{
				TargetId:    conver.TargetId,
				ChannelType: int32(conver.ChannelType),
				UnreadTag:   int(conver.UnreadTag),
			})
			rowsAffected := userConvers.UpdateUnreadTag(conver.TargetId, conver.ChannelType, int(conver.UnreadTag))
			if rowsAffected {
				userConvers.PersistConver(conver.TargetId, conver.ChannelType)
				affected = true
			}
		}
		if affected {
			//Notify other device to mark unread
			flag := msgdefines.SetCmdMsg(0)
			bs, _ := json.Marshal(markUnreadConvers)
			commonservices.AsyncPrivateMsg(ctx, userId, userId, &pbobjs.UpMsg{
				MsgType:    msgdefines.CmdMsgType_MarkUnread,
				MsgContent: bs,
				Flags:      flag,
			})
		}
	}
	return errs.IMErrorCode_SUCCESS
}

func SetTopConversV2(ctx context.Context, req *pbobjs.ConversationsReq) (errs.IMErrorCode, int64) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	userId := bases.GetRequesterIdFromCtx(ctx)
	topConvers := &TopConvers{
		Conversations: []*Conversation{},
	}
	var cmdMsgTime int64 = 0
	currTime := time.Now().UnixMilli()
	if len(req.Conversations) > 0 {
		userConvers := getUserConvers(appkey, userId)
		for _, conver := range req.Conversations {
			var t int64 = 0
			if conver.IsTop > 0 {
				t = currTime
			}
			affected := userConvers.UpdTopState(conver.TargetId, conver.ChannelType, int(conver.IsTop), t)
			if affected {
				userConvers.PersistConver(conver.TargetId, conver.ChannelType)
			}
			topConvers.Conversations = append(topConvers.Conversations, &Conversation{
				TargetId:      conver.TargetId,
				ChannelType:   int32(conver.ChannelType),
				IsTop:         conver.IsTop,
				TopUpdateTime: t,
			})
		}
		//notify other device
		flag := msgdefines.SetCmdMsg(0)
		bs, _ := json.Marshal(topConvers)
		code, _, msgTime, _ := commonservices.SyncPrivateMsg(ctx, userId, userId, &pbobjs.UpMsg{
			MsgType:    topConversMsgType,
			MsgContent: bs,
			Flags:      flag,
		})
		if code == errs.IMErrorCode_SUCCESS {
			cmdMsgTime = msgTime
		}
	}
	return errs.IMErrorCode_SUCCESS, cmdMsgTime
}

func QryTopConversV2(ctx context.Context, userId string, req *pbobjs.QryTopConversReq) (errs.IMErrorCode, *pbobjs.QryConversationsResp) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	resp := &pbobjs.QryConversationsResp{
		Conversations: []*pbobjs.Conversation{},
		IsFinished:    true,
	}
	userConvers := getUserConvers(appkey, userId)
	convers := userConvers.QryTopConvers(req.StartTime, 101, req.SortType, req.Order > 0)
	conversations := fillConvers(ctx, userId, convers, userConvers)
	resp.Conversations = append(resp.Conversations, conversations...)
	if len(resp.Conversations) > 100 {
		resp.Conversations = resp.Conversations[:100]
		resp.IsFinished = false
	}
	return errs.IMErrorCode_SUCCESS, resp
}

func QryTotalUnreadCountV2(ctx context.Context, userId string, req *pbobjs.QryTotalUnreadCountReq) *pbobjs.QryTotalUnreadCountResp {
	appkey := bases.GetAppKeyFromCtx(ctx)
	userConvers := getUserConvers(appkey, userId)
	return &pbobjs.QryTotalUnreadCountResp{
		TotalCount: userConvers.TotalUnreadCount(),
	}
}

func UndisturbConversV2(ctx context.Context, req *pbobjs.UndisturbConversReq) errs.IMErrorCode {
	appkey := bases.GetAppKeyFromCtx(ctx)
	userId := bases.GetTargetIdFromCtx(ctx)
	convers := &UndisturbConvers{
		Conversations: []*UndisturbConver{},
	}
	userConvers := getUserConvers(appkey, userId)
	needUnCacheConvers := []*pbobjs.Conversation{}
	for _, item := range req.Items {
		affected := userConvers.UpdateUndisturbType(item.TargetId, item.ChannelType, item.UndisturbType)
		if affected {
			userConvers.PersistConver(item.TargetId, item.ChannelType)
			needUnCacheConvers = append(needUnCacheConvers, &pbobjs.Conversation{
				TargetId:    item.TargetId,
				ChannelType: item.ChannelType,
			})
		}
		convers.Conversations = append(convers.Conversations, &UndisturbConver{
			TargetId:      item.TargetId,
			ChannelType:   int32(item.ChannelType),
			UndisturbType: item.UndisturbType,
		})
	}
	if len(needUnCacheConvers) > 0 {
		bases.AsyncRpcCall(ctx, "del_conver_cache", userId, &pbobjs.ConversationsReq{
			Conversations: needUnCacheConvers,
		})
	}
	//Notify other device to update undisturb
	flag := msgdefines.SetCmdMsg(0)
	bs, _ := json.Marshal(convers)
	commonservices.AsyncPrivateMsg(ctx, userId, userId, &pbobjs.UpMsg{
		MsgType:    undisturbMsgType,
		MsgContent: bs,
		Flags:      flag,
	})
	return errs.IMErrorCode_SUCCESS
}

func fillConvers(ctx context.Context, userId string, convers []*models.Conversation, userConvers *UserConversations) []*pbobjs.Conversation {
	appkey := bases.GetAppKeyFromCtx(ctx)
	retConvers := []*pbobjs.Conversation{}

	priConvers := []hisModels.ConverItem{}
	grpConvers := []hisModels.ConverItem{}
	for _, conver := range convers {
		if conver.IsDeleted > 0 {
			continue
		}
		if conver.ChannelType == pbobjs.ChannelType_Private {
			priConvers = append(priConvers, hisModels.ConverItem{
				ConverId: commonservices.GetConversationId(userId, conver.TargetId, conver.ChannelType),
				MsgId:    conver.LatestMsgId,
			})
		} else if conver.ChannelType == pbobjs.ChannelType_Group {
			grpConvers = append(grpConvers, hisModels.ConverItem{
				ConverId: commonservices.GetConversationId(userId, conver.TargetId, conver.ChannelType),
				MsgId:    conver.LatestMsgId,
			})
		}
	}
	priMsgs := qryPriMsgs(appkey, userId, priConvers)
	grpMsgs := qryGrpMsgs(appkey, userId, grpConvers)
	for _, conver := range convers {
		unreadCount := conver.LatestUnreadMsgIndex - conver.LatestReadMsgIndex
		conversation := &pbobjs.Conversation{
			UserId:            conver.UserId,
			TargetId:          conver.TargetId,
			ChannelType:       conver.ChannelType,
			SortTime:          conver.SortTime,
			SyncTime:          conver.SyncTime,
			UnreadCount:       unreadCount,
			IsTop:             int32(conver.IsTop),
			TopUpdatedTime:    conver.TopUpdatedTime,
			UndisturbType:     conver.UndisturbType,
			LatestUnreadIndex: conver.LatestUnreadMsgIndex,
			LatestReadIndex:   conver.LatestReadMsgIndex,
			IsDelete:          int32(conver.IsDeleted),
			UnreadTag:         int32(conver.UnreadTag),
		}
		//add latest msg
		var downMsg *pbobjs.DownMsg
		if conver.ChannelType == pbobjs.ChannelType_Private {
			if msg, exist := priMsgs[conver.LatestMsgId]; exist {
				downMsg = msg
				conversation.TargetUserInfo = msg.TargetUserInfo
				if userId == downMsg.SenderId {
					downMsg.IsSend = true
					downMsg.TargetId = conver.TargetId
					conversation.TargetUserInfo = commonservices.GetTargetDisplayUserInfo(ctx, downMsg.TargetId)
				}
			}
		} else if conver.ChannelType == pbobjs.ChannelType_Group {
			if msg, exist := grpMsgs[conver.LatestMsgId]; exist {
				downMsg = msg
				conversation.GroupInfo = msg.GroupInfo
				conversation.TargetUserInfo = msg.TargetUserInfo
			}
		}
		conversation.Msg = downMsg
		if conver.IsDeleted == 0 {
			//target userinfo/groupinfo
			// if conversation.ChannelType == pbobjs.ChannelType_Private {
			// 	conversation.TargetUserInfo = commonservices.GetTargetDisplayUserInfo(ctx, conversation.TargetId)
			// } else if conversation.ChannelType == pbobjs.ChannelType_Group {
			// 	conversation.GroupInfo = commonservices.GetGroupInfoFromCache(ctx, conversation.TargetId)
			// }
			//mentions
			mentionInfo := userConvers.GetMentionInfo(conver.TargetId, conver.ChannelType)
			if mentionInfo == nil {
				userConvers.AppendMention(conver.TargetId, conver.ChannelType, nil)
				mentionInfo = userConvers.GetMentionInfo(conver.TargetId, conver.ChannelType)
			}
			if mentionInfo != nil {
				conversation.Mentions = &pbobjs.Mentions{
					IsMentioned:     mentionInfo.IsMentioned,
					MentionMsgCount: int32(mentionInfo.MentionMsgCount),
					Senders:         []*pbobjs.UserInfo{},
					MentionMsgs:     []*pbobjs.MentionMsg{},
				}
				for _, senderId := range mentionInfo.SenderIds {
					userInfo := commonservices.GetTargetDisplayUserInfo(ctx, senderId)
					conversation.Mentions.Senders = append(conversation.Mentions.Senders, userInfo)
				}
				for _, mentionMsg := range mentionInfo.MentionMsgs {
					conversation.Mentions.MentionMsgs = append(conversation.Mentions.MentionMsgs, &pbobjs.MentionMsg{
						SenderId:    mentionMsg.SenderId,
						MsgId:       mentionMsg.MsgId,
						MsgTime:     mentionMsg.MsgTime,
						MentionType: mentionMsg.MentionType,
					})
				}
			}
			//conver tags
			if conver.ConverExts != nil && len(conver.ConverExts.ConverTags) > 0 {
				tagList := []*pbobjs.ConverTag{}
				for tag := range conver.ConverExts.ConverTags {
					tagList = append(tagList, &pbobjs.ConverTag{
						Tag: tag,
					})
				}
				conversation.ConverTags = tagList
			}
		}
		retConvers = append(retConvers, conversation)
	}
	return retConvers
}

func qryGrpMsgs(appkey, userId string, convers []hisModels.ConverItem) map[string]*pbobjs.DownMsg {
	storage := hisStorages.NewGroupHisMsgStorage()
	msgs, err := storage.FindByConvers(appkey, convers)
	ret := map[string]*pbobjs.DownMsg{}
	if err == nil {
		for _, msg := range msgs {
			var downMsg pbobjs.DownMsg
			err = tools.PbUnMarshal(msg.MsgBody, &downMsg)
			if err == nil {
				if userId == downMsg.SenderId {
					downMsg.IsSend = true
				}
				ret[downMsg.MsgId] = &downMsg
			}
		}
	}
	return ret
}

func qryPriMsgs(appkey, userId string, convers []hisModels.ConverItem) map[string]*pbobjs.DownMsg {
	storage := hisStorages.NewPrivateHisMsgStorage()
	msgs, err := storage.FindByConvers(appkey, convers)
	ret := map[string]*pbobjs.DownMsg{}
	if err == nil {
		for _, msg := range msgs {
			var downMsg pbobjs.DownMsg
			err = tools.PbUnMarshal(msg.MsgBody, &downMsg)
			if err == nil {
				if userId == msg.SenderId {
					downMsg.IsSend = true
					downMsg.TargetId = msg.ReceiverId
				}
				downMsg.IsRead = msg.IsRead > 0
				ret[downMsg.MsgId] = &downMsg
			}
		}
	}
	return ret
}

package services

import (
	"context"
	"encoding/json"
	"fmt"
	"im-server/commons/bases"
	"im-server/commons/caches"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/commonservices"
	"im-server/services/conversation/storages"
	"im-server/services/conversation/storages/models"
	"sort"
	"time"

	"google.golang.org/protobuf/proto"
)

var globalConverCache *caches.EphemeralCache

func init() {
	globalConverCache = caches.NewEphemeralCache(time.Millisecond*100, 5*time.Second, func(key, value interface{}) {
		conver, ok := value.(*GlobalConversationCacheItem)
		if ok && conver != nil {
			globalConverStorage := storages.NewGlobalConversationStorage()
			globalConverStorage.UpsertConversation(models.GlobalConver{
				ConverId:    conver.ConverId,
				SenderId:    conver.SenderId,
				TargetId:    conver.TargetId,
				ChannelType: conver.ChannelType,
				UpdatedTime: conver.MsgTime,
				AppKey:      conver.Appkey,
			})
		}
	})
}

type ConversationCacheItem struct {
	Appkey      string
	UserId      string
	TargetId    string
	ChannelType pbobjs.ChannelType

	LatestMsgId string
	LatestMsg   *pbobjs.DownMsg
	SortTime    int64
	SyncTime    int64
	UnReadIndex int64
}

func getGlobalConverCacheKey(appkey, converId string, channelType pbobjs.ChannelType) string {
	return fmt.Sprintf("%s_%s_%d", appkey, converId, channelType)
}

func saveConversationByCache(item *ConversationCacheItem) {
	if item == nil {
		return
	}
}

type GlobalConversationCacheItem struct {
	Appkey      string
	ConverId    string
	SenderId    string
	TargetId    string
	ChannelType pbobjs.ChannelType
	MsgTime     int64
}

func saveGlobalConversationByCache(item *GlobalConversationCacheItem) {
	if item == nil {
		return
	}
	key := getGlobalConverCacheKey(item.Appkey, item.ConverId, item.ChannelType)
	globalConverCache.Upsert(key, func(oldVal interface{}) interface{} {
		var converItem *GlobalConversationCacheItem
		if oldVal != nil {
			converItem = oldVal.(*GlobalConversationCacheItem)
			if item.MsgTime > converItem.MsgTime {
				converItem.MsgTime = item.MsgTime
			}
			converItem.SenderId = item.SenderId
			converItem.TargetId = item.TargetId
		} else {
			converItem = item
		}
		return converItem
	})
}

var addConverMsgType string = "jg:addconver"

type AddNilConver struct {
	Conversation *Conversation `json:"conversation"`
}

func SaveNilConversation(ctx context.Context, appkey string, userId string, targetId string, channelType pbobjs.ChannelType) (errs.IMErrorCode, *pbobjs.Conversation) {
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
	saveConversationByCache(&ConversationCacheItem{
		Appkey:      appkey,
		UserId:      userId,
		TargetId:    targetId,
		ChannelType: channelType,
		LatestMsgId: "",
		SortTime:    sortTime,
		SyncTime:    sortTime,
		UnReadIndex: 0,
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
	flag := commonservices.SetCmdMsg(0)
	bs, _ := json.Marshal(addConver)
	commonservices.AsyncPrivateMsg(ctx, userId, userId, &pbobjs.UpMsg{
		MsgType:    addConverMsgType,
		MsgContent: bs,
		Flags:      flag,
	})
	return errs.IMErrorCode_SUCCESS, resp
}

func SaveConversation(appkey string, userId string, msg *pbobjs.DownMsg) error {
	if msg == nil {
		return nil
	}
	if commonservices.IsStoreMsg(msg.Flags) {
		var sortTime int64 = msg.MsgTime
		var unreadIndex int64 = 0
		var err error
		if msg.IsSend {
			if commonservices.IsNoAffectSenderConver(msg.Flags) {
				sortTime = 0
			}
			isRecordGlobalConvers := false
			appinfo, exist := commonservices.GetAppInfo(appkey)
			if exist && appinfo != nil {
				isRecordGlobalConvers = appinfo.RecordGlobalConvers
			}
			if isRecordGlobalConvers {
				//record global conversation
				converId := GetGlobalConverId(userId, msg.TargetId, msg.ChannelType)
				saveGlobalConversationByCache(&GlobalConversationCacheItem{
					Appkey:      appkey,
					ConverId:    converId,
					SenderId:    userId,
					TargetId:    msg.TargetId,
					ChannelType: msg.ChannelType,
					MsgTime:     msg.MsgTime,
				})
			}
		} else {
			unreadIndex = msg.UnreadIndex
		}
		HandleMentionedMsg(appkey, userId, msg)
		saveConversationByCache(&ConversationCacheItem{
			Appkey:      appkey,
			UserId:      userId,
			TargetId:    msg.TargetId,
			ChannelType: msg.ChannelType,
			LatestMsgId: msg.MsgId,
			LatestMsg:   msg,
			SortTime:    sortTime,
			SyncTime:    msg.MsgTime,
			UnReadIndex: unreadIndex,
		})

		return err
	} else {
		return nil
	}
}

func HandleMentionedMsg(appkey string, userId string, msg *pbobjs.DownMsg) {
	if commonservices.IsMentionedMe(userId, msg) {
		isRead := 0
		if msg.IsSend {
			isRead = 1
		}
		//append to cache
		userConvers := getUserConvers(appkey, userId)
		userConvers.AppendMention(msg.TargetId, msg.ChannelType, &models.MentionMsg{
			SenderId: msg.SenderId,
			MsgId:    msg.MsgId,
			MsgTime:  msg.MsgTime,
			MsgIndex: msg.UnreadIndex,
		})
		//save mentioned msg
		storage := storages.NewMentionMsgStorage()
		storage.SaveMentionMsg(models.MentionMsg{
			UserId:      userId,
			TargetId:    msg.TargetId,
			ChannelType: msg.ChannelType,
			SenderId:    msg.SenderId,
			MentionType: msg.MentionInfo.MentionType,
			MsgId:       msg.MsgId,
			MsgTime:     msg.MsgTime,
			MsgIndex:    msg.UnreadIndex,
			AppKey:      appkey,
			IsRead:      isRead,
		})
	}
}

func QryTotalUnreadCount(ctx context.Context, userId string, req *pbobjs.QryTotalUnreadCountReq) *pbobjs.QryTotalUnreadCountResp {
	appkey := bases.GetAppKeyFromCtx(ctx)
	storage := storages.NewConversationStorage()

	channelTypes := []pbobjs.ChannelType{}
	excludeConvers := []*pbobjs.SimpleConversation{}
	tag := ""
	if req.Filter != nil {
		channelTypes = append(channelTypes, req.Filter.ChannelTypes...)
		excludeConvers = append(excludeConvers, req.Filter.ExcludeConvers...)
		tag = req.Filter.Tag
	}
	var totoalCount int64 = storage.TotalUnreadCount(appkey, userId, channelTypes, excludeConvers, tag)
	return &pbobjs.QryTotalUnreadCountResp{
		TotalCount: totoalCount,
	}
}

func QryGlobalConvers(ctx context.Context, req *pbobjs.QryGlobalConversReq) *pbobjs.QryGlobalConversResp {
	appkey := bases.GetAppKeyFromCtx(ctx)
	targetId := req.TargetId
	channelType := req.ChannelType
	startTime := req.Start
	isPositiveOrder := false
	if req.Order == 0 {
		if startTime <= 0 {
			startTime = time.Now().UnixMilli()
		}
	} else {
		isPositiveOrder = true
	}
	count := req.Count
	resp := &pbobjs.QryGlobalConversResp{
		Convers: []*pbobjs.GlobalConver{},
	}
	storage := storages.NewGlobalConversationStorage()
	dbConvers, err := storage.QryConversations(appkey, targetId, channelType, startTime, count, isPositiveOrder, req.ExcludeUserIds)
	if err == nil {
		for _, dbConver := range dbConvers {
			idStr, _ := tools.EncodeInt(dbConver.Id)
			resp.Convers = append(resp.Convers, &pbobjs.GlobalConver{
				Id:          idStr,
				ConverId:    dbConver.ConverId,
				SenderId:    dbConver.SenderId,
				TargetId:    dbConver.TargetId,
				ChannelType: dbConver.ChannelType,
				UpdatedTime: dbConver.UpdatedTime,
			})
		}
		if len(resp.Convers) < int(count) {
			resp.IsFinished = true
		}
	}
	return resp
}

func dbConver2Conversations(ctx context.Context, dbConver *models.Conversation) *pbobjs.Conversation {
	appkey := bases.GetAppKeyFromCtx(ctx)
	userId := dbConver.UserId

	var downMsg *pbobjs.DownMsg
	if len(dbConver.LatestMsg) > 0 {
		var dbDownMsg pbobjs.DownMsg
		err := tools.PbUnMarshal(dbConver.LatestMsg, &dbDownMsg)
		if err == nil {
			downMsg = &dbDownMsg
		}
	}
	if downMsg == nil {
		downMsgs := QryHisMsgByIds(ctx, userId, dbConver.TargetId, dbConver.ChannelType, []string{dbConver.LatestMsgId})
		if len(downMsgs) > 0 {
			downMsg = downMsgs[0]
		}
	}
	if downMsg != nil {
		downMsg.TargetUserInfo = nil
		downMsg.GroupInfo = nil
		if downMsg.ClientUid == "" {
			downMsg.ClientUid = tools.GenerateUUIDShort22()
		}
	}

	unreadCount := dbConver.LatestUnreadMsgIndex - dbConver.LatestReadMsgIndex
	conversation := &pbobjs.Conversation{
		UserId:            dbConver.UserId,
		TargetId:          dbConver.TargetId,
		ChannelType:       dbConver.ChannelType,
		SortTime:          dbConver.SortTime,
		SyncTime:          dbConver.SyncTime,
		UnreadCount:       unreadCount,
		Msg:               downMsg,
		IsTop:             int32(dbConver.IsTop),
		TopUpdatedTime:    dbConver.TopUpdatedTime,
		UndisturbType:     dbConver.UndisturbType,
		LatestUnreadIndex: dbConver.LatestUnreadMsgIndex,
		LatestReadIndex:   dbConver.LatestReadMsgIndex,
		IsDelete:          int32(dbConver.IsDeleted),
		UnreadTag:         int32(dbConver.UnreadTag),
	}
	if conversation.ChannelType == pbobjs.ChannelType_Private {
		conversation.TargetUserInfo = commonservices.GetTargetDisplayUserInfo(ctx, conversation.TargetId)
	} else if conversation.ChannelType == pbobjs.ChannelType_Group {
		conversation.GroupInfo = commonservices.GetGroupInfoFromCache(ctx, conversation.TargetId)
	}
	//mentions
	mentionStorage := storages.NewMentionMsgStorage()
	mentionMsgs, err := mentionStorage.QryMentionSenderIdsBaseIndex(appkey, userId, dbConver.TargetId, dbConver.ChannelType, dbConver.LatestReadMsgIndex, 100)
	if err == nil && len(mentionMsgs) > 0 {
		conversation.Mentions = &pbobjs.Mentions{
			IsMentioned:     true,
			MentionMsgCount: int32(len(mentionMsgs)),
			Senders:         []*pbobjs.UserInfo{},
			MentionMsgs:     []*pbobjs.MentionMsg{},
		}
		tmpMap := map[string]int{}
		for _, mentionMsg := range mentionMsgs {
			conversation.Mentions.MentionMsgs = append(conversation.Mentions.MentionMsgs, &pbobjs.MentionMsg{
				SenderId: mentionMsg.SenderId,
				MsgId:    mentionMsg.MsgId,
				MsgTime:  mentionMsg.MsgTime,
			})
			if _, exist := tmpMap[mentionMsg.SenderId]; !exist {
				tmpMap[mentionMsg.SenderId] = 1
				userInfo := commonservices.GetTargetDisplayUserInfo(ctx, mentionMsg.SenderId)
				conversation.Mentions.Senders = append(conversation.Mentions.Senders, userInfo)
			}
		}
	}
	//conver tags
	tagStorage := storages.NewUserConverTagStorage()
	tags, err := tagStorage.QryTagsByConver(appkey, userId, dbConver.TargetId, dbConver.ChannelType)
	if err == nil {
		converTags := []*pbobjs.ConverTag{}
		for _, tag := range tags {
			converTags = append(converTags, &pbobjs.ConverTag{
				Tag:     tag.Tag,
				TagName: tag.TagName,
				TagType: pbobjs.ConverTagType_UserConverTag,
			})
		}
		conversation.ConverTags = append(conversation.ConverTags, converTags...)
	}

	return conversation
}

func QryConversations(ctx context.Context, req *pbobjs.QryConversationsReq) *pbobjs.QryConversationsResp {
	appkey := bases.GetAppKeyFromCtx(ctx)
	userId := bases.GetRequesterIdFromCtx(ctx)
	targetId := req.TargetId
	channelType := req.ChannelType
	startTime := req.StartTime
	isPositiveOrder := false
	if req.Order == 0 { //0:倒序;1:正序
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
	converStorage := storages.NewConversationStorage()
	dbConvers, err := converStorage.QryConversations(appkey, userId, targetId, channelType, startTime, count+1, isPositiveOrder, req.Tag)
	if err == nil {
		if len(dbConvers) > int(count) {
			dbConvers = dbConvers[:count]
		} else {
			resp.IsFinished = true
		}
		for _, dbConver := range dbConvers {
			conversation := dbConver2Conversations(ctx, dbConver)
			resp.Conversations = append(resp.Conversations, conversation)
		}
		if isPositiveOrder {
			sort.Slice(resp.Conversations, func(i, j int) bool {
				return resp.Conversations[i].SortTime > resp.Conversations[j].SortTime
			})
		}
	}
	return resp
}

func SyncConversations(ctx context.Context, appkey, userId string, startTime int64, count int32) *pbobjs.QryConversationsResp {
	resp := &pbobjs.QryConversationsResp{
		Conversations: []*pbobjs.Conversation{},
	}
	converStorage := storages.NewConversationStorage()
	dbConvers, err := converStorage.SyncConversations(appkey, userId, startTime, count+1)
	if err == nil {
		if len(dbConvers) > int(count) {
			dbConvers = dbConvers[:count]
		} else {
			resp.IsFinished = true
		}
		for _, dbConver := range dbConvers {
			conversation := dbConver2Conversations(ctx, dbConver)
			resp.Conversations = append(resp.Conversations, conversation)
		}
	}
	return resp
}

func QryHisMsgByIds(ctx context.Context, userId, targetId string, channelType pbobjs.ChannelType, msgIds []string) []*pbobjs.DownMsg {
	if len(msgIds) > 0 {
		converId := commonservices.GetConversationId(userId, targetId, channelType)
		code, resp, err := bases.SyncRpcCall(ctx, "qry_hismsg_by_ids", converId, &pbobjs.QryHisMsgByIdsReq{
			TargetId:    targetId,
			ChannelType: channelType,
			MsgIds:      msgIds,
		}, func() proto.Message {
			return &pbobjs.DownMsgSet{}
		})
		if err == nil && code == errs.IMErrorCode_SUCCESS && resp != nil {
			msgs, ok := resp.(*pbobjs.DownMsgSet)
			if ok && len(msgs.Msgs) > 0 {
				return msgs.Msgs
			}
		}
	}
	return []*pbobjs.DownMsg{}
}

func DelConversation(ctx context.Context, userId string, convers []*pbobjs.Conversation) errs.IMErrorCode {
	retErr := errs.IMErrorCode_SUCCESS
	appkey := bases.GetAppKeyFromCtx(ctx)
	if len(convers) > 0 {
		delConvers := &DelConvers{
			Conversations: []*Conversation{},
		}
		converStorage := storages.NewConversationStorage()
		for _, conver := range convers {
			delConvers.Conversations = append(delConvers.Conversations, &Conversation{
				TargetId:    conver.TargetId,
				ChannelType: int32(conver.ChannelType),
			})
			err := converStorage.DelConversation(appkey, userId, conver.TargetId, conver.ChannelType)
			if err != nil {
				retErr = errs.IMErrorCode_MSG_DELFAILED
			}
		}
		bases.AsyncRpcCall(ctx, "del_conver_cache", userId, &pbobjs.ConversationsReq{
			Conversations: convers,
		})
		//notify other device
		flag := commonservices.SetCmdMsg(0)
		bs, _ := json.Marshal(delConvers)
		commonservices.AsyncPrivateMsg(ctx, userId, userId, &pbobjs.UpMsg{
			MsgType:    delConversMsgType,
			MsgContent: bs,
			Flags:      flag,
		})
	}
	return retErr
}

var delConversMsgType string = "jg:delconvers"

type DelConvers struct {
	Conversations []*Conversation `json:"conversations"`
}

type Conversation struct {
	TargetId       string     `json:"target_id"`
	ChannelType    int32      `json:"channel_type"`
	IsTop          int32      `json:"is_top,omitempty"`
	TopUpdateTime  int64      `json:"top_update_time"`
	SortTime       *int64     `json:"sort_time,omitempty"`
	SyncTime       *int64     `json:"sync_time,omitempty"`
	TargetUserInfo *UserInfo  `json:"target_user_info,omitempty"`
	GroupInfo      *GroupInfo `json:"group_info,omitempty"`
}

type UserInfo struct {
	UserId       string            `json:"user_id"`
	Nickname     string            `json:"nickname"`
	UserPortrait string            `json:"user_portrait"`
	ExtFields    map[string]string `json:"ext_fields"`
	UpdatedTime  int64             `json:"updated_time"`
}

type GroupInfo struct {
	GroupId       string            `json:"group_id"`
	GroupName     string            `json:"group_name"`
	GroupPortrait string            `json:"group_portrait"`
	IsMute        int               `json:"is_mute"`
	UpdatedTime   int64             `json:"updated_time"`
	ExtFields     map[string]string `json:"ext_fields"`
}

func SetTopConvers(ctx context.Context, req *pbobjs.ConversationsReq) (errs.IMErrorCode, int64) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	userId := bases.GetRequesterIdFromCtx(ctx)
	topConvers := &TopConvers{
		Conversations: []*Conversation{},
	}
	var cmdMsgTime int64 = 0
	currTime := time.Now().UnixMilli()
	if len(req.Conversations) > 0 {
		converStorage := storages.NewConversationStorage()
		for _, conver := range req.Conversations {
			var t int64 = 0
			if conver.IsTop > 0 {
				t = currTime
			}
			converStorage.UpdateIsTopState(appkey, userId, conver.TargetId, conver.ChannelType, int(conver.IsTop), t)
			topConvers.Conversations = append(topConvers.Conversations, &Conversation{
				TargetId:      conver.TargetId,
				ChannelType:   int32(conver.ChannelType),
				IsTop:         conver.IsTop,
				TopUpdateTime: t,
			})
		}
		//notify other device
		flag := commonservices.SetCmdMsg(0)
		bs, _ := json.Marshal(topConvers)
		code, _, msgTime, _ := commonservices.SyncPrivateMsg(ctx, userId, &pbobjs.UpMsg{
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

var topConversMsgType string = "jg:topconvers"

type TopConvers struct {
	Conversations []*Conversation `json:"conversations"`
}

func QryTopConvers(ctx context.Context, userId string, req *pbobjs.QryTopConversReq) (errs.IMErrorCode, *pbobjs.QryConversationsResp) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	resp := &pbobjs.QryConversationsResp{
		Conversations: []*pbobjs.Conversation{},
		IsFinished:    true,
	}
	converStorage := storages.NewConversationStorage()
	dbConvers, err := converStorage.QryTopConvers(appkey, userId, req.StartTime, 101, req.SortType, req.Order > 0)
	if err == nil {
		for _, dbConver := range dbConvers {
			conversation := dbConver2Conversations(ctx, dbConver)
			resp.Conversations = append(resp.Conversations, conversation)
		}
		if len(dbConvers) > 100 {
			resp.IsFinished = false
		}
	}
	return errs.IMErrorCode_SUCCESS, resp
}

func ClearTotalUnread(ctx context.Context, userId string) errs.IMErrorCode {
	appkey := bases.GetAppKeyFromCtx(ctx)
	storage := storages.NewConversationStorage()
	storage.ClearTotalUnreadCount(appkey, userId)
	commonservices.AsyncPrivateMsg(ctx, userId, userId, &pbobjs.UpMsg{
		MsgType:    "jg:cleartotalunread",
		MsgContent: []byte(fmt.Sprintf(`{"clear_time":%d}`, time.Now().UnixMilli())),
		Flags:      commonservices.SetCmdMsg(0),
	})
	return errs.IMErrorCode_SUCCESS
}

func ClearUnread(ctx context.Context, userId string, convers []*pbobjs.Conversation, noCmdMsg bool) errs.IMErrorCode {
	retErr := errs.IMErrorCode_SUCCESS
	appkey := bases.GetAppKeyFromCtx(ctx)
	if len(convers) > 0 {
		clearUnreadConvers := &ClearUnreadConvers{
			Conversations: []*ClearUnreadConver{},
		}
		converStorage := storages.NewConversationStorage()
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
				rowsAffected, err := converStorage.UpdateLatestReadMsgIndex(appkey, userId, conver.TargetId, conver.ChannelType, msgIndex, readMsgId, readMsgTime)
				if rowsAffected > 0 {
					affected = true
				}
				if err != nil {
					continue
				}
			} else {
				converStorage.ClearUnread(appkey, userId, conver.TargetId, conver.ChannelType)
			}
		}
		if affected && !noCmdMsg {
			//Notify other device to clear unread
			flag := commonservices.SetCmdMsg(0)
			bs, _ := json.Marshal(clearUnreadConvers)
			exts := bases.GetExtsFromCtx(ctx)
			if len(convers) == 1 {
				exts[commonservices.RpcExtKey_UniqTag] = fmt.Sprintf("%s_%d", convers[0].TargetId, convers[0].ChannelType)
			}
			commonservices.AsyncPrivateMsg(ctx, userId, userId, &pbobjs.UpMsg{
				MsgType:    commonservices.CmdMsgType_ClearUnread,
				MsgContent: bs,
				Flags:      flag,
			}, &bases.ExtsOption{Exts: exts})
		}
	}
	return retErr
}

type ClearUnreadConvers struct {
	Conversations []*ClearUnreadConver `json:"conversations"`
}

type ClearUnreadConver struct {
	TargetId           string `json:"target_id"`
	ChannelType        int32  `json:"channel_type"`
	LatestReadMsgIndex int64  `json:"latest_read_index"`
}

type MarkUnreadConvers struct {
	Conversations []*MarkUnreadConver `json:"conversations"`
}

type MarkUnreadConver struct {
	TargetId    string `json:"target_id"`
	ChannelType int32  `json:"channel_type"`
	UnreadTag   int    `json:"unread_tag"`
}

func MarkUnread(ctx context.Context, userId string, req *pbobjs.ConversationsReq) errs.IMErrorCode {
	appkey := bases.GetAppKeyFromCtx(ctx)
	if len(req.Conversations) > 0 {
		markUnreadConvers := &MarkUnreadConvers{
			Conversations: []*MarkUnreadConver{},
		}
		converStorage := storages.NewConversationStorage()
		affected := false
		for _, conver := range req.Conversations {
			markUnreadConvers.Conversations = append(markUnreadConvers.Conversations, &MarkUnreadConver{
				TargetId:    conver.TargetId,
				ChannelType: int32(conver.ChannelType),
				UnreadTag:   int(conver.UnreadTag),
			})
			rowsAffected, err := converStorage.UpdateUnreadTag(appkey, userId, conver.TargetId, conver.ChannelType)
			if rowsAffected > 0 {
				affected = true
			}
			if err != nil {
				continue
			}
		}
		if affected {
			//Notify other device to mark unread
			flag := commonservices.SetCmdMsg(0)
			bs, _ := json.Marshal(markUnreadConvers)
			commonservices.AsyncPrivateMsg(ctx, userId, userId, &pbobjs.UpMsg{
				MsgType:    commonservices.CmdMsgType_MarkUnread,
				MsgContent: bs,
				Flags:      flag,
			})
		}
	}
	return errs.IMErrorCode_SUCCESS
}

var undisturbMsgType string = "jg:undisturb"

type UndisturbConvers struct {
	Conversations []*UndisturbConver `json:"conversations"`
}
type UndisturbConver struct {
	TargetId      string `json:"target_id"`
	ChannelType   int32  `json:"channel_type"`
	UndisturbType int32  `json:"undisturb_type"`
}

func DoUndisturbConvers(ctx context.Context, req *pbobjs.UndisturbConversReq) errs.IMErrorCode {
	appkey := bases.GetAppKeyFromCtx(ctx)
	storage := storages.NewConversationStorage()
	userId := bases.GetTargetIdFromCtx(ctx)
	convers := &UndisturbConvers{
		Conversations: []*UndisturbConver{},
	}
	for _, item := range req.Items {
		storage.UpdateUndisturbType(appkey, userId, item.TargetId, item.ChannelType, item.UndisturbType)
		convers.Conversations = append(convers.Conversations, &UndisturbConver{
			TargetId:      item.TargetId,
			ChannelType:   int32(item.ChannelType),
			UndisturbType: item.UndisturbType,
		})
	}
	//Notify other device to update undisturb
	flag := commonservices.SetCmdMsg(0)
	bs, _ := json.Marshal(convers)
	commonservices.AsyncPrivateMsg(ctx, userId, userId, &pbobjs.UpMsg{
		MsgType:    undisturbMsgType,
		MsgContent: bs,
		Flags:      flag,
	})
	return errs.IMErrorCode_SUCCESS
}

func QryConver(ctx context.Context, userId string, req *pbobjs.QryConverReq) (errs.IMErrorCode, *pbobjs.Conversation) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	storage := storages.NewConversationStorage()
	conver, err := storage.FindOne(appkey, userId, req.TargetId, req.ChannelType)
	if err == nil && conver != nil {
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
			return errs.IMErrorCode_SUCCESS, dbConver2Conversations(ctx, conver)
		}
	}
	return errs.IMErrorCode_SUCCESS, &pbobjs.Conversation{
		TargetId:          req.TargetId,
		ChannelType:       req.ChannelType,
		UndisturbType:     0,
		LatestUnreadIndex: 0,
	}
}

func BatchQryConvers(ctx context.Context, req *pbobjs.QryConverReq) (errs.IMErrorCode, *pbobjs.QryConversationsResp) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	storage := storages.NewConversationStorage()
	reqConvers := []models.Conversation{}
	for _, uId := range req.UserIds {
		reqConvers = append(reqConvers, models.Conversation{
			UserId:      uId,
			TargetId:    req.TargetId,
			ChannelType: req.ChannelType,
		})
	}
	ret := &pbobjs.QryConversationsResp{
		Conversations: []*pbobjs.Conversation{},
	}
	convers, err := storage.BatchFind(appkey, reqConvers)
	if err == nil {
		for _, conver := range convers {
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

func GetGlobalConverId(senderId, targetId string, channelType pbobjs.ChannelType) string {
	converId := commonservices.GetConversationId(senderId, targetId, channelType)
	return tools.SHA1(converId)
}

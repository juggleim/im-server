package services

import (
	"context"
	"fmt"
	"im-server/commons/bases"
	"im-server/commons/caches"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/commonservices"
	"im-server/services/conversation/storages"
	"im-server/services/conversation/storages/models"
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

type GlobalConversationCacheItem struct {
	Appkey      string
	ConverId    string
	SenderId    string
	TargetId    string
	ChannelType pbobjs.ChannelType
	MsgTime     int64
}

func SaveGlobalConversationByCache(item *GlobalConversationCacheItem) {
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

func HandleMentionedMsg(appkey string, userId string, msg *pbobjs.DownMsg, userConvers *UserConversations) {
	if commonservices.IsMentionedMe(userId, msg) {
		isRead := 0
		if msg.IsSend {
			isRead = 1
		}
		//append to cache
		if userConvers == nil {
			userConvers = getUserConvers(appkey, userId)
		}
		userConvers.AppendMention(msg.TargetId, msg.ChannelType, &models.MentionMsg{
			SenderId:    msg.SenderId,
			MsgId:       msg.MsgId,
			MsgTime:     msg.MsgTime,
			MsgIndex:    msg.UnreadIndex,
			MentionType: msg.MentionInfo.MentionType,
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

var topConversMsgType string = "jg:topconvers"

type TopConvers struct {
	Conversations []*Conversation `json:"conversations"`
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

var undisturbMsgType string = "jg:undisturb"

type UndisturbConvers struct {
	Conversations []*UndisturbConver `json:"conversations"`
}
type UndisturbConver struct {
	TargetId      string `json:"target_id"`
	ChannelType   int32  `json:"channel_type"`
	UndisturbType int32  `json:"undisturb_type"`
}

func GetGlobalConverId(senderId, targetId string, channelType pbobjs.ChannelType) string {
	converId := commonservices.GetConversationId(senderId, targetId, channelType)
	return tools.SHA1(converId)
}

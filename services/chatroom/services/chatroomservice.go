package services

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/caches"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/chatroom/storages"
	"im-server/services/chatroom/storages/models"
	"im-server/services/commonservices"
	"math"
	"strings"
	"sync"
	"time"

	"github.com/bytedance/gopkg/collection/zset"
)

var (
	chatroomCache *caches.LruCache
	chatroomLocks *tools.SegmentatedLocks
)

type ChatroomStatus int

const (
	ChatroomStatus_Normal   = 0
	ChatroomStatus_Destroy  = 1
	ChatroomStatus_NotExist = 2
)

func init() {
	chatroomCache = caches.NewLruCacheWithAddReadTimeout(10000, nil, 30*time.Minute, 30*time.Minute)
	chatroomLocks = tools.NewSegmentatedLocks(256)
}

type ChatroomContainer struct {
	Appkey   string
	ChatId   string
	ChatName string
	Status   ChatroomStatus // 0:normal; 1: destroy; 2: not exist;
	IsMute   bool

	MsgTime     int64
	MsgSeq      int64
	Members     *zset.Float64Set
	Atts        map[string]*pbobjs.ChatAttItem
	AttTime     int64
	autoDelAtts map[string]map[string]bool

	BanUsers   sync.Map
	MuteUsers  sync.Map
	AllowUsers sync.Map
}

func (container *ChatroomContainer) GetMsgTimeSeq(currentTime int64, flag int32) (int64, int64) {
	lock := chatroomLocks.GetLocks(container.Appkey, container.ChatId)
	lock.Lock()
	defer lock.Unlock()
	if currentTime <= container.MsgTime {
		container.MsgTime = container.MsgTime + 1
	} else {
		container.MsgTime = currentTime
	}
	if commonservices.IsStoreMsg(flag) {
		container.MsgSeq = container.MsgSeq + 1
	}
	return container.MsgTime, container.MsgSeq
}

func (container *ChatroomContainer) GetMsgTime(currentTime int64) int64 {
	lock := chatroomLocks.GetLocks(container.Appkey, container.ChatId)
	lock.Lock()
	defer lock.Unlock()
	if currentTime <= container.MsgTime {
		container.MsgTime = container.MsgTime + 1
	} else {
		container.MsgTime = currentTime
	}
	return container.MsgTime
}

func (container *ChatroomContainer) Destroy() {
	lock := chatroomLocks.GetLocks(container.Appkey, container.ChatId)
	lock.Lock()
	defer lock.Unlock()
	container.Status = ChatroomStatus_Destroy
	container.MsgTime = 0
	container.MsgSeq = 0
	container.Members = zset.NewFloat64()
	container.Atts = make(map[string]*pbobjs.ChatAttItem)
	container.autoDelAtts = make(map[string]map[string]bool)
	container.AttTime = 0
	container.BanUsers = sync.Map{}
	container.MuteUsers = sync.Map{}
	container.AllowUsers = sync.Map{}
}

func (container *ChatroomContainer) AddMember(memberId string) bool {
	if !container.Members.Contains(memberId) {
		container.Members.Add(float64(time.Now().UnixMilli()), memberId)
		return true
	}
	return false
}

func (container *ChatroomContainer) DelMember(memberId string) bool {
	if container.Members.Contains(memberId) {
		container.Members.Remove(memberId)
		return true
	}
	return false
}

func (container *ChatroomContainer) CheckMemberExist(memberId string) bool {
	return container.Members.Contains(memberId)
}

func (container *ChatroomContainer) CheckMemberMute(memberId string) bool {
	if _, exist := container.MuteUsers.Load(memberId); exist {
		return true
	}
	return false
}

func (container *ChatroomContainer) CheckMemberAllow(memberId string) bool {
	if _, exist := container.AllowUsers.Load(memberId); exist {
		return true
	}
	return false
}

func (container *ChatroomContainer) CheckMemberBan(memberId string) bool {
	if _, exist := container.BanUsers.Load(memberId); exist {
		return true
	}
	return false
}

func (container *ChatroomContainer) AddAtt(userId, attKey, attValue string, isForce, isAutoDel bool) (errs.IMErrorCode, int64) {
	key := getChatroomKey(container.Appkey, container.ChatId)
	lock := chatroomLocks.GetLocks(key)
	lock.Lock()
	defer lock.Unlock()
	maxAttCount := getMaxAttCount(container.Appkey)
	if old, exist := container.Atts[attKey]; !exist {
		length := len(container.Atts)
		if length >= maxAttCount {
			return errs.IMErrorCode_CHATROOM_ATTFULL, 0
		}
	} else {
		if old.OptType != pbobjs.ChatAttOptType_ChatAttOpt_Del && !isForce && old.UserId != userId {
			return errs.IMErrorCode_CHATROOM_SEIZEFAILED, 0
		}
		if old.UserId != "" && old.UserId != userId {
			// check autoDelAtts
			if memberKeys, exist := container.autoDelAtts[old.UserId]; exist {
				delete(memberKeys, attKey)
				if len(memberKeys) <= 0 {
					delete(container.autoDelAtts, old.UserId)
				}
			}
		}
	}
	currentTime := time.Now().UnixMilli()
	if currentTime <= container.AttTime {
		container.AttTime = container.AttTime + 1
	} else {
		container.AttTime = currentTime
	}
	container.Atts[attKey] = &pbobjs.ChatAttItem{
		Key:     attKey,
		Value:   attValue,
		UserId:  userId,
		AttTime: container.AttTime,
		OptType: pbobjs.ChatAttOptType_ChatAttOpt_Add,
	}
	if isAutoDel {
		if index, exist := container.autoDelAtts[userId]; exist {
			index[attKey] = true
		} else {
			container.autoDelAtts[userId] = make(map[string]bool)
			container.autoDelAtts[userId][attKey] = true
		}
	}
	return errs.IMErrorCode_SUCCESS, container.AttTime
}

func (container *ChatroomContainer) DelAtt(userId, attKey string, isForce bool) (errs.IMErrorCode, int64) {
	key := getChatroomKey(container.Appkey, container.ChatId)
	lock := chatroomLocks.GetLocks(key)
	lock.Lock()
	defer lock.Unlock()
	if old, exist := container.Atts[attKey]; exist {
		if old.OptType == pbobjs.ChatAttOptType_ChatAttOpt_Del {
			return errs.IMErrorCode_CHATROOM_HASDELETED, 0
		}
		if !isForce && old.UserId != userId {
			return errs.IMErrorCode_CHATROOM_SEIZEFAILED, 0
		}
		currentTime := time.Now().UnixMilli()
		if currentTime <= container.AttTime {
			container.AttTime = container.AttTime + 1
		} else {
			container.AttTime = currentTime
		}
		old.AttTime = container.AttTime
		old.Value = ""
		old.UserId = ""
		old.OptType = pbobjs.ChatAttOptType_ChatAttOpt_Del
		//check autoDelAtts
		if memberKeys, exist := container.autoDelAtts[userId]; exist {
			delete(memberKeys, attKey)
			if len(memberKeys) <= 0 {
				delete(container.autoDelAtts, userId)
			}
		}
		return errs.IMErrorCode_SUCCESS, container.AttTime
	} else {
		return errs.IMErrorCode_CHATROOM_ATTNOTEXIST, 0
	}
}

func (container *ChatroomContainer) QryAutoDelKeysByMember(memberId string) []string {
	key := getChatroomKey(container.Appkey, container.ChatId)
	lock := chatroomLocks.GetLocks(key)
	lock.RLock()
	defer lock.RUnlock()
	delKeys := []string{}
	if memberKeys, exist := container.autoDelAtts[memberId]; exist {
		for key := range memberKeys {
			delKeys = append(delKeys, key)
		}
	}
	return delKeys
}

func (container *ChatroomContainer) GetAtts() []*pbobjs.ChatAttItem {
	key := getChatroomKey(container.Appkey, container.ChatId)
	lock := chatroomLocks.GetLocks(key)
	lock.RLock()
	defer lock.RUnlock()
	atts := []*pbobjs.ChatAttItem{}
	for _, v := range container.Atts {
		atts = append(atts, &pbobjs.ChatAttItem{
			Key:     v.Key,
			Value:   v.Value,
			AttTime: v.AttTime,
			UserId:  v.UserId,
			OptType: v.OptType,
		})
	}
	return atts
}

func getChatroomContainer(appkey, chatId string) (*ChatroomContainer, bool) {
	key := getChatroomKey(appkey, chatId)
	if cacheContainer, exist := chatroomCache.Get(key); exist {
		container := cacheContainer.(*ChatroomContainer)
		if container.Status == ChatroomStatus_Normal {
			return container, true
		}
		return container, false
	} else {
		lock := chatroomLocks.GetLocks(key)
		lock.Lock()
		defer lock.Unlock()
		if cacheContainer, exist := chatroomCache.Get(key); exist {
			container := cacheContainer.(*ChatroomContainer)
			if container.Status == ChatroomStatus_Normal {
				return container, true
			}
			return container, false
		} else {
			storage := storages.NewChatroomStorage()
			chatroom, err := storage.FindById(appkey, chatId)
			container := &ChatroomContainer{
				Appkey:  appkey,
				ChatId:  chatId,
				IsMute:  chatroom.IsMute > 0,
				Members: zset.NewFloat64(),
				Atts:    make(map[string]*pbobjs.ChatAttItem),
			}
			if err == nil && chatroom != nil {
				container.Status = ChatroomStatus_Normal
				//init members of chatroom
				memberStorages := storages.NewChatroomMemberStorage()
				var startId int64 = 0
				for {
					members, err := memberStorages.QryMembers(appkey, chatId, true, int64(startId), 1000)
					if err != nil || members == nil {
						break
					}
					for _, member := range members {
						startId = member.ID
						container.Members.Add(float64(member.CreatedTime.UnixMilli()), member.MemberId)
					}
					if len(members) < 1000 {
						break
					}
				}
				// init chatroom attributes
				extStorage := storages.NewChatroomExtStorage()
				exts, err := extStorage.QryExts(appkey, chatId)
				if err == nil {
					for _, ext := range exts {
						oType := pbobjs.ChatAttOptType_ChatAttOpt_Add
						if ext.IsDelete == 1 {
							oType = pbobjs.ChatAttOptType_ChatAttOpt_Del
						}
						container.Atts[ext.ItemKey] = &pbobjs.ChatAttItem{
							Key:     ext.ItemKey,
							Value:   ext.ItemValue,
							UserId:  ext.MemberId,
							AttTime: ext.ItemTime,
							OptType: oType,
						}
						if ext.ItemTime > container.AttTime {
							container.AttTime = ext.ItemTime
						}
					}
				}
				//init banusers
				foreachBanUsersFromDb(appkey, chatId, pbobjs.ChrmBanType_Ban, func(banUser *models.ChatroomBanUser) {
					container.BanUsers.Store(banUser.MemberId, banUser.CreatedTime)
				})
				//init muteusers
				foreachBanUsersFromDb(appkey, chatId, pbobjs.ChrmBanType_Mute, func(banUser *models.ChatroomBanUser) {
					container.MuteUsers.Store(banUser.MemberId, banUser.CreatedTime)
				})
				//init allowusers
				foreachBanUsersFromDb(appkey, chatId, pbobjs.ChrmBanType_Allow, func(banUser *models.ChatroomBanUser) {
					container.AllowUsers.Store(banUser.MemberId, banUser.CreatedTime)
				})
			} else {
				container.Status = ChatroomStatus_NotExist
			}
			chatroomCache.Add(key, container)
			return container, container.Status == ChatroomStatus_Normal
		}
	}
}

func initChatroomContainer(appkey, chatId, chatName string, isMute bool) {
	key := getChatroomKey(appkey, chatId)
	lock := chatroomLocks.GetLocks(key)
	lock.Lock()
	defer lock.Unlock()
	if !chatroomCache.Contains(key) {
		chatroomCache.Add(key, &ChatroomContainer{
			Appkey:     appkey,
			ChatId:     chatId,
			ChatName:   chatName,
			Status:     ChatroomStatus_Normal,
			IsMute:     isMute,
			Members:    zset.NewFloat64(),
			Atts:       make(map[string]*pbobjs.ChatAttItem),
			BanUsers:   sync.Map{},
			MuteUsers:  sync.Map{},
			AllowUsers: sync.Map{},
		})
	}
}

func getChatroomKey(appkey, chatId string) string {
	return strings.Join([]string{appkey, chatId}, "_")
}

func GetTopChatroomMembers(appkey, chatId string, count, order int32) []string {
	members := []string{}
	key := getChatroomKey(appkey, chatId)
	if cacheContainer, exist := chatroomCache.Get(key); exist {
		container := cacheContainer.(*ChatroomContainer)
		if container.Members != nil {
			nodeList := container.Members.Range(0, math.MaxInt)
			for _, node := range nodeList {
				members = append(members, node.Value)
			}
		}
	}
	return members
}

func CreateChatroom(ctx context.Context, req *pbobjs.ChatroomInfo) errs.IMErrorCode {
	appkey := bases.GetAppKeyFromCtx(ctx)
	//add for cache
	container, exist := getChatroomContainer(appkey, req.ChatId)
	if exist {
		if req.ChatName != "" {
			container.ChatName = req.ChatName
		}
		return errs.IMErrorCode_SUCCESS
	}
	initChatroomContainer(appkey, req.ChatId, req.ChatName, req.IsMute)
	//add to db
	storage := storages.NewChatroomStorage()
	storage.Create(models.Chatroom{
		ChatId:      req.ChatId,
		ChatName:    req.ChatName,
		CreatedTime: time.Now(),
		UpdatedTime: time.Now(),
		AppKey:      appkey,
	})
	//notify chatmsg nodes
	bases.Broadcast(ctx, "c_chrm_dispatch", &pbobjs.ChrmDispatchReq{
		ChatId:       req.ChatId,
		DispatchType: pbobjs.ChrmDispatchType_CreateChatroom,
	})
	return errs.IMErrorCode_SUCCESS
}

func DestroyChatroom(ctx context.Context, req *pbobjs.ChatroomInfo) errs.IMErrorCode {
	appkey := bases.GetAppKeyFromCtx(ctx)
	//destroy from cache
	key := getChatroomKey(appkey, req.ChatId)
	chatroomCache.Remove(key)
	container, exist := getChatroomContainer(appkey, req.ChatId)
	if exist {
		container.Destroy()
		//delete from db
		storage := storages.NewChatroomStorage()
		storage.Delete(appkey, req.ChatId)
		//delete members of chatroom
		memberStorage := storages.NewChatroomMemberStorage()
		memberStorage.ClearMembers(appkey, req.ChatId)
		//delete atts of chatroom
		extStorage := storages.NewChatroomExtStorage()
		extStorage.ClearExts(appkey, req.ChatId)
		//delete banusers of chatroom
		banUserStorage := storages.NewChatroomBanUserStorage()
		banUserStorage.ClearBanUsers(appkey, req.ChatId)

		//notify chatmsg nodes
		bases.Broadcast(ctx, "c_chrm_dispatch", &pbobjs.ChrmDispatchReq{
			ChatId:       req.ChatId,
			DispatchType: pbobjs.ChrmDispatchType_DestroyChatroom,
		})
	}

	return errs.IMErrorCode_SUCCESS
}

func JoinChatroom(ctx context.Context, chat *pbobjs.ChatroomInfo) errs.IMErrorCode {
	appkey := bases.GetAppKeyFromCtx(ctx)
	userId := bases.GetRequesterIdFromCtx(ctx)
	chatId := bases.GetTargetIdFromCtx(ctx)
	container, exist := getChatroomContainer(appkey, chatId)
	if !exist {
		return errs.IMErrorCode_CHATROOM_NOTEXIST
	}
	//check ban
	if container.CheckMemberBan(userId) {
		return errs.IMErrorCode_CHATROOM_BAN
	}
	//add to cache
	succ := container.AddMember(userId)
	if succ {
		//add to db
		storage := storages.NewChatroomMemberStorage()
		storage.Create(models.ChatroomMember{
			ChatId:      chatId,
			MemberId:    userId,
			CreatedTime: time.Now(),
			AppKey:      appkey,
		})
	}
	//notify owner node
	bases.AsyncRpcCall(ctx, "c_members_dispatch", userId, &pbobjs.ChatMembersDispatchReq{
		ChatId:       chatId,
		MemberIds:    []string{userId},
		DispatchType: pbobjs.ChatMembersDispatchType_JoinChatroom,
	})

	return errs.IMErrorCode_SUCCESS
}

func QuitChatroom(ctx context.Context, userId string, chat *pbobjs.ChatroomInfo) errs.IMErrorCode {
	appkey := bases.GetAppKeyFromCtx(ctx)
	chatId := chat.ChatId

	container, exist := getChatroomContainer(appkey, chatId)
	if !exist {
		return errs.IMErrorCode_CHATROOM_NOTEXIST
	}
	//remove from cache
	succ := container.DelMember(userId)
	if succ {
		//remove from db
		storage := storages.NewChatroomMemberStorage()
		storage.DeleteMember(appkey, chatId, userId)
	}
	//notify owner node
	bases.AsyncRpcCall(ctx, "c_members_dispatch", userId, &pbobjs.ChatMembersDispatchReq{
		ChatId:       chatId,
		MemberIds:    []string{userId},
		DispatchType: pbobjs.ChatMembersDispatchType_QuitChatroom,
	})
	//clean auto-del keys
	delKeys := container.QryAutoDelKeysByMember(userId)
	if len(delKeys) > 0 {
		atts := &pbobjs.ChatAtts{
			ChatId: chatId,
			Atts:   []*pbobjs.ChatAttItem{},
		}
		for _, key := range delKeys {
			code, attTime := container.DelAtt(userId, key, true)
			if code == errs.IMErrorCode_SUCCESS {
				atts.Atts = append(atts.Atts, &pbobjs.ChatAttItem{
					Key:     key,
					AttTime: attTime,
					UserId:  userId,
					OptType: pbobjs.ChatAttOptType_ChatAttOpt_Del,
				})
			}
		}
		if len(atts.Atts) > 0 {
			bases.Broadcast(ctx, "c_atts_dispatch", atts)
		}
	}

	return errs.IMErrorCode_SUCCESS
}

func SyncQuitChatroom(ctx context.Context, member *pbobjs.ChatroomMember) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	chatId := bases.GetTargetIdFromCtx(ctx)
	userId := member.MemberId
	container, exist := getChatroomContainer(appkey, chatId)
	if exist {
		succ := container.DelMember(userId)
		if succ {
			//remove from db
			storage := storages.NewChatroomMemberStorage()
			storage.DeleteMember(appkey, chatId, userId)
		}
		//clean auto-del keys
		delKeys := container.QryAutoDelKeysByMember(userId)
		if len(delKeys) > 0 {
			atts := &pbobjs.ChatAtts{
				ChatId: chatId,
				Atts:   []*pbobjs.ChatAttItem{},
			}
			for _, key := range delKeys {
				code, attTime := container.DelAtt(userId, key, true)
				if code == errs.IMErrorCode_SUCCESS {
					atts.Atts = append(atts.Atts, &pbobjs.ChatAttItem{
						Key:     key,
						AttTime: attTime,
						UserId:  userId,
						OptType: pbobjs.ChatAttOptType_ChatAttOpt_Del,
					})
				}
			}
			if len(atts.Atts) > 0 {
				bases.Broadcast(ctx, "c_atts_dispatch", atts)
			}
		}
	}
}

func GetPartialMembers(ctx context.Context, chatId, nodeName, method string) (*pbobjs.ChatroomInfo, errs.IMErrorCode) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	container, exist := getChatroomContainer(appkey, chatId)
	if !exist {
		return nil, errs.IMErrorCode_CHATROOM_NOTEXIST
	}
	ret := &pbobjs.ChatroomInfo{
		ChatId:   chatId,
		ChatName: container.ChatName,
		Members:  []*pbobjs.ChatroomMember{},
	}
	nodes := container.Members.RangeByScore(0, float64(time.Now().UnixMilli()))
	for _, node := range nodes {
		n := bases.GetCluster().GetTargetNode(method, node.Value)
		if n.Name == nodeName {
			ret.Members = append(ret.Members, &pbobjs.ChatroomMember{
				MemberId: node.Value,
			})
		}
	}
	return ret, errs.IMErrorCode_SUCCESS
}

func GetPartialInfo(ctx context.Context, chatId, nodeName, method string) (*pbobjs.ChatroomInfo, errs.IMErrorCode) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	container, exist := getChatroomContainer(appkey, chatId)
	if !exist {
		return nil, errs.IMErrorCode_CHATROOM_NOTEXIST
	}
	ret := &pbobjs.ChatroomInfo{
		ChatId:   chatId,
		ChatName: container.ChatName,
		Members:  []*pbobjs.ChatroomMember{},
		Atts:     []*pbobjs.ChatAttItem{},
	}
	nodes := container.Members.RangeByScore(0, float64(time.Now().UnixMilli()))
	for _, node := range nodes {
		n := bases.GetCluster().GetTargetNode(method, node.Value)
		if n.Name == nodeName {
			ret.Members = append(ret.Members, &pbobjs.ChatroomMember{
				MemberId: node.Value,
			})
		}
	}
	ret.Atts = container.GetAtts()
	return ret, errs.IMErrorCode_SUCCESS
}

func QryChatroomInfo(ctx context.Context, chatId string, count, order int) (errs.IMErrorCode, *pbobjs.ChatroomInfo) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	container, exist := getChatroomContainer(appkey, chatId)
	if !exist {
		return errs.IMErrorCode_CHATROOM_NOTEXIST, nil
	}
	ret := &pbobjs.ChatroomInfo{
		ChatId:      chatId,
		ChatName:    container.ChatName,
		IsMute:      container.IsMute,
		MemberCount: int32(container.Members.Len()),
		Members:     []*pbobjs.ChatroomMember{},
		Atts:        container.GetAtts(),
	}
	//members
	var nodes []zset.Float64Node
	if order == 0 {
		nodes = container.Members.RevRange(0, count)
	} else {
		nodes = container.Members.Range(0, count)
	}
	if len(nodes) > 0 {
		for _, node := range nodes {
			ret.Members = append(ret.Members, &pbobjs.ChatroomMember{
				MemberId:  node.Value,
				AddedTime: int64(node.Score),
			})
		}
	}
	return errs.IMErrorCode_SUCCESS, ret
}

func SetChrmMute(ctx context.Context, req *pbobjs.ChatroomInfo) errs.IMErrorCode {
	appkey := bases.GetAppKeyFromCtx(ctx)
	container, exist := getChatroomContainer(appkey, req.ChatId)
	if !exist {
		return errs.IMErrorCode_CHATROOM_NOTEXIST
	}
	if container.IsMute != req.IsMute {
		container.IsMute = req.IsMute
		//update db
		storage := storages.NewChatroomStorage()
		var isMute int = 0
		if req.IsMute {
			isMute = 1
		}
		storage.UpdateMute(appkey, req.ChatId, isMute)
	}

	return errs.IMErrorCode_SUCCESS
}

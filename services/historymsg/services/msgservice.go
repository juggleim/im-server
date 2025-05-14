package services

import (
	"fmt"
	"im-server/commons/caches"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/historymsg/storages"
	"time"
)

var (
	messageCache *caches.LruCache
	messageLocks *tools.SegmentatedLocks
)

func init() {
	messageCache = caches.NewLruCacheWithReadTimeout("hismessage_cache", 100000, nil, 10*time.Minute)
	messageLocks = tools.NewSegmentatedLocks(512)
}

type MsgInfo struct {
	Appkey      string
	ChannelType pbobjs.ChannelType
	ConverId    string
	MsgId       string
	SenderId    string
	//readinfo
	ReadMembers map[string]int64
	MemberCount int
	//ext
	MsgExtMap map[string]*pbobjs.MsgExtItem
	//exset
	MsgExsetMap     map[string][]*pbobjs.MsgExtItem
	MsgExsetUniqMap map[string]bool
}

func (info *MsgInfo) AddReadMembers(members map[string]int64) (bool, int) {
	isAdded := false
	for memberId, addedTime := range members {
		if _, exist := info.ReadMembers[memberId]; !exist {
			info.ReadMembers[memberId] = addedTime
			isAdded = true
		}
	}
	return isAdded, len(info.ReadMembers)
}

func (info *MsgInfo) AddReadMember(memberId string, addedTime int64) (bool, int) {
	key := getMsgInfoCacheKey(info.Appkey, info.ConverId, info.MsgId, info.ChannelType)
	lock := messageLocks.GetLocks(key)
	lock.Lock()
	defer lock.Unlock()
	if _, exist := info.ReadMembers[memberId]; exist {
		return false, len(info.ReadMembers)
	} else {
		info.ReadMembers[memberId] = addedTime
		return true, len(info.ReadMembers)
	}
}

func (info *MsgInfo) GetReadMemberCount() int {
	key := getMsgInfoCacheKey(info.Appkey, info.ConverId, info.MsgId, info.ChannelType)
	lock := messageLocks.GetLocks(key)
	lock.RLock()
	defer lock.RUnlock()
	return len(info.ReadMembers)
}

func (info *MsgInfo) SetMsgExt(ext *pbobjs.MsgExtItem) errs.IMErrorCode {
	key := getMsgInfoCacheKey(info.Appkey, info.ConverId, info.MsgId, info.ChannelType)
	lock := messageLocks.GetLocks(key)
	lock.Lock()
	defer lock.Unlock()
	if _, exist := info.MsgExtMap[ext.Key]; exist {
		return errs.IMErrorCode_MSG_MSGEXTDUPLICATE
	}
	if len(info.MsgExtMap) >= 100 {
		return errs.IMErrorCode_MSG_MSGEXTOVERLIMIT
	}
	info.MsgExtMap[ext.Key] = &pbobjs.MsgExtItem{
		Key:       ext.Key,
		Value:     ext.Value,
		Timestamp: ext.Timestamp,
		UserInfo:  ext.UserInfo,
	}
	return errs.IMErrorCode_SUCCESS
}

func (info *MsgInfo) DelMsgExt(extKey string) bool {
	key := getMsgInfoCacheKey(info.Appkey, info.ConverId, info.MsgId, info.ChannelType)
	lock := messageLocks.GetLocks(key)
	lock.Lock()
	defer lock.Unlock()
	if _, exist := info.MsgExtMap[extKey]; exist {
		delete(info.MsgExtMap, extKey)
		return true
	}
	return false
}

func (info *MsgInfo) ForeachMsgExt(f func(key string, ext *pbobjs.MsgExtItem)) {
	key := getMsgInfoCacheKey(info.Appkey, info.ConverId, info.MsgId, info.ChannelType)
	lock := messageLocks.GetLocks(key)
	lock.RLock()
	defer lock.RUnlock()
	for k, val := range info.MsgExtMap {
		f(k, val)
	}
}

func (info *MsgInfo) AddMsgExset(ext *pbobjs.MsgExtItem) errs.IMErrorCode {
	key := getMsgInfoCacheKey(info.Appkey, info.ConverId, info.MsgId, info.ChannelType)
	lock := messageLocks.GetLocks(key)
	lock.Lock()
	defer lock.Unlock()
	uniqKey := fmt.Sprintf("%s_%s", ext.Key, ext.Value)
	if _, exist := info.MsgExsetUniqMap[uniqKey]; exist {
		return errs.IMErrorCode_MSG_MSGEXTDUPLICATE
	}
	if list, exist := info.MsgExsetMap[ext.Key]; exist {
		if len(info.MsgExsetMap[ext.Key]) >= 100 {
			return errs.IMErrorCode_MSG_MSGEXTOVERLIMIT
		}
		info.MsgExsetMap[ext.Key] = append(list, ext)
	} else {
		if len(info.MsgExsetMap) >= 20 {
			return errs.IMErrorCode_MSG_MSGEXTOVERLIMIT
		}
		info.MsgExsetMap[ext.Key] = []*pbobjs.MsgExtItem{ext}
	}
	info.MsgExsetUniqMap[uniqKey] = true
	return errs.IMErrorCode_SUCCESS
}

func (info *MsgInfo) DelMsgExset(extKey, extVal string) bool {
	key := getMsgInfoCacheKey(info.Appkey, info.ConverId, info.MsgId, info.ChannelType)
	lock := messageLocks.GetLocks(key)
	lock.Lock()
	defer lock.Unlock()
	uniqKey := fmt.Sprintf("%s_%s", extKey, extVal)
	if _, exist := info.MsgExsetUniqMap[uniqKey]; !exist {
		return false
	}
	delete(info.MsgExsetUniqMap, uniqKey)
	list, exist := info.MsgExsetMap[extKey]
	if exist && len(list) > 0 {
		newList := []*pbobjs.MsgExtItem{}
		for _, item := range list {
			if item.Value != extVal {
				newList = append(newList, item)
			}
		}
		info.MsgExsetMap[extKey] = newList
		return true
	} else {
		return false
	}
}

func (info *MsgInfo) ForeachMsgExset(f func(key string, exts []*pbobjs.MsgExtItem)) {
	key := getMsgInfoCacheKey(info.Appkey, info.ConverId, info.MsgId, info.ChannelType)
	lock := messageLocks.GetLocks(key)
	lock.RLock()
	defer lock.RUnlock()
	for k, val := range info.MsgExsetMap {
		f(k, val)
	}
}

func getMsgInfoCacheKey(appkey, converId, msgId string, channelType pbobjs.ChannelType) string {
	return fmt.Sprintf("%s_%s_%s_%d", appkey, msgId, converId, channelType)
}

func GetMsgInfo(appkey, converId, msgId string, channelType pbobjs.ChannelType) *MsgInfo {
	key := getMsgInfoCacheKey(appkey, converId, msgId, channelType)
	if info, exist := messageCache.Get(key); exist {
		return info.(*MsgInfo)
	} else {
		lock := messageLocks.GetLocks(key)
		lock.Lock()
		defer lock.Unlock()
		if info, exist := messageCache.Get(key); exist {
			return info.(*MsgInfo)
		} else {
			msgInfo := &MsgInfo{
				Appkey:      appkey,
				ChannelType: channelType,
				ConverId:    converId,
				MsgId:       msgId,

				ReadMembers:     make(map[string]int64),
				MsgExtMap:       make(map[string]*pbobjs.MsgExtItem),
				MsgExsetMap:     make(map[string][]*pbobjs.MsgExtItem),
				MsgExsetUniqMap: make(map[string]bool),
			}
			msgExt := []byte{}
			msgExset := []byte{}
			if channelType == pbobjs.ChannelType_Group {
				hisStorage := storages.NewGroupHisMsgStorage()
				grpMsg, err := hisStorage.FindById(appkey, converId, msgId)
				if err == nil && grpMsg != nil {
					msgExt = grpMsg.MsgExt
					msgExset = grpMsg.MsgExset
					msgInfo.MemberCount = grpMsg.MemberCount
					msgInfo.SenderId = grpMsg.SenderId

					//readinfo
					readInfoStorage := storages.NewReadInfoStorage()
					infos, err := readInfoStorage.QryReadInfosByMsgId(appkey, converId, channelType, msgId, 0, GrpMsgReadInfoLimit)
					if err == nil {
						members := make(map[string]int64)
						for _, info := range infos {
							members[info.MemberId] = info.CreatedTime
						}
						msgInfo.AddReadMembers(members)
					}
				}
			} else if channelType == pbobjs.ChannelType_Private {
				hisStorage := storages.NewPrivateHisMsgStorage()
				priMsg, err := hisStorage.FindById(appkey, converId, msgId)
				if err == nil && priMsg != nil {
					msgExt = priMsg.MsgExt
					msgExset = priMsg.MsgExset
					msgInfo.SenderId = priMsg.SenderId
				}
			}
			if len(msgExt) > 0 {
				extItems := MsgExtBs2Pb(msgExt)
				if extItems != nil && len(extItems.Exts) > 0 {
					for _, extItem := range extItems.Exts {
						if extItem.Key != "" {
							val := &pbobjs.MsgExtItem{
								Key:       extItem.Key,
								Value:     extItem.Value,
								Timestamp: extItem.Timestamp,
								UserInfo:  extItem.UserInfo,
							}
							msgInfo.MsgExtMap[extItem.Key] = val
						}
					}
				}
			}
			if len(msgExset) > 0 {
				extItems := MsgExtBs2Pb(msgExset)
				if extItems != nil && len(extItems.Exts) > 0 {
					for _, extItem := range extItems.Exts {
						if extItem.Key != "" {
							val := &pbobjs.MsgExtItem{
								Key:       extItem.Key,
								Value:     extItem.Value,
								Timestamp: extItem.Timestamp,
								UserInfo:  extItem.UserInfo,
							}
							uniqKey := fmt.Sprintf("%s_%s", extItem.Key, extItem.Value)
							if _, exist := msgInfo.MsgExsetUniqMap[uniqKey]; !exist {
								msgInfo.MsgExsetUniqMap[uniqKey] = true
								if list, exist := msgInfo.MsgExsetMap[extItem.Key]; exist {
									msgInfo.MsgExsetMap[extItem.Key] = append(list, val)
								} else {
									msgInfo.MsgExsetMap[extItem.Key] = []*pbobjs.MsgExtItem{val}
								}
							}
						}
					}
				}
			}
			messageCache.Add(key, msgInfo)
			return msgInfo
		}
	}
}

func MsgExtBs2Pb(bs []byte) *pbobjs.MsgExtItems {
	if len(bs) > 0 {
		items := &pbobjs.MsgExtItems{}
		err := tools.PbUnMarshal(bs, items)
		if err == nil {
			return items
		}
	}
	return nil
}

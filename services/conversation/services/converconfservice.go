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
	"im-server/services/commonservices/convercache"
	converStorage "im-server/services/conversation/storages"
	"im-server/services/conversation/storages/models"
	"time"

	"google.golang.org/protobuf/proto"
)

type ConverConfItem struct {
	AppKey      string
	ConverId    string
	ChannelType pbobjs.ChannelType
	SubChannel  string

	//rtc room
	RtcRoomId                 string
	ActivedRtcRoom            *pbobjs.ActivedRtcRoom
	LatestActivedRtcRoomCheck int64

	GlobalConverTags map[string]bool
}

var converConfCache *caches.LruCache
var converConfLocks *tools.SegmentatedLocks

func init() {
	converConfCache = caches.NewLruCacheWithAddReadTimeout("converconf_cache", 100000, nil, 5*time.Minute, 5*time.Minute)
	converConfLocks = tools.NewSegmentatedLocks(256)
}

func getConverConfCacheKey(appkey, converId string, channelType pbobjs.ChannelType, subChannel string) string {
	return fmt.Sprintf("%s_%s_%d_%s", appkey, converId, channelType, subChannel)
}

func GetConverConf(appkey, converId string, channelType pbobjs.ChannelType, subChannel string) *ConverConfItem {
	key := getConverConfCacheKey(appkey, converId, channelType, subChannel)
	if val, exist := converConfCache.Get(key); exist {
		return val.(*ConverConfItem)
	} else {
		l := converConfLocks.GetLocks(key)
		l.Lock()
		defer l.Unlock()
		if val, exist := converConfCache.Get(key); exist {
			return val.(*ConverConfItem)
		} else {
			item := &ConverConfItem{
				AppKey:      appkey,
				ConverId:    converId,
				ChannelType: channelType,
				SubChannel:  subChannel,
			}
			globalConverTagsConf, err := converStorage.NewConverConfStorage().Find(
				appkey,
				converId,
				channelType,
				subChannel,
				string(commonservices.AttItemKey_GlobalConverConf_GlobalConverTags),
			)
			if err == nil && globalConverTagsConf != nil &&
				globalConverTagsConf.ItemType == int(commonservices.AttItemType_Setting) &&
				globalConverTagsConf.ItemValue != "" {
				globalConverTags := map[string]bool{}
				if err := json.Unmarshal([]byte(globalConverTagsConf.ItemValue), &globalConverTags); err == nil {
					item.GlobalConverTags = globalConverTags
				}
			}
			converConfCache.Add(key, item)
			return item
		}
	}
}

func (conf *ConverConfItem) SetRtcRoomId(rtcRoomId string) bool {
	l := converConfLocks.GetLocks(getConverConfCacheKey(conf.AppKey, conf.ConverId, conf.ChannelType, conf.SubChannel))
	l.Lock()
	defer l.Unlock()
	if conf.RtcRoomId == "" || conf.RtcRoomId == rtcRoomId {
		conf.RtcRoomId = rtcRoomId
		return true
	}
	return false
}

func (conf *ConverConfItem) ClearRtcRoomId() {
	l := converConfLocks.GetLocks(getConverConfCacheKey(conf.AppKey, conf.ConverId, conf.ChannelType, conf.SubChannel))
	l.Lock()
	defer l.Unlock()
	conf.RtcRoomId = ""
	conf.ActivedRtcRoom = nil
	conf.LatestActivedRtcRoomCheck = 0
}

func (conf *ConverConfItem) SetGlobalConverTagsIfChanged(globalConverTags map[string]bool, persist func() error) error {
	l := converConfLocks.GetLocks(getConverConfCacheKey(conf.AppKey, conf.ConverId, conf.ChannelType, conf.SubChannel))
	l.Lock()
	defer l.Unlock()

	if globalConverTagsEqual(conf.GlobalConverTags, globalConverTags) {
		return nil
	}
	if err := persist(); err != nil {
		return err
	}
	conf.GlobalConverTags = make(map[string]bool, len(globalConverTags))
	for tag, enabled := range globalConverTags {
		conf.GlobalConverTags[tag] = enabled
	}
	return nil
}

func (conf *ConverConfItem) GetGlobalConverTags() map[string]bool {
	l := converConfLocks.GetLocks(getConverConfCacheKey(conf.AppKey, conf.ConverId, conf.ChannelType, conf.SubChannel))
	l.RLock()
	defer l.RUnlock()

	ret := make(map[string]bool, len(conf.GlobalConverTags))
	for tag, enabled := range conf.GlobalConverTags {
		ret[tag] = enabled
	}
	return ret
}

func globalConverTagsEqual(left, right map[string]bool) bool {
	if len(left) != len(right) {
		return false
	}
	for tag, enabled := range left {
		if rightEnabled, exists := right[tag]; !exists || rightEnabled != enabled {
			return false
		}
	}
	return true
}

func SetGlobalConverConf(ctx context.Context, req *pbobjs.SetConverConfReq) errs.IMErrorCode {
	if req == nil || req.ConverId == "" || req.ChannelType == pbobjs.ChannelType_Unknown || req.ItemKey == "" {
		return errs.IMErrorCode_MSG_PARAM_ILLEGAL
	}
	if req.ItemKey != string(commonservices.AttItemKey_GlobalConverConf_GlobalConverTags) {
		return errs.IMErrorCode_SUCCESS
	}
	if req.ItemType != int32(commonservices.AttItemType_Setting) {
		return errs.IMErrorCode_MSG_PARAM_ILLEGAL
	}
	globalConverTags := map[string]bool{}
	if err := json.Unmarshal([]byte(req.ItemValue), &globalConverTags); err != nil || globalConverTags == nil {
		return errs.IMErrorCode_MSG_PARAM_ILLEGAL
	}

	appkey := bases.GetAppKeyFromCtx(ctx)
	storage := converStorage.NewConverConfStorage()
	confModel := models.ConverConf{
		AppKey:     appkey,
		ConverId:   req.ConverId,
		ConverType: req.ChannelType,
		SubChannel: req.SubChannel,
		ItemKey:    req.ItemKey,
		ItemValue:  req.ItemValue,
		ItemType:   int(req.ItemType),
	}
	converConf := GetConverConf(appkey, req.ConverId, req.ChannelType, req.SubChannel)
	if converConf == nil {
		return errs.IMErrorCode_DEFAULT
	}
	if err := converConf.SetGlobalConverTagsIfChanged(globalConverTags, func() error {
		return storage.Upsert(confModel)
	}); err != nil {
		return errs.IMErrorCode_DEFAULT
	}
	convercache.RemoveMsgConverCache(appkey, req.ConverId, req.SubChannel, req.ChannelType)
	return errs.IMErrorCode_SUCCESS
}

func (conf *ConverConfItem) GetActivedRtcRoom(ctx context.Context) *pbobjs.ActivedRtcRoom {
	if conf.RtcRoomId == "" {
		conf.ActivedRtcRoom = nil
		conf.LatestActivedRtcRoomCheck = 0
	}
	if (conf.ActivedRtcRoom == nil && conf.RtcRoomId != "") || time.Now().UnixMilli()-conf.LatestActivedRtcRoomCheck > 10*1000 {
		l := converConfLocks.GetLocks(getConverConfCacheKey(conf.AppKey, conf.ConverId, conf.ChannelType, conf.SubChannel))
		l.Lock()
		defer l.Unlock()
		if conf.RtcRoomId != "" {
			code, resp, err := bases.SyncRpcCall(ctx, "rtc_qry", conf.RtcRoomId, &pbobjs.Nil{}, func() proto.Message {
				return &pbobjs.RtcRoom{}
			})
			if err == nil {
				if code == errs.IMErrorCode_SUCCESS {
					if resp != nil {
						rtcRoom := resp.(*pbobjs.RtcRoom)
						conf.ActivedRtcRoom = &pbobjs.ActivedRtcRoom{
							RoomType:     rtcRoom.RoomType,
							RoomId:       rtcRoom.RoomId,
							Owner:        rtcRoom.Owner,
							RtcChannel:   rtcRoom.RtcChannel,
							RtcMediaType: rtcRoom.RtcMediaType,
							Ext:          rtcRoom.Ext,
							Members:      rtcRoom.Members,
						}
						conf.LatestActivedRtcRoomCheck = time.Now().UnixMilli()
					}
				} else {
					conf.RtcRoomId = ""
					conf.LatestActivedRtcRoomCheck = 0
				}
			}
		} else {
			conf.ActivedRtcRoom = nil
			conf.LatestActivedRtcRoomCheck = 0
		}
	}
	return conf.ActivedRtcRoom
}

func RtcConverBind(ctx context.Context, req *pbobjs.RtcConverBindReq) errs.IMErrorCode {
	appkey := bases.GetAppKeyFromCtx(ctx)
	converConf := GetConverConf(appkey, req.ConverId, req.ChannelType, req.SubChannel)
	if converConf != nil {
		if req.Action == pbobjs.RtcConverBindAction_RtcBind {
			succ := converConf.SetRtcRoomId(req.RtcRoomId)
			if !succ {
				return errs.IMErrorCode_RTCROOM_BINDCONVERFAIL
			}
		} else if req.Action == pbobjs.RtcConverBindAction_RtcUnBind {
			converConf.ClearRtcRoomId()
		}
	}
	return errs.IMErrorCode_SUCCESS
}

func QryConverConf(ctx context.Context, req *pbobjs.ConverIndex) (errs.IMErrorCode, *pbobjs.ConverConf) {
	ret := &pbobjs.ConverConf{}
	appkey := bases.GetAppKeyFromCtx(ctx)
	userId := bases.GetRequesterIdFromCtx(ctx)
	converId := commonservices.GetConversationId(userId, req.TargetId, req.ChannelType)
	converConf := GetConverConf(appkey, converId, req.ChannelType, req.SubChannel)
	if converConf != nil {
		ret.ActivedRtcRoom = converConf.GetActivedRtcRoom(ctx)
	}
	return errs.IMErrorCode_SUCCESS, ret
}

func QryGlobalConverConf(ctx context.Context, req *pbobjs.ConverConfReq) (errs.IMErrorCode, *pbobjs.ConverConf) {
	ret := req.ConverConf
	if ret == nil {
		ret = &pbobjs.ConverConf{}
	}
	appkey := bases.GetAppKeyFromCtx(ctx)
	converConf := GetConverConf(appkey, req.ConverId, req.ChannelType, req.SubChannel)
	if converConf != nil {
		ret.ActivedRtcRoom = converConf.GetActivedRtcRoom(ctx)
		ret.GlobalConverTags = converConf.GetGlobalConverTags()
	}
	return errs.IMErrorCode_SUCCESS, ret
}

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
	rtcStorage "im-server/services/rtcroom/storages"
	"time"

	"google.golang.org/protobuf/proto"
)

type ConverConfItem struct {
	AppKey      string
	ConverId    string
	ChannelType pbobjs.ChannelType
	SubChannel  string

	//rtc room
	RtcRoomId         string
	LatestRtcPingTime int64
}

var converConfCache *caches.LruCache
var converConfLocks *tools.SegmentatedLocks

func init() {
	converConfCache = caches.NewLruCacheWithReadTimeout("converconf_cache", 100000, nil, 10*time.Minute)
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
			storage := storages.NewConverConfStorage()
			confs, err := storage.QryConverConfs(appkey, converId, subChannel, int32(channelType))
			if err == nil && len(confs) > 0 {
				for itemKey, conf := range confs {
					if itemKey == string(models.ConverConfItemKey_RtcRoomId) {
						item.RtcRoomId = conf.ItemValue
						item.LatestRtcPingTime = conf.UpdatedTime
					}
				}
			}
			converConfCache.Add(key, item)
			return item
		}
	}
}

func (conf *ConverConfItem) SetRtcRoomId(rtcRoomId string, t int64) bool {
	l := converConfLocks.GetLocks(getConverConfCacheKey(conf.AppKey, conf.ConverId, conf.ChannelType, conf.SubChannel))
	l.Lock()
	defer l.Unlock()
	if conf.RtcRoomId == "" || conf.RtcRoomId == rtcRoomId {
		conf.RtcRoomId = rtcRoomId
		if t <= 0 {
			t = time.Now().UnixMilli()
		}
		conf.LatestRtcPingTime = t
		return true
	}
	return false
}

func (conf *ConverConfItem) ClearRtcRoomId() {
	l := converConfLocks.GetLocks(getConverConfCacheKey(conf.AppKey, conf.ConverId, conf.ChannelType, conf.SubChannel))
	l.Lock()
	defer l.Unlock()
	conf.RtcRoomId = ""
	conf.LatestRtcPingTime = 0
}

func QryConverConf(ctx context.Context, req *pbobjs.ConverIndex) (errs.IMErrorCode, *pbobjs.ConverConf) {
	ret := &pbobjs.ConverConf{}
	appkey := bases.GetAppKeyFromCtx(ctx)
	userId := bases.GetRequesterIdFromCtx(ctx)
	converId := commonservices.GetConversationId(userId, req.TargetId, req.ChannelType)
	storage := rtcStorage.NewRtcRoomStorage()
	room, err := storage.FindByConver(appkey, converId, req.ChannelType, req.SubChannel)
	if err == nil && room != nil {
		code, resp, err := bases.SyncRpcCall(ctx, "rtc_qry", room.RoomId, &pbobjs.Nil{}, func() proto.Message {
			return &pbobjs.RtcRoom{}
		})
		if err == nil && code == errs.IMErrorCode_SUCCESS && resp != nil {
			rtcRoom := resp.(*pbobjs.RtcRoom)
			ret.ActivedRtcRoom = &pbobjs.ActivedRtcRoom{
				RoomType:     rtcRoom.RoomType,
				RoomId:       rtcRoom.RoomId,
				Owner:        rtcRoom.Owner,
				RtcChannel:   rtcRoom.RtcChannel,
				RtcMediaType: rtcRoom.RtcMediaType,
				Ext:          rtcRoom.Ext,
				Members:      rtcRoom.Members,
			}
		}
	}
	return errs.IMErrorCode_SUCCESS, ret
}

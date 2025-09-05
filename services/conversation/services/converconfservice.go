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

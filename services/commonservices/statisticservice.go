package commonservices

import (
	"fmt"
	"im-server/commons/caches"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/commonservices/dbs"
	"im-server/services/usermanager/storages"
	"sync/atomic"
	"time"
)

type StatType int
type ConnectType int

var (
	statCache           *caches.LruCache
	userActivitiesCache *caches.LruCache
	statLocks           *tools.SegmentatedLocks

	StatType_Up       StatType = 1
	StatType_Dispatch StatType = 2
	StatType_Down     StatType = 3

	ConnectType_Connect     ConnectType = 0
	ConnectType_ChrmConnect ConnectType = 1
)

func init() {
	statCache = caches.NewLruCacheWithAddReadTimeout("msgstat_cache", 1000, func(key, value interface{}) {
		counter := value.(*Counter)
		dao := dbs.MsgStatDao{}
		dao.IncrByStep(counter.Appkey, int(counter.StateType), int(counter.ChannelType), getDbTimeMark(), counter.Count)
	}, 10*time.Minute, 10*time.Minute)
	userActivitiesCache = caches.NewLruCacheWithAddReadTimeout("useractivity_cache", 10000, func(key, value interface{}) {
		counter := value.(*UserActivityCounter)
		dao := dbs.UserActivityDao{}
		dao.IncrByStep(counter.Appkey, counter.UserId, getDbTimeMark(), counter.Count)
	}, 10*time.Minute, 10*time.Minute)
	statLocks = tools.NewSegmentatedLocks(128)
}

type Statistics struct {
	Items          []interface{} `json:"items"`
	TotalUserCount *int64        `json:"total_user_count,omitempty"`
}

type StatisticMsgItem struct {
	Count    int64 `json:"count"`
	TimeMark int64 `json:"time_mark"`
}

type MsgStatistics struct {
	MsgUp       *Statistics `json:"msg_up,omitempty"`
	MsgDown     *Statistics `json:"msg_down,omitempty"`
	MsgDispatch *Statistics `json:"msg_dispatch,omitempty"`
}

func QryMsgStatistic(appkey string, statTypes []StatType, channelType pbobjs.ChannelType, start, end int64) *MsgStatistics {
	ret := &MsgStatistics{}
	intStateTypes := []int{}
	for _, stateType := range statTypes {
		intStateTypes = append(intStateTypes, int(stateType))
		switch stateType {
		case StatType_Up:
			ret.MsgUp = &Statistics{
				Items: []interface{}{},
			}
		case StatType_Down:
			ret.MsgDown = &Statistics{
				Items: []interface{}{},
			}
		case StatType_Dispatch:
			ret.MsgDispatch = &Statistics{
				Items: []interface{}{},
			}
		}
	}
	dao := dbs.MsgStatDao{}
	list := dao.QryStats(appkey, intStateTypes, int(channelType), start, end)
	for _, item := range list {
		switch item.StatType {
		case int(StatType_Up):
			ret.MsgUp.Items = append(ret.MsgUp.Items, &StatisticMsgItem{
				Count:    item.Count,
				TimeMark: item.TimeMark * 1000,
			})
		case int(StatType_Down):
			ret.MsgDown.Items = append(ret.MsgDown.Items, &StatisticMsgItem{
				Count:    item.Count,
				TimeMark: item.TimeMark * 1000,
			})
		case int(StatType_Dispatch):
			ret.MsgDispatch.Items = append(ret.MsgDispatch.Items, &StatisticMsgItem{
				Count:    item.Count,
				TimeMark: item.TimeMark * 1000,
			})
		}
	}
	return ret
}

var oneDay int64 = 24 * 60 * 60 * 1000

type UserActivityItem struct {
	Count    int64 `json:"count"`
	TimeMark int64 `json:"time_mark"`
}

func QryUserActivities(appkey string, start, end int64) *Statistics {
	ret := &Statistics{
		Items: []interface{}{},
	}
	timeMarks := []int64{}
	for s := start / oneDay * oneDay; s <= end; {
		if s >= start {
			timeMarks = append(timeMarks, s)
		}
		s = s + oneDay
	}
	dao := dbs.UserActivityDao{}
	for _, timemark := range timeMarks {
		ret.Items = append(ret.Items, &UserActivityItem{
			TimeMark: timemark,
			Count:    dao.CountUserActivities(appkey, timemark),
		})
	}
	return ret
}

func QryUserRegiste(appkey string, start, end int64) *Statistics {
	ret := &Statistics{
		Items: []interface{}{},
	}
	timeMarks := []int64{}
	for s := start / oneDay * oneDay; s <= end; {
		timeMarks = append(timeMarks, s)
		s = s + oneDay
	}
	storage := storages.NewUserStorage()
	for _, timemark := range timeMarks {
		ret.Items = append(ret.Items, &UserActivityItem{
			TimeMark: timemark,
			Count:    storage.CountByTime(appkey, timemark, timemark+oneDay),
		})
	}
	totalCount := storage.Count(appkey)
	if totalCount > 0 {
		ret.TotalUserCount = tools.Int64Ptr(int64(totalCount))
	}
	return ret
}

func ReportUserLogin(appkey string, userId string) {
	counter := getUserActivityCounter(appkey, userId)
	counter.Incry()
}

func ReportUpMsg(appkey string, channelType pbobjs.ChannelType, step int64) {
	counter := getCounter(appkey, StatType_Up, channelType)
	counter.IncrByStep(step)
}

func ReportDispatchMsg(appkey string, channelType pbobjs.ChannelType, step int64) {
	counter := getCounter(appkey, StatType_Dispatch, channelType)
	counter.IncrByStep(step)
}

func ReportDownMsg(appkey string, channelType pbobjs.ChannelType, step int64) {
	counter := getCounter(appkey, StatType_Down, channelType)
	counter.IncrByStep(step)
}

func getCounter(appkey string, statType StatType, channelType pbobjs.ChannelType) *Counter {
	key := fmt.Sprintf("%s_%d_%d", appkey, channelType, statType)
	if counterObj, exist := statCache.Get(key); exist {
		counter := counterObj.(*Counter)
		return counter
	} else {
		lock := statLocks.GetLocks(key)
		lock.Lock()
		defer lock.Unlock()

		if counterObj, exist := statCache.Get(key); exist {
			counter := counterObj.(*Counter)
			return counter
		} else {
			counter := NewCounter(appkey, int64(statType), int64(channelType))
			statCache.Add(key, counter)
			return counter
		}
	}
}

func getUserActivityCounter(appkey, userId string) *UserActivityCounter {
	key := fmt.Sprintf("%s_%s", appkey, userId)
	if obj, exist := userActivitiesCache.Get(key); exist {
		return obj.(*UserActivityCounter)
	} else {
		lock := statLocks.GetLocks(key)
		lock.Lock()
		defer lock.Unlock()

		if obj, exist := userActivitiesCache.Get(key); exist {
			return obj.(*UserActivityCounter)
		} else {
			counter := &UserActivityCounter{
				Appkey: appkey,
				UserId: userId,
				Count:  0,
			}
			userActivitiesCache.Add(key, counter)
			return counter
		}
	}
}

type Counter struct {
	Appkey      string
	StateType   int64
	ChannelType int64
	Count       int64
}

func NewCounter(appkey string, stateType, channelType int64) *Counter {
	return &Counter{
		Appkey:      appkey,
		StateType:   stateType,
		ChannelType: channelType,
		Count:       0,
	}
}

func (c *Counter) Incry() {
	c.IncrByStep(1)
}

func (c *Counter) IncrByStep(step int64) {
	atomic.AddInt64(&c.Count, step)
}

func getDbTimeMark() int64 {
	current := time.Now().Unix()
	var day int64 = 24 * 60 * 60
	return current / day * day
}

type UserActivityCounter struct {
	Appkey string
	UserId string
	Count  int64
}

func (c *UserActivityCounter) Incry() {
	atomic.AddInt64(&c.Count, 1)
}

func ReportConcurrentConnectCount(appkey string, count int64) {
	current := time.Now().UnixMilli()
	timeMark := current / 30000 * 30000
	dao := dbs.ConnectCountDao{}
	dao.IncrByStep(appkey, int(ConnectType_Connect), timeMark, count)
}

type ConcurrentConnectItem struct {
	TimeMark int64 `json:"time_mark"`
	Count    int64 `json:"count"`
}

func QryConnect(appkey string, start, end int64) *Statistics {
	ret := &Statistics{
		Items: []interface{}{},
	}
	dao := dbs.ConnectCountDao{}
	list := dao.QryStats(appkey, int(ConnectType_Connect), start, end)
	for _, item := range list {
		ret.Items = append(ret.Items, &ConcurrentConnectItem{
			TimeMark: item.TimeMark,
			Count:    item.Count,
		})
	}
	return ret
}

func QryMaxConnect(appkey string, start, end int64) *Statistics {
	ret := &Statistics{
		Items: []interface{}{},
	}
	timeMarks := []int64{}
	for s := start / oneDay * oneDay; s <= end; {
		timeMarks = append(timeMarks, s)
		s = s + oneDay
	}
	dao := dbs.ConnectCountDao{}
	for _, timemark := range timeMarks {
		item := dao.MaxByTime(appkey, int(ConnectType_Connect), timemark, timemark+oneDay)
		var count int64 = 0
		if item != nil {
			count = item.Count
		}
		ret.Items = append(ret.Items, &ConcurrentConnectItem{
			TimeMark: timemark,
			Count:    count,
		})
	}
	return ret
}

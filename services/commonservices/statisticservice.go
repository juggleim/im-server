package commonservices

import (
	"fmt"
	"im-server/commons/caches"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/commonservices/dbs"
	userDbs "im-server/services/usermanager/dbs"
	"sync/atomic"
	"time"
)

type StatType int

var (
	statCache           *caches.LruCache
	userActivitiesCache *caches.LruCache
	statLocks           *tools.SegmentatedLocks

	StatType_Up       StatType = 1
	StatType_Dispatch StatType = 2
	StatType_Down     StatType = 3
)

func init() {
	statCache = caches.NewLruCacheWithAddReadTimeout(1000, func(key, value interface{}) {
		counter := value.(*Counter)
		dao := dbs.MsgStatDao{}
		dao.IncrByStep(counter.Appkey, int(counter.StateType), int(counter.ChannelType), getDbTimeMark(), counter.Count)
	}, 10*time.Minute, 10*time.Minute)
	userActivitiesCache = caches.NewLruCacheWithAddReadTimeout(10000, func(key, value interface{}) {
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

func QryMsgStatistic(appkey string, statType StatType, channelType pbobjs.ChannelType, start, end int64) *Statistics {
	ret := &Statistics{
		Items: []interface{}{},
	}
	dao := dbs.MsgStatDao{}
	list := dao.QryStats(appkey, int(statType), int(channelType), start, end)
	for _, item := range list {
		ret.Items = append(ret.Items, &StatisticMsgItem{
			Count:    item.Count,
			TimeMark: item.TimeMark,
		})
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
	dao := userDbs.UserDao{}
	for _, timemark := range timeMarks {
		ret.Items = append(ret.Items, &UserActivityItem{
			TimeMark: timemark,
			Count:    dao.CountByTime(appkey, timemark, timemark+oneDay),
		})
	}
	totalCount := dao.Count(appkey)
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

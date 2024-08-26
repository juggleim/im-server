package commonservices

import (
	"fmt"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/commonservices/dbs"
	"sync"
	"time"
)

type StatType int

var (
	statCache *sync.Map
	statLocks *tools.SegmentatedLocks

	StatType_Up       StatType = 1
	StatType_Dispatch StatType = 2
	StatType_Down     StatType = 3
)

func init() {
	statCache = &sync.Map{}
	statLocks = tools.NewSegmentatedLocks(8)
}

func ReportUpMsg(appkey string, channelType pbobjs.ChannelType, step int64) {
	// counter := getCounter(appkey, StatType_Up, channelType)
	// counter.IncrByStep(step)
}

func ReportDispatchMsg(appkey string, channelType pbobjs.ChannelType, step int64) {
	// counter := getCounter(appkey, StatType_Dispatch, channelType)
	// counter.IncrByStep(step)
}

func ReportDownMsg(appkey string, channelType pbobjs.ChannelType, step int64) {
	// counter := getCounter(appkey, StatType_Down, channelType)
	// counter.IncrByStep(step)
}

func getCounter(appkey string, statType StatType, channelType pbobjs.ChannelType) *Counter {
	key := fmt.Sprintf("%s_%d_%d", appkey, channelType, statType)
	if counterObj, exist := statCache.Load(key); exist {
		counter := counterObj.(*Counter)
		return counter
	} else {
		lock := statLocks.GetLocks(key)
		lock.Lock()
		defer lock.Unlock()

		if counterObj, exist := statCache.Load(key); exist {
			counter := counterObj.(*Counter)
			return counter
		} else {
			counter := NewCounter(func(count, timeMark int64) {
				if statType == StatType_Up {
					dao := dbs.UpStatDao{}
					dao.IncrByStep(appkey, int(channelType), timeMark, count)
				} else if statType == StatType_Dispatch {
					dao := dbs.DispStatDao{}
					dao.IncrByStep(appkey, int(channelType), timeMark, count)
				} else if statType == StatType_Down {
					dao := dbs.DownStatDao{}
					dao.IncrByStep(appkey, int(channelType), timeMark, count)
				}
			})
			statCache.Store(key, counter)
			return counter
		}
	}
}

type Counter struct {
	Count int64

	interval int64 //second
	timeMark int64
	lock     sync.RWMutex
	report   func(count, timeMark int64)
}

func NewCounter(report func(count, timeMark int64)) *Counter {
	return &Counter{
		Count:    0,
		interval: 30,
		report:   report,
	}
}

func (c *Counter) Incry() {
	c.IncrByStep(1)
}

func (c *Counter) IncrByStep(step int64) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.timeMark == 0 {
		c.timeMark = c.getTimeMark()
	} else {
		newTimeMark := c.getTimeMark()
		if newTimeMark > c.timeMark {
			if c.report != nil {
				go c.report(c.Count, c.timeMark)
			}
			c.timeMark = newTimeMark
			c.Count = 0
		}
	}

	c.Count = c.Count + step
}

func (c *Counter) getTimeMark() int64 {
	current := time.Now().Unix()
	return current / c.interval * c.interval
}

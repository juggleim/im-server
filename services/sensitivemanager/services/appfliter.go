package services

import (
	"fmt"
	"im-server/commons/caches"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/sensitivemanager/dbs"
	"time"

	"github.com/samber/lo"
)

var (
	filterCache *caches.LruCache
	filterLocks *tools.SegmentatedLocks
)

func init() {
	filterCache = caches.NewLruCacheWithAddReadTimeout("filter_cache", 10000, nil, 8*time.Minute, 10*time.Minute)
	filterLocks = tools.NewSegmentatedLocks(128)

	filterCache.SetValueCreator(func(key interface{}) interface{} {
		s := NewSensitiveService()
		start := time.Now()
		loadAppWords(s, key.(string))
		fmt.Println("load app words cost:", time.Since(start))
		return s
	})
}

func GetAppFilter(appKey string) *SensitiveService {
	lock := filterLocks.GetLocks(appKey)
	lock.Lock()
	defer lock.Unlock()

	v, ok := filterCache.GetByCreator(appKey, nil)
	if !ok {
		return nil
	}
	return v.(*SensitiveService)
}

func loadAppWords(service *SensitiveService, appKey string) (err error) {
	var (
		startId  int64 = 0
		pageSize int64 = 1000
	)
	dao := dbs.SensitiveWordDao{}
	for {
		list, err := dao.QrySensitiveWords(appKey, pageSize, startId)
		if err != nil {
			return err
		}
		words := lo.Map(list, func(item *dbs.SensitiveWordDao, index int) *pbobjs.SensitiveWord {
			if startId < item.ID {
				startId = item.ID
			}
			return &pbobjs.SensitiveWord{
				Word:     item.Word,
				WordType: pbobjs.SensitiveWordType(item.WordType),
			}
		})

		service.AddWord(words...)
		if len(list) < int(pageSize) {
			break
		}
	}
	return nil
}

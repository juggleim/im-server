package sensitive

import (
	"fmt"
	"im-server/commons/caches"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/commonservices/logs"
	"im-server/services/sensitivemanager/dbs"
	"sync"
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

func GetAppSensitiveFilter(appKey string) *SensitiveService {
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
			logs.NewLogEntity().Error(err.Error())
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

type SensitiveService struct {
	replaceFilter *Filter
	denyFilter    *Filter
	loadLock      *sync.RWMutex
}

func NewSensitiveService() *SensitiveService {
	return &SensitiveService{
		replaceFilter: NewFilter(),
		denyFilter:    NewFilter(),
		loadLock:      &sync.RWMutex{},
	}
}

func (s *SensitiveService) ReplaceSensitiveWords(text string) (isDeny bool, replacedText string) {
	s.loadLock.RLock()
	defer s.loadLock.RUnlock()

	if s.denyFilter != nil {
		var ok bool
		ok, _ = s.denyFilter.FindIn(text)
		if ok {
			isDeny = true
			return
		}
	}
	if s.replaceFilter != nil {
		replacedText = s.replaceFilter.Replace(text, '*')
	}

	return
}

func (s *SensitiveService) AddWord(words ...*pbobjs.SensitiveWord) {
	s.loadLock.Lock()
	defer s.loadLock.Unlock()
	if s.replaceFilter == nil {
		s.replaceFilter = NewFilter()
	}
	if s.denyFilter == nil {
		s.denyFilter = NewFilter()
	}
	for _, word := range words {
		if word.WordType == pbobjs.SensitiveWordType_deny_word {
			s.denyFilter.AddWord(word.Word)
		} else {
			s.replaceFilter.AddWord(word.Word)
		}
	}
}

func (s *SensitiveService) DelWord(words ...string) {
	s.loadLock.Lock()
	defer s.loadLock.Unlock()
	if s.denyFilter != nil {
		s.denyFilter.DelWord(words...)
	}
	if s.replaceFilter != nil {
		s.replaceFilter.DelWord(words...)
	}
}

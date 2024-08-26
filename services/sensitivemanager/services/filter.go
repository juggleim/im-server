package services

import (
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/sensitivemanager/sensitive"
	"sync"
)

type SensitiveService struct {
	replaceFilter *sensitive.Filter
	denyFilter    *sensitive.Filter
	loadLock      *sync.RWMutex
}

func NewSensitiveService() *SensitiveService {
	return &SensitiveService{
		replaceFilter: sensitive.New(),
		denyFilter:    sensitive.New(),
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
		s.replaceFilter = sensitive.New()
	}
	if s.denyFilter == nil {
		s.denyFilter = sensitive.New()
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

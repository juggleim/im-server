package services

import (
	"fmt"
	"im-server/commons/caches"
	"im-server/commons/errs"
	"im-server/services/conversation/storages"
	"im-server/services/conversation/storages/models"
	"sort"
	"time"
)

var userConverTagsCache *caches.LruCache

func init() {
	userConverTagsCache = caches.NewLruCacheWithAddReadTimeout("user_conver_tags_cache", 100000, nil, time.Hour, time.Hour)
}

type UserConverTags struct {
	Appkey string
	UserId string

	TagMap map[string]*models.UserConverTag
}

func cloneUserConverTag(tag *models.UserConverTag) *models.UserConverTag {
	if tag == nil {
		return nil
	}
	ret := *tag
	return &ret
}

func getUserConverTagsCacheKey(appkey, userId string) string {
	return fmt.Sprintf("%s_%s", appkey, userId)
}

func sortUserConverTags(tags []*models.UserConverTag) {
	if len(tags) <= 1 {
		return
	}
	sort.Slice(tags, func(i, j int) bool {
		if tags[i].TagOrder == tags[j].TagOrder {
			return tags[i].CreatedTime < tags[j].CreatedTime
		}
		return tags[i].TagOrder < tags[j].TagOrder
	})
}

func GetUserConverTags(appkey, userId string) *UserConverTags {
	cacheKey := getUserConverTagsCacheKey(appkey, userId)
	if obj, exist := userConverTagsCache.Get(cacheKey); exist {
		return obj.(*UserConverTags)
	} else {
		l := userLocks.GetLocks(cacheKey)
		l.Lock()
		defer l.Unlock()

		if obj, exist := userConverTagsCache.Get(cacheKey); exist {
			return obj.(*UserConverTags)
		}
		userConverTags := &UserConverTags{
			Appkey: appkey,
			UserId: userId,
			TagMap: make(map[string]*models.UserConverTag),
		}
		storage := storages.NewUserConverTagStorage()
		dbTags, err := storage.QryTags(appkey, userId)
		if err == nil {
			tagMap := make(map[string]*models.UserConverTag)
			for _, tag := range dbTags {
				tagMap[tag.Tag] = &models.UserConverTag{
					AppKey:      tag.AppKey,
					UserId:      tag.UserId,
					Tag:         tag.Tag,
					TagName:     tag.TagName,
					TagOrder:    tag.TagOrder,
					CreatedTime: tag.CreatedTime,
				}
			}
			userConverTags.TagMap = tagMap
		}
		userConverTagsCache.Add(cacheKey, userConverTags)
		return userConverTags
	}
}

func (uc *UserConverTags) AddTagsWithBackup(tags []models.UserConverTag, maxCount int) ([]models.UserConverTag, map[string]*models.UserConverTag, errs.IMErrorCode) {
	l := userLocks.GetLocks(getUserConverTagsCacheKey(uc.Appkey, uc.UserId))
	l.Lock()
	defer l.Unlock()
	if len(tags) == 0 {
		return []models.UserConverTag{}, map[string]*models.UserConverTag{}, errs.IMErrorCode_SUCCESS
	}

	// Pre-check added count; if exceed maxCount, reject all changes.
	needAddSet := make(map[string]struct{})
	for _, item := range tags {
		if _, exist := uc.TagMap[item.Tag]; !exist {
			needAddSet[item.Tag] = struct{}{}
		}
	}
	if len(uc.TagMap)+len(needAddSet) > maxCount {
		return []models.UserConverTag{}, map[string]*models.UserConverTag{}, errs.IMErrorCode_CONVER_TAGEXCEEDMAXCOUNT
	}

	changed := make([]models.UserConverTag, 0, len(tags))
	changedIdx := make(map[string]int)
	backup := make(map[string]*models.UserConverTag)
	ensureBackup := func(tagKey string) {
		if _, ok := backup[tagKey]; ok {
			return
		}
		if old, exist := uc.TagMap[tagKey]; exist {
			backup[tagKey] = cloneUserConverTag(old)
		} else {
			backup[tagKey] = nil
		}
	}
	for _, item := range tags {
		newItem := models.UserConverTag{
			AppKey:      item.AppKey,
			UserId:      item.UserId,
			Tag:         item.Tag,
			TagName:     item.TagName,
			TagOrder:    item.TagOrder,
			CreatedTime: item.CreatedTime,
			IsAdd:       false,
		}
		cacheTag, exist := uc.TagMap[item.Tag]
		if exist {
			if cacheTag.TagName == item.TagName && cacheTag.TagOrder == item.TagOrder {
				continue
			}
			ensureBackup(item.Tag)
			newItem.IsAdd = false
			uc.TagMap[item.Tag] = &newItem

			if idx, ok := changedIdx[item.Tag]; ok {
				changed[idx] = newItem
			} else {
				changedIdx[item.Tag] = len(changed)
				changed = append(changed, newItem)
			}
			continue
		}
		// New tag
		ensureBackup(item.Tag)
		newItem.IsAdd = true
		uc.TagMap[item.Tag] = &newItem
		if idx, ok := changedIdx[item.Tag]; ok {
			changed[idx] = newItem
		} else {
			changedIdx[item.Tag] = len(changed)
			changed = append(changed, newItem)
		}
	}
	return changed, backup, errs.IMErrorCode_SUCCESS
}

func (uc *UserConverTags) RollbackTags(backup map[string]*models.UserConverTag) {
	if len(backup) == 0 {
		return
	}
	l := userLocks.GetLocks(getUserConverTagsCacheKey(uc.Appkey, uc.UserId))
	l.Lock()
	defer l.Unlock()
	for tagKey, old := range backup {
		if old == nil {
			delete(uc.TagMap, tagKey)
		} else {
			uc.TagMap[tagKey] = cloneUserConverTag(old)
		}
	}
}

func (uc *UserConverTags) DelTag(tag string) bool {
	l := userLocks.GetLocks(getUserConverTagsCacheKey(uc.Appkey, uc.UserId))
	l.Lock()
	defer l.Unlock()
	if _, exist := uc.TagMap[tag]; !exist {
		return false
	}
	delete(uc.TagMap, tag)
	return true
}

func (uc *UserConverTags) QryTags() []*models.UserConverTag {
	l := userLocks.GetLocks(getUserConverTagsCacheKey(uc.Appkey, uc.UserId))
	l.RLock()
	defer l.RUnlock()
	tags := make([]*models.UserConverTag, 0, len(uc.TagMap))
	for _, tag := range uc.TagMap {
		tags = append(tags, tag)
	}
	sortUserConverTags(tags)
	return tags
}

func (uc *UserConverTags) ContainsTag(tag string) bool {
	l := userLocks.GetLocks(getUserConverTagsCacheKey(uc.Appkey, uc.UserId))
	l.RLock()
	defer l.RUnlock()
	_, exist := uc.TagMap[tag]
	return exist
}

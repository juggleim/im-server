package services

import (
	"fmt"
	"im-server/commons/caches"
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

func (uc *UserConverTags) AddTag(tag *models.UserConverTag) bool {
	l := userLocks.GetLocks(getUserConverTagsCacheKey(uc.Appkey, uc.UserId))
	l.Lock()
	defer l.Unlock()
	if cacheTag, exist := uc.TagMap[tag.Tag]; exist {
		if cacheTag.TagName == tag.TagName && cacheTag.TagOrder == tag.TagOrder {
			return false
		}
	}
	uc.TagMap[tag.Tag] = tag
	return true
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

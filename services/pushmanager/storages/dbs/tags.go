package dbs

import (
	"context"
	"github.com/jinzhu/gorm"
	"im-server/services/pushmanager/storages/models"
	"time"
)

var _ models.IUserTagStorage = (*UserTagsDao)(nil)

type UserTag struct {
	ID          int       `gorm:"primary_key"`
	AppKey      string    `gorm:"app_key"`
	UserID      string    `gorm:"user_id"`
	Tag         string    `gorm:"tag"`
	CreatedTime time.Time `gorm:"created_time"`
}

func (UserTag) TableName() string {
	return "user_tags"
}

type UserTagsDao struct {
	db *gorm.DB
}

func NewUserTagsDao(db *gorm.DB) *UserTagsDao {
	return &UserTagsDao{db: db}
}

func (u *UserTagsDao) AddUserTags(ctx context.Context, appKey string, userID string, tags ...string) error {
	sql := "INSERT INTO user_tags (app_key, user_id, tag) VALUES "

	for i, tag := range tags {
		if i == len(tags)-1 {
			sql += "('" + appKey + "', '" + userID + "', '" + tag + "')"
		} else {
			sql += "('" + appKey + "', '" + userID + "', '" + tag + "'), "

		}
	}
	return u.db.Exec(sql).Error
}

func (u *UserTagsDao) DeleteUserTags(ctx context.Context, appKey string, userID string, tags ...string) error {
	return u.db.Where("app_key = ? AND user_id = ? AND tag IN (?)", appKey, userID, tags).Delete(&UserTag{}).Error
}

func (u *UserTagsDao) ClearUserTag(ctx context.Context, appKey string, userID string) error {
	return u.db.Where("app_key = ? AND user_id = ?", appKey, userID).Delete(UserTag{}).Error
}

func (u *UserTagsDao) GetUserTags(ctx context.Context, appKey string, userID string) ([]string, error) {
	var tags []string
	err := u.db.Model(&UserTag{}).Where("app_key = ? AND user_id = ?", appKey, userID).Pluck("tag", &tags).Error
	return tags, err
}

func (u *UserTagsDao) GetUserWithTags(ctx context.Context, appKey string, condition models.Condition, page int, perPage int) (userIDs []string, err error) {
	var (
		results []UserTag
		query   *gorm.DB
	)

	if len(condition.TagsAnd) > 0 {
		query = u.db.Table("user_tags").
			Select("user_id").
			Group("user_id").
			Having("COUNT(DISTINCT tag) = ?", len(condition.TagsAnd)).
			Where("app_key = ?", appKey).
			Where("tag IN (?)", condition.TagsAnd).
			Offset((page - 1) * perPage).
			Limit(perPage)
	} else {
		query = u.db.Table("user_tags").Where("app_key = ?", appKey).
			Where("tag IN (?)", condition.TagsOr).
			Offset((page - 1) * perPage).
			Limit(perPage)
	}

	err = query.Scan(&results).Error
	if err != nil {
		return nil, err
	}

	userIDMap := make(map[string]struct{})
	for _, result := range results {
		userIDMap[result.UserID] = struct{}{}
	}

	for userID := range userIDMap {
		userIDs = append(userIDs, userID)
	}

	return userIDs, nil
}

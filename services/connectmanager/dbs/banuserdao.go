package dbs

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"im-server/commons/dbcommons"
	"sort"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserBanScope string

const (
	UserBanScopeDefault  UserBanScope = "default"
	UserBanScopePlatform UserBanScope = "platform"
	UserBanScopeDevice   UserBanScope = "device"
	UserBanScopeIp       UserBanScope = "ip"
)

type BanUserDao struct {
	ID             int64     `gorm:"primary_key"`
	UserId         string    `gorm:"user_id"`
	CreatedTime    time.Time `gorm:"created_time"`
	EndTime        int64     `gorm:"end_time"`
	ScopeKey       string    `gorm:"scope_key"`
	ScopeValue     string    `gorm:"scope_value"`
	ScopeValueHash string    `gorm:"scope_value_hash"`
	Ext            string    `gorm:"ext"`
	AppKey         string    `gorm:"app_key"`
}

func (user BanUserDao) TableName() string {
	return "banuserrules"
}

func (user BanUserDao) Upsert(item BanUserDao) error {
	return user.UpsertMany([]BanUserDao{item})
}

func (user BanUserDao) UpsertMany(items []BanUserDao) error {
	return dbcommons.GetDb().Transaction(func(tx *gorm.DB) error {
		type legacyScopeKey struct {
			appkey, userId, scopeKey string
		}
		legacyValues := make(map[legacyScopeKey][]string)
		for _, item := range items {
			if err := upsertBanUserRule(tx, item); err != nil {
				return err
			}
			key := legacyScopeKey{appkey: item.AppKey, userId: item.UserId, scopeKey: item.ScopeKey}
			legacyValues[key] = append(legacyValues[key], item.ScopeValue)
		}
		for key, scopeValues := range legacyValues {
			if err := deleteLegacyBanValues(tx, key.appkey, key.userId, key.scopeKey, scopeValues); err != nil {
				return err
			}
		}
		return nil
	})
}

func upsertBanUserRule(db *gorm.DB, item BanUserDao) error {
	if item.CreatedTime.IsZero() {
		item.CreatedTime = time.Now()
	}
	item.ScopeValueHash = hashBanScopeValue(item.ScopeValue)
	return db.Exec("INSERT INTO banuserrules (user_id, created_time, end_time, scope_key, scope_value, scope_value_hash, ext, app_key)VALUES(?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE end_time=VALUES(end_time), scope_value=VALUES(scope_value), ext=VALUES(ext)",
		item.UserId, item.CreatedTime, item.EndTime, item.ScopeKey, item.ScopeValue, item.ScopeValueHash, item.Ext, item.AppKey).Error
}

func (user BanUserDao) FindById(appkey, userId string) ([]*BanUserDao, error) {
	var items []*BanUserDao
	err := dbcommons.GetDb().Where("app_key=? and user_id=?", appkey, userId).Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (user BanUserDao) FindByUserIds(appkey string, userIds []string) ([]*BanUserDao, error) {
	var items []*BanUserDao
	err := dbcommons.GetDb().Where("app_key=? and user_id in (?) and (end_time=0 or end_time>?)", appkey, userIds, time.Now().UnixMilli()).Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (user BanUserDao) DelBanUser(appkey, userId, scopeKey string, scopeValues []string) error {
	if scopeKey == "" {
		return dbcommons.GetDb().Where("app_key=? and user_id=?", appkey, userId).Delete(&BanUserDao{}).Error
	}
	query := dbcommons.GetDb().Where("app_key=? and user_id=? and scope_key=?", appkey, userId, scopeKey)
	if len(scopeValues) > 0 {
		hashes := make([]string, 0, len(scopeValues))
		for _, scopeValue := range scopeValues {
			hashes = append(hashes, hashBanScopeValue(scopeValue))
		}
		query = query.Where("scope_value_hash in (?)", hashes)
	}
	return query.Delete(&BanUserDao{}).Error
}

func (user BanUserDao) CleanBaseTime(appkey, userId string, endTime int64) error {
	return dbcommons.GetDb().Where("app_key=? and user_id=? and end_time>0 and end_time<=?", appkey, userId, endTime).Delete(&BanUserDao{}).Error
}

func (user BanUserDao) QryBanUsers(appkey string, limit, startId int64) ([]*BanUserDao, error) {
	var items []*BanUserDao
	err := dbcommons.GetDb().Where("app_key=? and (end_time=0 or end_time>?) and id>?", appkey, time.Now().UnixMilli(), startId).Order("id asc").Limit(int(limit)).Find(&items).Error
	return items, err
}

func hashBanScopeValue(scopeValue string) string {
	sum := sha256.Sum256([]byte(scopeValue))
	return hex.EncodeToString(sum[:])
}

type LegacyBanUserDao struct {
	ID          int64     `gorm:"primary_key"`
	UserId      string    `gorm:"user_id"`
	CreatedTime time.Time `gorm:"created_time"`
	EndTime     int64     `gorm:"end_time"`
	ScopeKey    string    `gorm:"scope_key"`
	ScopeValue  string    `gorm:"scope_value"`
	Ext         string    `gorm:"ext"`
	AppKey      string    `gorm:"app_key"`
}

func (user LegacyBanUserDao) TableName() string {
	return "banusers"
}

func (user LegacyBanUserDao) FindById(appkey, userId string) ([]*LegacyBanUserDao, error) {
	var items []*LegacyBanUserDao
	err := dbcommons.GetDb().Where("app_key=? and user_id=?", appkey, userId).Find(&items).Error
	return items, err
}

func (user LegacyBanUserDao) FindByUserIds(appkey string, userIds []string) ([]*LegacyBanUserDao, error) {
	var items []*LegacyBanUserDao
	err := dbcommons.GetDb().Where("app_key=? and user_id in (?) and (end_time=0 or end_time>?)", appkey, userIds, time.Now().UnixMilli()).Find(&items).Error
	return items, err
}

func (user LegacyBanUserDao) QryBanUsers(appkey string, limit, startId int64) ([]*LegacyBanUserDao, error) {
	var items []*LegacyBanUserDao
	err := dbcommons.GetDb().Where("app_key=? and (end_time=0 or end_time>?) and id>?", appkey, time.Now().UnixMilli(), startId).Order("id asc").Limit(int(limit)).Find(&items).Error
	return items, err
}

func (user LegacyBanUserDao) DelBanUser(appkey, userId, scopeKey string, scopeValues []string) error {
	return dbcommons.GetDb().Transaction(func(tx *gorm.DB) error {
		return deleteLegacyBanValues(tx, appkey, userId, scopeKey, scopeValues)
	})
}

func deleteLegacyBanValues(tx *gorm.DB, appkey, userId, scopeKey string, scopeValues []string) error {
	query := tx.Where("app_key=? and user_id=?", appkey, userId)
	if scopeKey == "" {
		return query.Delete(&LegacyBanUserDao{}).Error
	}
	if scopeKey == string(UserBanScopeDefault) {
		query = query.Where("scope_key in (?)", []string{"", string(UserBanScopeDefault)})
	} else {
		query = query.Where("scope_key=?", scopeKey)
	}
	if len(scopeValues) == 0 {
		return query.Delete(&LegacyBanUserDao{}).Error
	}
	var item LegacyBanUserDao
	err := query.Clauses(clause.Locking{Strength: "UPDATE"}).First(&item).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	if err != nil {
		return err
	}
	removed := make(map[string]struct{}, len(scopeValues))
	for _, scopeValue := range scopeValues {
		removed[scopeValue] = struct{}{}
	}
	remaining := make([]string, 0)
	for _, scopeValue := range SplitLegacyBanScopeValues(scopeKey, item.ScopeValue) {
		if _, exists := removed[scopeValue]; !exists {
			remaining = append(remaining, scopeValue)
		}
	}
	if len(remaining) == 0 {
		return tx.Where("id=?", item.ID).Delete(&LegacyBanUserDao{}).Error
	}
	return tx.Model(&LegacyBanUserDao{}).Where("id=?", item.ID).Update("scope_value", strings.Join(remaining, ",")).Error
}

func (user LegacyBanUserDao) CleanBaseTime(appkey, userId string, endTime int64) error {
	return dbcommons.GetDb().Where("app_key=? and user_id=? and end_time>0 and end_time<=?", appkey, userId, endTime).Delete(&LegacyBanUserDao{}).Error
}

func SplitLegacyBanScopeValues(scopeKey, scopeValue string) []string {
	if scopeKey == string(UserBanScopeDefault) || scopeKey == "" {
		return []string{""}
	}
	values := make(map[string]struct{})
	for _, value := range strings.Split(scopeValue, ",") {
		value = strings.TrimSpace(value)
		if value != "" {
			values[value] = struct{}{}
		}
	}
	result := make([]string, 0, len(values))
	for value := range values {
		result = append(result, value)
	}
	sort.Strings(result)
	return result
}

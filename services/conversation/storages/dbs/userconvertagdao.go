package dbs

import (
	"fmt"
	"im-server/commons/dbcommons"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/conversation/storages/models"
	"strings"
	"time"
)

type UserConverTagDao struct {
	ID          int64     `gorm:"primary_key"`
	UserId      string    `gorm:"user_id"`
	Tag         string    `gorm:"tag"`
	TagName     string    `gorm:"tag_name"`
	CreatedTime time.Time `gorm:"created_time"`
	AppKey      string    `gorm:"app_key"`
}

func (utag *UserConverTagDao) TableName() string {
	return "userconvertags"
}

func (utag *UserConverTagDao) Upsert(item models.UserConverTag) error {
	if item.TagName != "" {
		sql := fmt.Sprintf("INSERT INTO %s (app_key,user_id,tag,tag_name) VALUES (?,?,?,?) ON DUPLICATE KEY UPDATE tag_name=?", utag.TableName())
		return dbcommons.GetDb().Exec(sql, item.AppKey, item.UserId, item.Tag, item.TagName, item.TagName).Error
	} else {
		sql := fmt.Sprintf("INSERT IGNORE INTO %s (app_key,user_id,tag) VALUES (?,?,?)", utag.TableName())
		return dbcommons.GetDb().Exec(sql, item.AppKey, item.UserId, item.Tag, item.TagName, item.TagName).Error
	}
}

func (utag *UserConverTagDao) Delete(appkey, userId, tag string) error {
	return dbcommons.GetDb().Where("app_key=? and user_id=? and tag=?", appkey, userId, tag).Delete(&UserConverTagDao{}).Error
}

func (utag *UserConverTagDao) QryTags(appkey, userId string) ([]*models.UserConverTag, error) {
	var items []*UserConverTagDao
	ret := []*models.UserConverTag{}
	err := dbcommons.GetDb().Where("app_key=? and user_id=?", appkey, userId).Find(&items).Error
	if err == nil {
		for _, item := range items {
			ret = append(ret, &models.UserConverTag{
				UserId:      item.UserId,
				Tag:         item.Tag,
				TagName:     item.TagName,
				CreatedTime: item.CreatedTime.UnixMilli(),
				AppKey:      item.AppKey,
			})
		}
	}
	return ret, err
}

func (utag *UserConverTagDao) QryTagsByConver(appkey, userId, targetId string, channelType pbobjs.ChannelType) ([]*models.UserConverTag, error) {
	var items []*UserConverTagDao

	tagRel := &ConverTagRelDao{}
	params := []interface{}{}
	var sqlBuilder strings.Builder
	sqlBuilder.WriteString("select ")
	sqlBuilder.WriteString(utag.TableName())
	sqlBuilder.WriteString(".tag,")
	sqlBuilder.WriteString(utag.TableName())
	sqlBuilder.WriteString(".tag_name from ")
	sqlBuilder.WriteString(utag.TableName())
	sqlBuilder.WriteString(" right join ")
	sqlBuilder.WriteString(tagRel.TableName())
	sqlBuilder.WriteString(" on (")
	sqlBuilder.WriteString(utag.TableName())
	sqlBuilder.WriteString(".app_key=")
	sqlBuilder.WriteString(tagRel.TableName())
	sqlBuilder.WriteString(".app_key and ")
	sqlBuilder.WriteString(utag.TableName())
	sqlBuilder.WriteString(".user_id=")
	sqlBuilder.WriteString(tagRel.TableName())
	sqlBuilder.WriteString(".user_id and ")
	sqlBuilder.WriteString(utag.TableName())
	sqlBuilder.WriteString(".tag=")
	sqlBuilder.WriteString(tagRel.TableName())
	sqlBuilder.WriteString(".tag) where ")

	sqlBuilder.WriteString(tagRel.TableName())
	sqlBuilder.WriteString(".app_key=? and ")
	params = append(params, appkey)

	sqlBuilder.WriteString(tagRel.TableName())
	sqlBuilder.WriteString(".user_id=? and ")
	params = append(params, userId)

	sqlBuilder.WriteString(tagRel.TableName())
	sqlBuilder.WriteString(".target_id=? and ")
	params = append(params, targetId)

	sqlBuilder.WriteString(tagRel.TableName())
	sqlBuilder.WriteString(".channel_type=?")
	params = append(params, channelType)

	err := dbcommons.GetDb().Raw(sqlBuilder.String(), params...).Find(&items).Error
	if err != nil {
		return []*models.UserConverTag{}, err
	}
	ret := []*models.UserConverTag{}
	for _, tag := range items {
		ret = append(ret, &models.UserConverTag{
			UserId:  userId,
			AppKey:  appkey,
			Tag:     tag.Tag,
			TagName: tag.TagName,
		})
	}
	return ret, nil
}

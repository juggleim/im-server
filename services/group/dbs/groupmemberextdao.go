package dbs

import (
	"bytes"
	"fmt"
	"im-server/commons/dbcommons"
	"time"
)

type GroupMemberExtDao struct {
	ID          int64     `gorm:"primary_key"`
	GroupId     string    `gorm:"group_id"`
	MemberId    string    `gorm:"member_id"`
	ItemKey     string    `gorm:"item_key"`
	ItemValue   string    `gorm:"item_value"`
	ItemType    int       `gorm:"item_type"`
	UpdatedTime time.Time `gorm:"updated_time"`
	AppKey      string    `gorm:"app_key"`
}

func (ext GroupMemberExtDao) TableName() string {
	return "groupmemberexts"
}

func (ext GroupMemberExtDao) BatchCreate(items []GroupMemberExtDao) error {
	var buffer bytes.Buffer
	sql := fmt.Sprintf("insert into %s (`app_key`,`group_id`,`member_id`,`item_key`,`item_value`,`item_type`)values", ext.TableName())
	params := []interface{}{}

	buffer.WriteString(sql)
	for i, item := range items {
		if i == len(items)-1 {
			buffer.WriteString("(?,?,?,?,?,?);")
		} else {
			buffer.WriteString("(?,?,?,?,?,?),")
		}
		params = append(params, item.AppKey, item.GroupId, item.MemberId, item.ItemKey, item.ItemValue, item.ItemType)
	}

	err := dbcommons.GetDb().Exec(buffer.String(), params...).Error
	return err
}

func (ext GroupMemberExtDao) QryExtFields(appkey, groupId, memberId string) ([]*GroupMemberExtDao, error) {
	var items []*GroupMemberExtDao
	err := dbcommons.GetDb().Where("app_key=? and group_id=? and member_id=?", appkey, groupId, memberId).Find(&items).Error
	return items, err
}

func (ext GroupMemberExtDao) QryExtFieldsByMemberIds(appkey, groupId string, memberIds []string) (map[string][]*GroupMemberExtDao, error) {
	var items []*GroupMemberExtDao
	err := dbcommons.GetDb().Where("app_key=? and group_id=? and member_id in (?)", appkey, groupId, memberIds).Find(&items).Error
	ret := map[string][]*GroupMemberExtDao{}
	if err != nil {
		return nil, err
	}
	for _, item := range items {
		if arr, exist := ret[item.MemberId]; exist {
			arr = append(arr, item)
			ret[item.MemberId] = arr
		} else {
			ret[item.MemberId] = []*GroupMemberExtDao{item}
		}
	}
	return ret, nil
}

func (ext GroupMemberExtDao) Upsert(appkey, groupId, memberId, itemKey, itemValue string, itemType int) error {
	return dbcommons.GetDb().Exec(fmt.Sprintf("INSERT INTO %s (app_key,group_id,member_id,item_key,item_value,item_type)VALUES(?,?,?,?,?,?) ON DUPLICATE KEY UPDATE item_value=?,updated_time=?", ext.TableName()), appkey, groupId, memberId, itemKey, itemValue, itemType, itemValue, time.Now()).Error
}

func (ext GroupMemberExtDao) BatchDelete(appkey, groupId string, memberIds []string) error {
	return dbcommons.GetDb().Where("app_key=? and group_id=? and member_id in (?)", appkey, groupId, memberIds).Delete(&GroupMemberExtDao{}).Error
}

func (ext GroupMemberExtDao) DeleteByGroupId(appkey, groupId string) error {
	return dbcommons.GetDb().Where("app_key=? and group_id=?", appkey, groupId).Delete(&GroupMemberExtDao{}).Error
}

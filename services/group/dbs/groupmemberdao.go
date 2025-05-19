package dbs

import (
	"bytes"
	"fmt"
	"im-server/commons/dbcommons"
	"time"
)

type GroupMemberDao struct {
	ID          int64     `gorm:"primary_key"`
	GroupId     string    `gorm:"group_id"`
	MemberId    string    `gorm:"member_id"`
	MemberType  int       `gorm:"member_type"`
	CreatedTime time.Time `gorm:"created_time"`
	AppKey      string    `gorm:"app_key"`
	IsMute      int       `gorm:"is_mute"`
	IsAllow     int       `gorm:"is_allow"`
	MuteEndAt   int64     `gorm:"mute_end_at"`
}

func (msg GroupMemberDao) TableName() string {
	return "groupmembers"
}

func (msg GroupMemberDao) Create(item GroupMemberDao) error {
	err := dbcommons.GetDb().Create(&item).Error
	return err
}

func (member GroupMemberDao) Find(appkey, groupId, memberId string) (*GroupMemberDao, error) {
	var item GroupMemberDao
	err := dbcommons.GetDb().Where("app_key=? and group_id=? and member_id=?", appkey, groupId, memberId).Take(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (member GroupMemberDao) BatchCreate(items []GroupMemberDao) error {
	var buffer bytes.Buffer
	sql := fmt.Sprintf("insert into %s (`group_id`,`member_id`,`app_key`)values", member.TableName())

	buffer.WriteString(sql)
	params := []interface{}{}
	for i, item := range items {
		if i == len(items)-1 {
			buffer.WriteString("(?,?,?);")
		} else {
			buffer.WriteString("(?,?,?),")
		}
		params = append(params, item.GroupId, item.MemberId, item.AppKey)
	}

	err := dbcommons.GetDb().Exec(buffer.String(), params...).Error
	return err
}

func (member GroupMemberDao) QueryMembers(appkey, groupId string, startId, limit int64) ([]*GroupMemberDao, error) {
	var items []*GroupMemberDao
	err := dbcommons.GetDb().Where("app_key=? and group_id=? and id>?", appkey, groupId, startId).Order("id asc").Limit(limit).Find(&items).Error
	return items, err
}

func (member GroupMemberDao) QueryGroupsByMemberId(appkey, memberId string, startId, limit int64) ([]*GroupMemberDao, error) {
	var items []*GroupMemberDao
	err := dbcommons.GetDb().Where("app_key=? and member_id=? and id>?", appkey, memberId, startId).Order("id asc").Limit(limit).Find(&items).Error
	return items, err
}

func (member GroupMemberDao) BatchDelete(appkey, groupId string, memberIds []string) error {
	return dbcommons.GetDb().Where("app_key=? and group_id=? and member_id in (?)", appkey, groupId, memberIds).Delete(&GroupMemberDao{}).Error
}

func (member GroupMemberDao) DeleteByGroupId(appkey, groupId string) error {
	return dbcommons.GetDb().Where("app_key=? and group_id=?", appkey, groupId).Delete(&GroupMemberDao{}).Error
}

func (member GroupMemberDao) UpdateMute(appkey, groupId string, isMute int, memberIds []string, muteEndAt int64) error {
	upd := map[string]interface{}{}
	upd["is_mute"] = isMute
	if isMute == 0 {
		upd["mute_end_at"] = 0
	} else {
		upd["mute_end_at"] = muteEndAt
	}
	return dbcommons.GetDb().Model(&GroupMemberDao{}).Where("app_key=? and group_id=? and member_id in (?)", appkey, groupId, memberIds).Update(upd).Error
}

func (member GroupMemberDao) UpdateAllow(appkey, groupId string, isAllow int, memberIds []string) error {
	return dbcommons.GetDb().Model(&member).Where("app_key=? and group_id=? and member_id in (?)", appkey, groupId, memberIds).Update("is_allow", isAllow).Error
}

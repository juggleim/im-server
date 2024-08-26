package dbs

import (
	"errors"
	"fmt"
	"im-server/commons/dbcommons"
	"time"

	"github.com/jinzhu/gorm"
)

type GroupDao struct {
	ID            int64     `gorm:"primary_key"`
	GroupId       string    `gorm:"group_id"`
	GroupName     string    `gorm:"group_name"`
	GroupPortrait string    `gorm:"group_portrait"`
	CreatedTime   time.Time `gorm:"created_time"`
	UpdatedTime   time.Time `gorm:"updated_time"`
	AppKey        string    `gorm:"app_key"`
	IsMute        int       `gorm:"is_mute"`
}

func (group GroupDao) TableName() string {
	return "groupinfos"
}
func (group GroupDao) Create(item GroupDao) error {
	err := dbcommons.GetDb().Create(&item).Error
	return err
}

func (group GroupDao) IsExist(appkey, groupId string) (bool, error) {
	var item GroupDao
	err := dbcommons.GetDb().Where("app_key=? and group_id=?", appkey, groupId).Take(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (group GroupDao) FindById(appkey, groupId string) (*GroupDao, error) {
	var item GroupDao
	err := dbcommons.GetDb().Where("app_key=? and group_id=?", appkey, groupId).Take(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (group GroupDao) Delete(appkey, groupId string) error {
	return dbcommons.GetDb().Where("app_key=? and group_id=?", appkey, groupId).Delete(&GroupMemberDao{}).Error
}

func (group GroupDao) UpdateGroupMuteStatus(appkey, groupId string, isMute int32) error {
	upd := map[string]interface{}{}
	upd["is_mute"] = isMute
	return dbcommons.GetDb().Model(&GroupDao{}).Where("app_key=? and group_id=?", appkey, groupId).Update(upd).Error
}

func (group GroupDao) UpdateGrpName(appkey, groupId, groupName, groupPortrait string) error {
	upd := map[string]interface{}{}
	if groupName != "" {
		upd["group_name"] = groupName
	}
	if groupPortrait != "" {
		upd["group_portrait"] = groupPortrait
	}
	if len(upd) > 0 {
		upd["updated_time"] = time.Now()
	} else {
		return fmt.Errorf("do nothing")
	}
	err := dbcommons.GetDb().Model(&GroupDao{}).Where("app_key=? and group_id=?", appkey, groupId).Update(upd).Error
	return err
}

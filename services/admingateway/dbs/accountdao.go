package dbs

import (
	"im-server/commons/dbcommons"
	"time"

	"github.com/jinzhu/gorm"
)

type AccountDao struct {
	ID            int64     `gorm:"primary_key"`
	Account       string    `gorm:"account"`
	Password      string    `gorm:"password"`
	CreatedTime   time.Time `gorm:"created_time"`
	UpdatedTime   time.Time `gorm:"updated_time"`
	State         int       `gorm:"state"` //0:normal; 1:forbidden
	RoleType      int       `gorm:"role_type"`
	ParentAccount string    `gorm:"parent_account"`
}

func (admin AccountDao) TableName() string {
	return "accounts"
}

func (admin AccountDao) Create(item AccountDao) error {
	err := dbcommons.GetDb().Create(&item).Error
	return err
}

func (admin AccountDao) FindByAccountPassword(account, password string) (*AccountDao, error) {
	var item AccountDao
	err := dbcommons.GetDb().Where("account=? and password=?", account, password).Take(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (admin AccountDao) FindByAccount(account string) (*AccountDao, error) {
	var item AccountDao
	err := dbcommons.GetDb().Where("account=?", account).Take(&item).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

func (admin AccountDao) UpdateState(accounts []string, state int) error {
	return dbcommons.GetDb().Model(&AccountDao{}).Where("account in (?)", accounts).Update("state", state).Error
}

func (admin AccountDao) UpdatePassword(account, password string) error {
	return dbcommons.GetDb().Model(&AccountDao{}).Where("account=?", account).Update("password", password).Error
}

func (admin AccountDao) QryAccounts(limit int64, offset int64) ([]*AccountDao, error) {
	var list []*AccountDao
	err := dbcommons.GetDb().Where("id > ?", offset).Order("id asc").Limit(limit).Find(&list).Error
	return list, err
}

func (admin AccountDao) DeleteAccounts(accounts []string) error {
	return dbcommons.GetDb().Where("account in (?)", accounts).Delete(&admin).Error
}

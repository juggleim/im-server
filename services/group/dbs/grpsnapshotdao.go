package dbs

import "im-server/commons/dbcommons"

type GrpSnapshotDao struct {
	ID          int64  `gorm:"primary_key"`
	AppKey      string `gorm:"app_key"`
	GroupId     string `gorm:"group_id"`
	CreatedTime int64  `gorm:"created_time"`
	Snapshot    []byte `gorm:"snapshot"`
}

func (snapshot GrpSnapshotDao) TableName() string {
	return "grpsnapshots"
}
func (snapshot GrpSnapshotDao) Create(item GrpSnapshotDao) error {
	err := dbcommons.GetDb().Create(&item).Error
	return err
}

func (snapshot GrpSnapshotDao) FindNearlySnapshot(appkey, groupId string, nearlyTime int64) (*GrpSnapshotDao, error) {
	var items []*GrpSnapshotDao
	err := dbcommons.GetDb().Where("app_key=? and group_id=? and created_time<?", appkey, groupId, nearlyTime).Order("created_time desc").Limit(1).Find(&items).Error
	if err == nil && len(items) > 0 {
		return items[0], nil
	}
	err = dbcommons.GetDb().Where("app_key=? and group_id=? and created_time>?", appkey, groupId, nearlyTime).Order("created_time asc").Limit(1).Find(&items).Error
	if err == nil && len(items) > 0 {
		return items[0], nil
	}
	return nil, err
}

func (snapshot GrpSnapshotDao) Exist(appkey, groupId string) bool {
	var item GrpSnapshotDao
	err := dbcommons.GetDb().Where("app_key=? and group_id=?", appkey, groupId).Select("id").Take(&item).Error
	if err == nil && item.ID > 0 {
		return true
	}
	return false
}

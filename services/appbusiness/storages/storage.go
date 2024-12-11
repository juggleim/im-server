package storages

import (
	"im-server/services/appbusiness/storages/dbs"
	"im-server/services/appbusiness/storages/models"
)

func NewFriendApplicationStorage() models.IFriendApplicationStorage {
	return &dbs.FriendApplicationDao{}
}

func NewGrpApplicationStorage() models.IGrpApplicationStorage {
	return &dbs.GrpApplicationDao{}
}

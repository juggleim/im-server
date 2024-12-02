package storages

import (
	"im-server/services/appbusiness/storages/dbs"
	"im-server/services/appbusiness/storages/models"
)

func NewFriendRelStorage() models.IFriendRelStorage {
	return &dbs.FriendRelDao{}
}

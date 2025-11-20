package storages

import (
	"im-server/services/friendmanager/storages/dbs"
	"im-server/services/friendmanager/storages/models"
)

func NewFriendRelStorage() models.IFriendRelStorage {
	return &dbs.FriendRelDao{}
}

package storages

import (
	"im-server/services/appbusiness/storages/models"
	"im-server/services/friends/storages/dbs"
)

func NewFriendRelStorage() models.IFriendRelStorage {
	return &dbs.FriendRelDao{}
}

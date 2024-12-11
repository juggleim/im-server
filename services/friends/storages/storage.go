package storages

import (
	"im-server/services/friends/storages/dbs"
	"im-server/services/friends/storages/models"
)

func NewFriendRelStorage() models.IFriendRelStorage {
	return &dbs.FriendRelDao{}
}

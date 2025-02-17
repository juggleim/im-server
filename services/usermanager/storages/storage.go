package storages

import (
	"im-server/services/usermanager/storages/dbs"
	"im-server/services/usermanager/storages/models"
)

func NewUserStorage() models.IUserStorage {
	return &dbs.UserDao{}
}

func NewUserExtStorage() models.IUserExtStorage {
	return &dbs.UserExtDao{}
}

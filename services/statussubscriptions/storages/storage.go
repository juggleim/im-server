package storages

import (
	"im-server/services/statussubscriptions/storages/dbs"
	"im-server/services/statussubscriptions/storages/models"
)

func NewUserSubRelStorage() models.IUserSubRelStorage {
	return &dbs.UserSubRelDao{}
}

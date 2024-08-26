package storages

import (
	"im-server/commons/configures"
	"im-server/commons/dbcommons"
	"im-server/commons/mongocommons"
	"im-server/services/pushmanager/storages/dbs"
	"im-server/services/pushmanager/storages/models"
	"im-server/services/pushmanager/storages/mongodbs"
)

func NewTagStorage() models.IUserTagStorage {
	switch configures.Config.MsgStoreEngine {
	case configures.MsgStoreEngine_MySQL:
		return dbs.NewUserTagsDao(dbcommons.GetDb())
	case configures.MsgStoreEngine_Mongo:
		return mongodbs.NewUserTagsDao(mongocommons.GetMongoDatabase())
	default:
		return dbs.NewUserTagsDao(dbcommons.GetDb())
	}
}

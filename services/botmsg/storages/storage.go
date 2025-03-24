package storages

import (
	"im-server/services/botmsg/storages/dbs"
	"im-server/services/botmsg/storages/models"
)

func NewBotConverStorage() models.IBotConverStorage {
	return &dbs.BotConverDao{}
}

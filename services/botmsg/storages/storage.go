package storages

import (
	"im-server/services/botmsg/storages/dbs"
	"im-server/services/botmsg/storages/models"
)

func NewBotConfStorage() models.IBotConfStorage {
	return &dbs.BotConfDao{}
}

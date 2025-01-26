package storages

import (
	"im-server/services/botmsg/storages/dbs"
	"im-server/services/botmsg/storages/models"
)

func NewBotConfStorage() models.IBotConfStorage {
	return &dbs.BotConfDao{}
}

func NewAssistantStorage() models.IAssistantStorage {
	return &dbs.AssistantDao{}
}

func NewBotConverStorage() models.IBotConverStorage {
	return &dbs.BotConverDao{}
}

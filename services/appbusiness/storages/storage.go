package storages

import (
	"im-server/services/appbusiness/storages/dbs"
	"im-server/services/appbusiness/storages/models"
)

func NewFriendApplicationStorage() models.IFriendApplicationStorage {
	return &dbs.FriendApplicationDao{}
}

func NewGrpApplicationStorage() models.IGrpApplicationStorage {
	return &dbs.GrpApplicationDao{}
}

func NewQrCodeRecordStorage() models.IQrCodeRecordStorage {
	return &dbs.QrCodeRecordDao{}
}

func NewSmsRecordStorage() models.ISmsRecordStorage {
	return &dbs.SmsRecordDao{}
}

func NewPromptStorage() models.IPromptStorage {
	return &dbs.PromptDao{}
}

func NewAiEngineStorage() models.IAiEngineStorage {
	return &dbs.AiEngineDao{}
}

func NewBotConfStorage() models.IBotConfStorage {
	return &dbs.BotConfDao{}
}

package models

import botModels "im-server/services/botmsg/storages/models"

type Assistant struct {
	ID          int64
	AssistantId string
	OwnerId     string
	Nickname    string
	Portrait    string
	Description string
	BotType     botModels.BotType
	BotConf     string
	Status      int
	AppKey      string
}

type IAssistantStorage interface {
	Create(item Assistant) error
	FindByAssistantId(appkey, assistantId string) (*Assistant, error)
	FindEnableAssistant(appkey string) (*Assistant, error)
}

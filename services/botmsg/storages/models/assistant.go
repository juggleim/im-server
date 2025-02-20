package models

type Assistant struct {
	ID          int64
	AssistantId string
	OwnerId     string
	Nickname    string
	Portrait    string
	Description string
	BotType     BotType
	BotConf     string
	Status      int
	AppKey      string
}

type IAssistantStorage interface {
	Create(item Assistant) error
	FindByAssistantId(appkey, assistantId string) (*Assistant, error)
}

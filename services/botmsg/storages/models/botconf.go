package models

type BotType int

var (
	BotType_Default     BotType = 0
	BotType_Custom      BotType = 1
	BotType_Dify        BotType = 2
	BotType_Coze        BotType = 3
	BotType_Minmax      BotType = 4
	BotType_SiliconFlow BotType = 5
)

type BotStatus int

var (
	BotStatus_Disable BotStatus = 0
	BotStatus_Enable  BotStatus = 1
)

type BotConf struct {
	ID          int64
	AppKey      string
	BotId       string
	Nickname    string
	BotPortrait string
	Description string
	BotType     BotType
	BotConf     string
	Status      BotStatus
}

type IBotConfStorage interface {
	Upsert(item BotConf) error
	FindById(appkey, botId string) (*BotConf, error)
	QryBotConfs(appkey string, startId, limit int64) ([]*BotConf, error)
	QryBotConfsWithStatus(appkey string, status BotStatus, startId, limit int64) ([]*BotConf, error)
}

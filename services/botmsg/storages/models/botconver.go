package models

type BotConverType int

var (
	BotConverType_Default BotConverType = 0
	BotConverType_Coze    BotConverType = 1
)

type BotConver struct {
	AppKey      string
	ConverType  BotConverType
	ConverKey   string
	ConverId    string
	UpdatedTime int64
}

type IBotConverStorage interface {
	Upsert(item BotConver) error
	Find(appkey string, converType BotConverType, converKey string) (*BotConver, error)
}

package models

type EngineType int

var (
	EngineType_SiliconFlow EngineType = 0
)

type AiEngine struct {
	ID         int64
	EngineType EngineType
	EngineConf string
	Status     int
	AppKey     string
}

type IAiEngineStorage interface {
	FindEnableAiEngine(appkey string) (*AiEngine, error)
}

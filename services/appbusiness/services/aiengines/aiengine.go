package aiengines

import "context"

type AssistantEngineType int

var (
	AssistantEngineType_SiliconFlow AssistantEngineType = 1
)

type IAiEngine interface {
	StreamChat(ctx context.Context, senderId, converId string, prompt string, question string, f func(answerPart string, isEnd bool))
}

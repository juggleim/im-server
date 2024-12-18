package botengines

import "context"

var DefaultBotEngine IBotEngine = &NilBotEngine{}

type IBotEngine interface {
	StreamChat(ctx context.Context, senderId, converId string, question string, f func(answerPart string, isEnd bool))
}

type NilBotEngine struct{}

func (engine *NilBotEngine) StreamChat(ctx context.Context, senderId, converId string, question string, f func(answerPart string, isEnd bool)) {
}

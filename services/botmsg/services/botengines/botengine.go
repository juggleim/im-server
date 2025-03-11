package botengines

import (
	"context"
	"im-server/commons/pbdefines/pbobjs"
)

var DefaultBotEngine IBotEngine = &NilBotEngine{}

type IBotEngine interface {
	StreamChat(ctx context.Context, senderId, converId string, channelType pbobjs.ChannelType, question string, f func(answerPart string, sectionStart, sectionEnd, isEnd bool))
	Chat(ctx context.Context, senderId, converId string, channelType pbobjs.ChannelType, question string) string
}

type NilBotEngine struct{}

func (engine *NilBotEngine) StreamChat(ctx context.Context, senderId, converKey string, channelType pbobjs.ChannelType, question string, f func(answerPart string, sectionStart, sectionEnd, isEnd bool)) {
}

func (engine *NilBotEngine) Chat(ctx context.Context, senderId, converKey string, channelType pbobjs.ChannelType, question string) string {
	return ""
}

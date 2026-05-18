package botengines

import (
	"context"
	"im-server/commons/pbdefines/pbobjs"
)

var DefaultBotEngine IBotEngine = &NilBotEngine{}

type IBotEngine interface {
	IsStreamChat() bool
	StreamChat(ctx context.Context, senderId, targetId string, msg *pbobjs.DownMsg, f func(answerPart string, sectionStart, sectionEnd, isEnd bool))
	Chat(ctx context.Context, senderId, targetId string, msg *pbobjs.DownMsg) string
}

type NilBotEngine struct{}

func (engine *NilBotEngine) IsStreamChat() bool {
	return false
}

func (engine *NilBotEngine) StreamChat(ctx context.Context, senderId, targetId string, msg *pbobjs.DownMsg, f func(answerPart string, sectionStart, sectionEnd, isEnd bool)) {
}

func (engine *NilBotEngine) Chat(ctx context.Context, senderId, targetId string, msg *pbobjs.DownMsg) string {
	return ""
}

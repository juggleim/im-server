package botclient

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/pbdefines/pbobjs"
)

func SendMsg2Bot(ctx context.Context, botId string, msg *pbobjs.DownMsg) {
	bases.AsyncRpcCall(ctx, "bot_msg", botId, msg)
}

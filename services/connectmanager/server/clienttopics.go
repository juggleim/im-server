package server

import (
	"time"

	"im-server/commons/errs"
	"im-server/services/connectmanager/server/codec"
	"im-server/services/connectmanager/server/imcontext"
	"im-server/services/connectmanager/services"
)

func allowClientPublish(ctx imcontext.WsHandleContext, msg *codec.PublishMsgBody) bool {
	if services.CheckClientTopic(msg.Topic) {
		return true
	}
	ctx.Write(codec.NewUserPublishAckMessage(&codec.PublishAckMsgBody{
		Index: msg.Index, Code: int32(errs.IMErrorCode_CONNECT_UNSUPPORTEDTOPIC), Timestamp: time.Now().UnixMilli(),
	}))
	return false
}

func allowClientQuery(ctx imcontext.WsHandleContext, msg *codec.QueryMsgBody) bool {
	if services.CheckClientQueryTopic(msg.Topic) {
		return true
	}
	ctx.Write(codec.NewQueryAckMessage(&codec.QueryAckMsgBody{
		Index: msg.Index, Code: int32(errs.IMErrorCode_CONNECT_UNSUPPORTEDTOPIC), Timestamp: time.Now().UnixMilli(),
	}, codec.QoS_NoAck))
	return false
}

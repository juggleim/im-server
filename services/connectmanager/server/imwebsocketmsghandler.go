package server

import (
	"im-server/services/connectmanager/server/codec"
	"im-server/services/connectmanager/server/imcontext"
)

type IMWebsocketMsgHandler struct {
	listener ImListener
}

func (handler IMWebsocketMsgHandler) HandleActive(ctx imcontext.WsHandleContext) {
	if handler.listener != nil {
		handler.listener.Create(ctx)
	}
}

func (handler IMWebsocketMsgHandler) HandleRead(ctx imcontext.WsHandleContext, message interface{}) {
	if handler.listener != nil {
		wsMsg, ok := message.(*codec.ImWebsocketMsg)
		if ok {
			switch wsMsg.Cmd {
			case int32(codec.Cmd_Connect):
				handler.listener.Connected(wsMsg.GetConnectMsgBody(), ctx)
			case int32(codec.Cmd_Disconnect):
				if !imcontext.CheckConnected(ctx) {
					ctx.Close(nil)
				}
				handler.listener.Diconnected(wsMsg.GetDisconnectMsgBody(), ctx)
			case int32(codec.Cmd_Ping):
				if !imcontext.CheckConnected(ctx) {
					ctx.Close(nil)
				}
				handler.listener.PingArrived(ctx)
			case int32(codec.Cmd_Publish):
				if !imcontext.CheckConnected(ctx) {
					ctx.Close(nil)
				}
				handler.listener.PublishArrived(wsMsg.GetPublishMsgBody(), int(wsMsg.GetQos()), ctx)
			case int32(codec.Cmd_PublishAck):
				if !imcontext.CheckConnected(ctx) {
					ctx.Close(nil)
				}
				handler.listener.PubAckArrived(wsMsg.GetPubAckMsgBody(), ctx)
			case int32(codec.Cmd_Query):
				if !imcontext.CheckConnected(ctx) {
					ctx.Close(nil)
				}
				handler.listener.QueryArrived(wsMsg.GetQryMsgBody(), ctx)
			case int32(codec.Cmd_QueryConfirm):
				if !imcontext.CheckConnected(ctx) {
					ctx.Close(nil)
				}
				handler.listener.QueryConfirmArrived(wsMsg.GetQryConfirmMsgBody(), ctx)
			default:
				break
			}
		}
	}
}

func (handler IMWebsocketMsgHandler) HandleInactive(ctx imcontext.WsHandleContext, ex error) {
	if handler.listener != nil {
		handler.listener.Close(ctx)
	}
	ctx.Close(ex)
}

func (handler IMWebsocketMsgHandler) HandleException(ctx imcontext.WsHandleContext, ex error) {
	if handler.listener != nil {
		handler.listener.ExceptionCaught(ctx, ex)
	}
	ctx.Close(ex)
}

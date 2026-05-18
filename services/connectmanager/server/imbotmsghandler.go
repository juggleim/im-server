package server

import (
	"fmt"
	"im-server/commons/errs"
	"im-server/services/connectmanager/server/codec"
	"im-server/services/connectmanager/server/imcontext"
)

type ImBotWebsocketMsgHandler struct {
	listener ImListener
}

func (handler ImBotWebsocketMsgHandler) HandleRead(ctx imcontext.WsHandleContext, message interface{}) {
	if handler.listener != nil {
		wsMsg, ok := message.(*codec.ImWebsocketMsg)
		if ok {
			switch wsMsg.Cmd {
			case int32(codec.Cmd_Connect):
				handler.listener.Connected(wsMsg.GetConnectMsgBody(), ctx)
			case int32(codec.Cmd_Disconnect):
				if !imcontext.CheckConnected(ctx) {
					ctx.Close(nil)
					return
				}
				disconnectMsg := wsMsg.GetDisconnectMsgBody()
				if disconnectMsg == nil {
					disconnectMsg = &codec.DisconnectMsgBody{}
				}
				handler.listener.Diconnected(disconnectMsg, ctx)
			case int32(codec.Cmd_Ping):
				if !imcontext.CheckConnected(ctx) {
					ctx.Close(nil)
					return
				}
				handler.listener.PingArrived(ctx)
			case int32(codec.Cmd_Publish):
				if !imcontext.CheckConnected(ctx) {
					ctx.Close(nil)
					return
				}
				handler.listener.PublishArrived(wsMsg.GetPublishMsgBody(), int(wsMsg.GetQos()), ctx)
			case int32(codec.Cmd_PublishAck):
				if !imcontext.CheckConnected(ctx) {
					ctx.Close(nil)
					return
				}
				handler.listener.PubAckArrived(wsMsg.GetPubAckMsgBody(), ctx)
			case int32(codec.Cmd_Query):
				if !imcontext.CheckConnected(ctx) {
					ctx.Close(nil)
					return
				}
				handler.listener.QueryArrived(wsMsg.GetQryMsgBody(), ctx)
			case int32(codec.Cmd_QueryConfirm):
				if !imcontext.CheckConnected(ctx) {
					ctx.Close(nil)
					return
				}
				handler.listener.QueryConfirmArrived(wsMsg.GetQryConfirmMsgBody(), ctx)
			default:
				ctx.Close(nil)
				if imcontext.CheckConnected(ctx) {
					handler.listener.ExceptionCaught(ctx, errs.IMErrorCode_CONNECT_CLOSE_DATA_ILLEGAL, fmt.Errorf("not support cmd:%d", wsMsg.Cmd))
				}
				return
			}
		}
	}
}

func (handler ImBotWebsocketMsgHandler) HandleException(ctx imcontext.WsHandleContext, code errs.IMErrorCode, ex error) {
	if handler.listener != nil {
		handler.listener.ExceptionCaught(ctx, code, ex)
	}
}

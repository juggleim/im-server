package codec

import (
	"im-server/commons/tools"
	"im-server/services/connectmanager/server/imcontext"
)

func (x *ImWebsocketMsg) Encrypt(ctx imcontext.WsHandleContext) {
	obfCode, exist := GetObfuscationCodeFromCtx(ctx)
	if exist {
		var payload []byte
		switch x.Cmd {
		case int32(Cmd_ConnectAck):
			payload, _ = tools.PbMarshal(x.GetConnectAckMsgBody())
		case int32(Cmd_Disconnect):
			payload, _ = tools.PbMarshal(x.GetDisconnectMsgBody())
		case int32(Cmd_Publish):
			payload, _ = tools.PbMarshal(x.GetPublishMsgBody())
		case int32(Cmd_PublishAck):
			payload, _ = tools.PbMarshal(x.GetPubAckMsgBody())
		case int32(Cmd_Query):
			payload, _ = tools.PbMarshal(x.GetQryMsgBody())
		case int32(Cmd_QueryAck):
			payload, _ = tools.PbMarshal(x.GetQryAckMsgBody())
		case int32(Cmd_QueryConfirm):
			payload, _ = tools.PbMarshal(x.GetQryConfirmMsgBody())
		}
		DoObfuscation(obfCode, payload)
		x.Payload = payload
		x.Testof = nil
	}
}

func (x *ImWebsocketMsg) Decrypt(ctx imcontext.WsHandleContext) {
	if x.Cmd == int32(Cmd_Connect) {
		if len(x.Payload) > 0 {
			obfCode := CalObfuscationCode(x.Payload)
			imcontext.SetContextAttr(ctx, imcontext.StateKey_ObfuscationCode, obfCode)
			DoObfuscation(obfCode, x.Payload)

			var msg ConnectMsgBody
			err := tools.PbUnMarshal(x.Payload, &msg)
			if err == nil {
				x.Testof = &ImWebsocketMsg_ConnectMsgBody{
					ConnectMsgBody: &msg,
				}
			}
		}
	} else {
		obfCode, exist := GetObfuscationCodeFromCtx(ctx)
		if exist {
			DoObfuscation(obfCode, x.Payload)
			switch x.Cmd {
			case int32(Cmd_ConnectAck):
				var msg ConnectAckMsgBody
				err := tools.PbUnMarshal(x.Payload, &msg)
				if err == nil {
					x.Testof = &ImWebsocketMsg_ConnectAckMsgBody{
						ConnectAckMsgBody: &msg,
					}
				}
			case int32(Cmd_Disconnect):
				var msg DisconnectMsgBody
				err := tools.PbUnMarshal(x.Payload, &msg)
				if err == nil {
					x.Testof = &ImWebsocketMsg_DisconnectMsgBody{
						DisconnectMsgBody: &msg,
					}
				}
			case int32(Cmd_Publish):
				var msg PublishMsgBody
				err := tools.PbUnMarshal(x.Payload, &msg)
				if err == nil {
					x.Testof = &ImWebsocketMsg_PublishMsgBody{
						PublishMsgBody: &msg,
					}
				}
			case int32(Cmd_PublishAck):
				var msg PublishAckMsgBody
				err := tools.PbUnMarshal(x.Payload, &msg)
				if err == nil {
					x.Testof = &ImWebsocketMsg_PubAckMsgBody{
						PubAckMsgBody: &msg,
					}
				}
			case int32(Cmd_Query):
				var msg QueryMsgBody
				err := tools.PbUnMarshal(x.Payload, &msg)
				if err == nil {
					x.Testof = &ImWebsocketMsg_QryMsgBody{
						QryMsgBody: &msg,
					}
				}
			case int32(Cmd_QueryAck):
				var msg QueryAckMsgBody
				err := tools.PbUnMarshal(x.Payload, &msg)
				if err == nil {
					x.Testof = &ImWebsocketMsg_QryAckMsgBody{
						QryAckMsgBody: &msg,
					}
				}
			case int32(Cmd_QueryConfirm):
				var msg QueryConfirmMsgBody
				err := tools.PbUnMarshal(x.Payload, &msg)
				if err == nil {
					x.Testof = &ImWebsocketMsg_QryConfirmMsgBody{
						QryConfirmMsgBody: &msg,
					}
				}
			}
		}
	}
}

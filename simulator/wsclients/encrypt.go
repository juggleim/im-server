package wsclients

import (
	"im-server/commons/tools"
	"im-server/services/connectmanager/server/codec"
)

func Encrypt(wsMsg *codec.ImWebsocketMsg, client *WsImClient) {
	if client.isEncrypt {
		var payload []byte
		switch wsMsg.Cmd {
		case int32(codec.Cmd_Connect):
			payload, _ = tools.PbMarshal(wsMsg.GetConnectMsgBody())
		case int32(codec.Cmd_Disconnect):
			payload, _ = tools.PbMarshal(wsMsg.GetDisconnectMsgBody())
		case int32(codec.Cmd_PublishAck):
			payload, _ = tools.PbMarshal(wsMsg.GetPubAckMsgBody())
		case int32(codec.Cmd_Publish):
			payload, _ = tools.PbMarshal(wsMsg.GetPublishMsgBody())
		case int32(codec.Cmd_Query):
			payload, _ = tools.PbMarshal(wsMsg.GetQryMsgBody())
		}

		codec.DoObfuscation(client.obfCode, payload)
		wsMsg.Payload = payload
		wsMsg.Testof = nil
	}
}

func Decrypt(wsImMsg *codec.ImWebsocketMsg, client *WsImClient) {
	if client.isEncrypt {
		var err error
		codec.DoObfuscation(client.obfCode, wsImMsg.Payload)
		switch wsImMsg.Cmd {
		case int32(codec.Cmd_ConnectAck):
			var msg codec.ConnectAckMsgBody
			err = tools.PbUnMarshal(wsImMsg.Payload, &msg)
			if err == nil {
				wsImMsg.Testof = &codec.ImWebsocketMsg_ConnectAckMsgBody{
					ConnectAckMsgBody: &msg,
				}
			}
		case int32(codec.Cmd_Disconnect):
			var msg codec.DisconnectMsgBody
			err = tools.PbUnMarshal(wsImMsg.Payload, &msg)
			if err == nil {
				wsImMsg.Testof = &codec.ImWebsocketMsg_DisconnectMsgBody{
					DisconnectMsgBody: &msg,
				}
			}
		case int32(codec.Cmd_Publish):
			var msg codec.PublishMsgBody
			err = tools.PbUnMarshal(wsImMsg.Payload, &msg)
			if err == nil {
				wsImMsg.Testof = &codec.ImWebsocketMsg_PublishMsgBody{
					PublishMsgBody: &msg,
				}
			}
		case int32(codec.Cmd_PublishAck):
			var msg codec.PublishAckMsgBody
			err = tools.PbUnMarshal(wsImMsg.Payload, &msg)
			if err == nil {
				wsImMsg.Testof = &codec.ImWebsocketMsg_PubAckMsgBody{
					PubAckMsgBody: &msg,
				}
			}
		case int32(codec.Cmd_QueryAck):
			var msg codec.QueryAckMsgBody
			err = tools.PbUnMarshal(wsImMsg.Payload, &msg)
			if err == nil {
				wsImMsg.Testof = &codec.ImWebsocketMsg_QryAckMsgBody{
					QryAckMsgBody: &msg,
				}
			}
		}
	}
}

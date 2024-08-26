package codec

import (
	"im-server/commons/tools"
)

type DisconnectMessage struct {
	MsgHeader
	MsgBody *DisconnectMsgBody
}

func NewDisconnectMessage(msgBody *DisconnectMsgBody) *DisconnectMessage {
	msg := &DisconnectMessage{
		MsgHeader: MsgHeader{
			Version: Version_1,
		},
		MsgBody: msgBody,
	}
	msg.SetCmd(Cmd_Disconnect)
	msg.SetQoS(QoS_NoAck)
	return msg
}
func NewDisconnectMessageWithHeader(header *MsgHeader) *DisconnectMessage {
	msg := &DisconnectMessage{
		MsgHeader: MsgHeader{
			Version:     Version_1,
			HeaderCode:  header.HeaderCode,
			Checksum:    header.Checksum,
			MsgBodySize: header.MsgBodySize,
		},
	}
	msg.SetCmd(Cmd_Disconnect)
	msg.SetQoS(QoS_NoAck)
	return msg
}

func (msg *DisconnectMessage) EncodeBody() ([]byte, error) {
	if msg.MsgBody != nil {
		return tools.PbMarshal(msg.MsgBody)
	}
	return nil, &CodecError{"MsgBody's length is 0."}
}

func (msg *DisconnectMessage) DecodeBody(msgBodyBytes []byte) error {
	msg.MsgBody = &DisconnectMsgBody{}
	return tools.PbUnMarshal(msgBodyBytes, msg.MsgBody)
}

func (msg *DisconnectMessage) ToImWebsocketMsg() *ImWebsocketMsg {
	return &ImWebsocketMsg{
		Version: int32(msg.MsgHeader.Version),
		Cmd:     int32(msg.GetCmd()),
		Qos:     int32(msg.GetQoS()),
		Testof: &ImWebsocketMsg_DisconnectMsgBody{
			DisconnectMsgBody: msg.MsgBody,
		},
	}
}

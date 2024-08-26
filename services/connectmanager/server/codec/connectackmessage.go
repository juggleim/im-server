package codec

import (
	"im-server/commons/tools"
)

type ConnectAckMessage struct {
	MsgHeader
	MsgBody *ConnectAckMsgBody
}

func NewConnectAckMessage(msgBody *ConnectAckMsgBody) *ConnectAckMessage {
	msg := &ConnectAckMessage{
		MsgHeader: MsgHeader{
			Version: Version_1,
		},
		MsgBody: msgBody,
	}
	msg.SetQoS(QoS_NoAck)
	msg.SetCmd(Cmd_ConnectAck)
	return msg
}
func NewConnectAckMessageWithHeader(header *MsgHeader) *ConnectAckMessage {
	msg := &ConnectAckMessage{
		MsgHeader: MsgHeader{
			Version:     Version_1,
			HeaderCode:  header.HeaderCode,
			Checksum:    header.Checksum,
			MsgBodySize: header.MsgBodySize,
		},
	}
	msg.SetCmd(Cmd_ConnectAck)
	msg.SetQoS(QoS_NoAck)
	return msg
}

func (msg *ConnectAckMessage) EncodeBody() ([]byte, error) {
	if msg.MsgBody != nil {
		return tools.PbMarshal(msg.MsgBody)
	}
	return nil, &CodecError{"MsgBody's length is 0."}
}

func (msg *ConnectAckMessage) DecodeBody(msgBodyBytes []byte) error {
	msg.MsgBody = &ConnectAckMsgBody{}
	return tools.PbUnMarshal(msgBodyBytes, msg.MsgBody)
}
func (msg *ConnectAckMessage) ToImWebsocketMsg() *ImWebsocketMsg {
	return &ImWebsocketMsg{
		Version: int32(msg.MsgHeader.Version),
		Cmd:     int32(msg.GetCmd()),
		Qos:     int32(msg.GetQoS()),
		Testof: &ImWebsocketMsg_ConnectAckMsgBody{
			ConnectAckMsgBody: msg.MsgBody,
		},
	}
}

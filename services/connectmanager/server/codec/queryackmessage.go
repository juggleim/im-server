package codec

import (
	"im-server/commons/tools"
)

type QueryAckMessage struct {
	MsgHeader
	MsgBody *QueryAckMsgBody
}

func NewQueryAckMessage(msgBody *QueryAckMsgBody, qos int) *QueryAckMessage {
	msg := &QueryAckMessage{
		MsgHeader: MsgHeader{
			Version: Version_1,
		},
		MsgBody: msgBody,
	}
	msg.SetCmd(Cmd_QueryAck)
	msg.SetQoS(qos)
	return msg
}
func NewQueryAckMessageWithHeader(header *MsgHeader) *QueryAckMessage {
	msg := &QueryAckMessage{
		MsgHeader: MsgHeader{
			Version:     Version_1,
			HeaderCode:  header.HeaderCode,
			Checksum:    header.Checksum,
			MsgBodySize: header.MsgBodySize,
		},
	}
	msg.SetCmd(Cmd_QueryAck)
	msg.SetQoS(QoS_NoAck)
	return msg
}

func (msg *QueryAckMessage) EncodeBody() ([]byte, error) {
	if msg.MsgBody != nil {
		return tools.PbMarshal(msg.MsgBody)
	}
	return nil, &CodecError{"MsgBody's length is 0."}
}

func (msg *QueryAckMessage) DecodeBody(msgBodyBytes []byte) error {
	msg.MsgBody = &QueryAckMsgBody{}
	return tools.PbUnMarshal(msgBodyBytes, msg.MsgBody)
}

func (msg *QueryAckMessage) ToImWebsocketMsg() *ImWebsocketMsg {
	return &ImWebsocketMsg{
		Version: int32(msg.MsgHeader.Version),
		Cmd:     int32(msg.GetCmd()),
		Qos:     int32(msg.GetQoS()),
		Testof: &ImWebsocketMsg_QryAckMsgBody{
			QryAckMsgBody: msg.MsgBody,
		},
	}
}

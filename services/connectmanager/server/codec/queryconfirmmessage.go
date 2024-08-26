package codec

import (
	"im-server/commons/tools"
)

type QueryConfirmMessage struct {
	MsgHeader
	MsgBody *QueryConfirmMsgBody
}

func NewWsQueryConfirmMessage(msgBody *QueryConfirmMsgBody) *ImWebsocketMsg {
	msg := &ImWebsocketMsg{
		Version: int32(Version_1),
		Cmd:     int32(Cmd_QueryConfirm),
		Qos:     int32(QoS_NoAck),
		Testof: &ImWebsocketMsg_QryConfirmMsgBody{
			QryConfirmMsgBody: msgBody,
		},
	}
	return msg
}
func NewQueryConfirmMessageWithHeader(header *MsgHeader) *QueryConfirmMessage {
	msg := &QueryConfirmMessage{
		MsgHeader: MsgHeader{
			Version:     Version_1,
			HeaderCode:  header.HeaderCode,
			Checksum:    header.Checksum,
			MsgBodySize: header.MsgBodySize,
		},
	}
	return msg
}

func (msg *QueryConfirmMessage) EncodeBody() ([]byte, error) {
	if msg.MsgBody != nil {
		return tools.PbMarshal(msg.MsgBody)
	}
	return nil, &CodecError{"MsgBody's length is 0."}
}

func (msg *QueryConfirmMessage) DecodeBody(msgBodyBytes []byte) error {
	msg.MsgBody = &QueryConfirmMsgBody{}
	return tools.PbUnMarshal(msgBodyBytes, msg.MsgBody)
}

func (msg *QueryConfirmMessage) ToImWebsocketMsg() *ImWebsocketMsg {
	return &ImWebsocketMsg{
		Version: int32(msg.MsgHeader.Version),
		Cmd:     int32(msg.GetCmd()),
		Qos:     int32(msg.GetQoS()),
		Testof: &ImWebsocketMsg_QryConfirmMsgBody{
			QryConfirmMsgBody: msg.MsgBody,
		},
	}
}

package codec

import (
	"im-server/commons/tools"
)

type UserPublishMessage struct {
	MsgHeader
	MsgBody *PublishMsgBody
}
type ServerPublishMessage struct {
	MsgHeader
	MsgBody *PublishMsgBody
}

func NewWsServerPublishMessage(msgBody *PublishMsgBody, qos int32) *ImWebsocketMsg {
	msg := &ImWebsocketMsg{
		Version: int32(Version_1),
		Cmd:     int32(Cmd_Publish),
		Qos:     qos,
		Testof: &ImWebsocketMsg_PublishMsgBody{
			PublishMsgBody: msgBody,
		},
	}
	return msg
}

func NewServerPublishMessage(msgBody *PublishMsgBody, qos int) *ServerPublishMessage {
	msg := &ServerPublishMessage{
		MsgHeader: MsgHeader{
			Version: Version_1,
		},
		MsgBody: msgBody,
	}
	msg.SetCmd(Cmd_Publish)
	msg.SetQoS(qos)
	return msg
}
func NewServerPublishMessageWithHeader(header *MsgHeader) *ServerPublishMessage {
	msg := &ServerPublishMessage{
		MsgHeader: MsgHeader{
			Version:     Version_1,
			HeaderCode:  header.HeaderCode,
			Checksum:    header.Checksum,
			MsgBodySize: header.MsgBodySize,
		},
	}
	return msg
}

func (msg *ServerPublishMessage) EncodeBody() ([]byte, error) {
	if msg.MsgBody != nil {
		return tools.PbMarshal(msg.MsgBody)
	}
	return nil, &CodecError{"MsgBody's length is 0."}
}

func (msg *ServerPublishMessage) DecodeBody(msgBodyBytes []byte) error {
	msg.MsgBody = &PublishMsgBody{}
	return tools.PbUnMarshal(msgBodyBytes, msg.MsgBody)
}
func (msg *ServerPublishMessage) ToImWebsocketMsg() *ImWebsocketMsg {
	return &ImWebsocketMsg{
		Version: int32(msg.MsgHeader.Version),
		Cmd:     int32(msg.GetCmd()),
		Qos:     int32(msg.GetQoS()),
		Testof: &ImWebsocketMsg_PublishMsgBody{
			PublishMsgBody: msg.MsgBody,
		},
	}
}

func NewUserPublishMessage(msgBody *PublishMsgBody) *UserPublishMessage {
	msg := &UserPublishMessage{
		MsgHeader: MsgHeader{
			Version: Version_1,
		},
		MsgBody: msgBody,
	}
	msg.SetCmd(Cmd_Publish)
	msg.SetQoS(QoS_NeedAck)
	return msg
}
func NewUserPublishMessageWithHeader(header *MsgHeader) *UserPublishMessage {
	msg := &UserPublishMessage{
		MsgHeader: MsgHeader{
			Version:     Version_1,
			HeaderCode:  header.HeaderCode,
			Checksum:    header.Checksum,
			MsgBodySize: header.MsgBodySize,
		},
	}
	return msg
}

func (msg *UserPublishMessage) EncodeBody() ([]byte, error) {
	if msg.MsgBody != nil {
		return tools.PbMarshal(msg.MsgBody)
	}
	return nil, &CodecError{"MsgBody's length is 0."}
}

func (msg *UserPublishMessage) DecodeBody(msgBodyBytes []byte) error {
	msg.MsgBody = &PublishMsgBody{}
	return tools.PbUnMarshal(msgBodyBytes, msg.MsgBody)
}

func (msg *UserPublishMessage) ToImWebsocketMsg() *ImWebsocketMsg {
	return &ImWebsocketMsg{
		Version: int32(msg.MsgHeader.Version),
		Cmd:     int32(msg.GetCmd()),
		Qos:     int32(msg.GetQoS()),
		Testof: &ImWebsocketMsg_PublishMsgBody{
			PublishMsgBody: msg.MsgBody,
		},
	}
}

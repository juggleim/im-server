package codec

type PingMessage struct {
	MsgHeader
}

func NewPingMessage() *PingMessage {
	msg := &PingMessage{
		MsgHeader: MsgHeader{
			Version: Version_1,
		},
	}
	msg.SetCmd(Cmd_Ping)
	msg.SetQoS(QoS_NeedAck)
	return msg
}
func NewPingMessageWithHeader(header *MsgHeader) *PingMessage {
	msg := &PingMessage{
		MsgHeader: MsgHeader{
			Version:     Version_1,
			HeaderCode:  header.HeaderCode,
			Checksum:    header.Checksum,
			MsgBodySize: header.MsgBodySize,
		},
	}
	msg.SetCmd(Cmd_Ping)
	msg.SetQoS(QoS_NeedAck)
	return msg
}
func (msg *PingMessage) EncodeBody() ([]byte, error) {
	return []byte{}, nil
}

func (msg *PingMessage) DecodeBody(msgBodyBytes []byte) error {
	return nil
}
func (msg *PingMessage) ToImWebsocketMsg() *ImWebsocketMsg {
	return &ImWebsocketMsg{
		Version: int32(msg.MsgHeader.Version),
		Cmd:     int32(msg.GetCmd()),
		Qos:     int32(msg.GetQoS()),
	}
}

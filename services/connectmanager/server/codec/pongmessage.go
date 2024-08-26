package codec

type PongMessage struct {
	MsgHeader
}

func NewPongMessageWithHeader(header *MsgHeader) *PongMessage {
	msg := &PongMessage{
		MsgHeader: MsgHeader{
			Version:     Version_1,
			HeaderCode:  header.HeaderCode,
			Checksum:    header.Checksum,
			MsgBodySize: header.MsgBodySize,
		},
	}
	msg.SetCmd(Cmd_Pong)
	msg.SetQoS(QoS_NoAck)
	return msg
}

func NewPongMessage() *PongMessage {
	msg := &PongMessage{
		MsgHeader: MsgHeader{
			Version: Version_1,
		},
	}
	msg.SetCmd(Cmd_Pong)
	msg.SetQoS(QoS_NoAck)
	return msg
}

func (msg *PongMessage) EncodeBody() ([]byte, error) {
	return []byte{}, nil
}

func (msg *PongMessage) DecodeBody(msgBodyBytes []byte) error {
	return nil
}

func (msg *PongMessage) ToImWebsocketMsg() *ImWebsocketMsg {
	return &ImWebsocketMsg{
		Version: int32(msg.MsgHeader.Version),
		Cmd:     int32(msg.GetCmd()),
		Qos:     int32(msg.GetQoS()),
	}
}

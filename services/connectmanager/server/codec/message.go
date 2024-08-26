package codec

import (
	"bytes"
	"io"
)

const (
	Version_1 byte   = byte(1)
	ProtoId   string = "jug9le1m"

	QoS_NoAck   int = 0
	QoS_NeedAck int = 1

	Cmd_Connect      int = 0
	Cmd_ConnectAck   int = 1
	Cmd_Disconnect   int = 2
	Cmd_Publish      int = 3
	Cmd_PublishAck   int = 4
	Cmd_Query        int = 5
	Cmd_QueryAck     int = 6
	Cmd_QueryConfirm int = 7
	Cmd_Ping         int = 8
	Cmd_Pong         int = 9
)

type CodecError struct {
	ErrMsg string
}

func (err *CodecError) Error() string {
	return err.ErrMsg
}

type IMessage interface {
	GetQoS() int
	SetQoS(int)
	GetCmd() int
	SetCmd(int)
	EncodeHeader(buf *bytes.Buffer, bodyBytes []byte)
	DecodeHeader(reader io.Reader)
	EncodeBody() ([]byte, error)
	DecodeBody(msgBodyBytes []byte) error

	ToImWebsocketMsg() *ImWebsocketMsg
}

type MsgHeader struct {
	Version     byte
	HeaderCode  byte
	Checksum    byte
	MsgBodySize int
}

// type BaseMessage struct {
// 	MsgHeader
// 	Name string
// 	Age  int
// }

func (msg *MsgHeader) GetQoS() int {
	var qos byte
	qos = msg.HeaderCode >> 2
	qos = qos & 0x03
	return int(qos)
}
func (msg *MsgHeader) SetQoS(qos int) {
	tmpQoS := byte(qos)
	tmpQoS = tmpQoS & 0x03
	tmpQoS = tmpQoS << 2
	msg.HeaderCode = msg.HeaderCode | tmpQoS
}

func (msg *MsgHeader) GetCmd() int {
	var cmd byte
	cmd = msg.HeaderCode >> 4
	cmd = cmd & 0x0f
	return int(cmd)
}
func (msg *MsgHeader) SetCmd(cmd int) {
	tmpCmd := byte(cmd)
	tmpCmd = tmpCmd & 0x0f
	tmpCmd = tmpCmd << 4
	msg.HeaderCode = msg.HeaderCode | tmpCmd
}

func (msg *MsgHeader) EncodeHeader(buf *bytes.Buffer, bodyBytes []byte) {
	buf.WriteByte(msg.Version)
	buf.WriteByte(msg.HeaderCode)

	msg.MsgBodySize = len(bodyBytes)
	msg.Checksum = calChecksum(msg.HeaderCode, bodyBytes)
	buf.WriteByte(msg.Checksum)

	if msg.GetCmd() != Cmd_Ping && msg.GetCmd() != Cmd_Pong {
		//write body size
		buf.Write(MsgBodySize2Bytes(msg.MsgBodySize))
	}
}

func (msg *MsgHeader) DecodeHeader(reader io.Reader) {
	tmpBs := make([]byte, 1)
	reader.Read(tmpBs)
	msg.HeaderCode = tmpBs[0]

	reader.Read(tmpBs)
	msg.Checksum = tmpBs[0]

	if msg.GetCmd() != Cmd_Ping && msg.GetCmd() != Cmd_Pong {
		msg.MsgBodySize = Bytes2MsgBodySize(reader)
	}
}
func calChecksum(initChecksum byte, bodyBytes []byte) byte {
	checksum := initChecksum
	for _, b := range bodyBytes {
		checksum = checksum ^ b
	}
	return checksum
}
func (msg *MsgHeader) ValidateChecksum(bodyBytes []byte) bool {
	checksum := calChecksum(msg.HeaderCode, bodyBytes)
	if checksum == msg.Checksum {
		return true
	} else {
		return false
	}
}

func MsgBodySize2Bytes(msgBodySize int) []byte {
	buf := bytes.NewBuffer([]byte{})
	msgLength := msgBodySize
	if msgLength > 0 {
		for {
			b := byte(msgLength & 0x7f)
			msgLength = msgLength >> 7
			if msgLength > 0 {
				b = b | 0x80
			}
			buf.WriteByte(b)
			if msgLength <= 0 {
				break
			}
		}
	} else {
		buf.WriteByte(0)
	}
	return buf.Bytes()
}
func Bytes2MsgBodySize(reader io.Reader) int {
	msgLength := 0
	multiplier := 1
	digit := 0
	for {
		readByte, err := readByte(reader)
		if err != nil {
			return 0
		}
		digit = int(readByte)
		msgLength += (digit & 0x7f) * multiplier
		multiplier *= 128
		if (digit & 0x80) <= 0 {
			break
		}
	}
	return msgLength
}

func readByte(reader io.Reader) (byte, error) {
	bs := make([]byte, 1)
	_, err := reader.Read(bs)
	return bs[0], err
}

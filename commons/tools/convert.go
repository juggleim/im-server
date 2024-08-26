package tools

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"strconv"

	"google.golang.org/protobuf/proto"
)

func Bool2String(b bool) string {
	if b {
		return "1"
	} else {
		return "0"
	}
}

func String2Bool(str string) bool {
	b, err := strconv.ParseBool(str)
	if err == nil {
		return b
	}
	return false
}

func String2Bytes(s string) []byte {
	// reader := strings.NewReader(s)
	// bytes := make([]byte, reader.Size())
	// reader.ReadAt(bytes, 0)
	bytes := []byte(s)
	return bytes
}

func Bytes2String(bytes []byte) string {
	// sb := strings.Builder{}
	// sb.Write(bytes)
	// return sb.String()
	return string(bytes)
}

func Bytes2Int(b []byte) int {
	buf := bytes.NewBuffer(b)
	var tmp uint32
	binary.Read(buf, binary.BigEndian, &tmp)
	return int(tmp)
}

func Int2Bytes(i int) []byte {
	buf := bytes.NewBuffer([]byte{})
	tmp := uint32(i)
	binary.Write(buf, binary.BigEndian, tmp)
	return buf.Bytes()
}

func String2Int64(str string) (int64, error) {
	return strconv.ParseInt(str, 10, 64)
}

func Int642String(i int64) string {
	return strconv.FormatInt(i, 64)
}

func Int64ToBytes(i int64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	return buf
}

func BytesToInt64(buf []byte) int64 {
	return int64(binary.BigEndian.Uint64(buf))
}

func UInt64ToBytes(i uint64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, i)
	return buf
}

func BytesToUInt64(buf []byte) uint64 {
	return binary.BigEndian.Uint64(buf)
}

func PbMarshal(obj proto.Message) ([]byte, error) {
	bytes, err := proto.Marshal(obj)
	return bytes, err
}
func PbUnMarshal(bytes []byte, typeScope proto.Message) error {
	err := proto.Unmarshal(bytes, typeScope)
	return err
}

func JsonMarshal(obj interface{}) ([]byte, error) {
	return json.Marshal(obj)
}

func JsonUnMarshal(bytes []byte, obj interface{}) error {
	return json.Unmarshal(bytes, obj)
}

func IntPtr(v int) *int {
	return &v
}
func Int32Ptr(v int32) *int32 {
	return &v
}
func Int64Ptr(systemID int64) *int64 {
	return &systemID
}
func BoolPtr(v bool) *bool {
	return &v
}
func StringPtr(v string) *string {
	return &v
}

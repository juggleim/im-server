package utils

import (
	"encoding/binary"
	"encoding/json"
	"strings"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

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
func GenerateUUID() uuid.UUID {
	uid := uuid.New()
	return uid
}
func UUID2Bytes(uuid uuid.UUID) []byte {
	bytes := make([]byte, 16)
	for i := 0; i < 16; i++ {
		bytes[i] = uuid[i]
	}
	return bytes
}
func GenerateUUIDBytes() []byte {
	return UUID2Bytes(GenerateUUID())
}

func GenerateUUIDShortString() string {
	return UUID2ShortString(GenerateUUID())
}

func Bytes2ShortString(uuid []byte) string {
	var bs [16]byte
	for i := 0; i < 16; i++ {
		bs[i] = uuid[i]
	}
	return UUIDBytes2ShortString(bs)
}

func UUID2ShortString(uuid uuid.UUID) string {
	return UUIDBytes2ShortString(uuid)
}

func UUIDBytes2ShortString(uuid [16]byte) string {
	mostBits := make([]byte, 8)
	leastBits := make([]byte, 8)
	for i := 0; i < 8; i++ {
		mostBits[i] = uuid[i]
	}
	for i := 8; i < 16; i++ {
		leastBits[i-8] = uuid[i]
	}
	return strings.Join([]string{toIdString(BytesToUInt64(mostBits)), toIdString(BytesToUInt64(leastBits))}, "")
}

var DIGITS64 []byte = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-_")

func toIdString(l uint64) string {
	buf := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	var length int = 11
	var least uint64 = 63 //0x3f

	for {
		length--
		buf[length] = DIGITS64[int(l&least)]
		l = l >> 6
		if l == 0 {
			break
		}
	}
	return string(buf)
}

/**
 *
 *
 *
**/
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

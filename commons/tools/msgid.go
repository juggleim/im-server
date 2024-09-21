package tools

import (
	"bytes"
	"encoding/binary"
	"hash/crc32"
	"im-server/commons/pbdefines/pbobjs"
	"sync/atomic"
)

const (
	base32EncodeChars string = "abcdefghjklmnpqrstuvwxyz23456789"
)

var (
	base32DecodeChars map[string]int64 = map[string]int64{
		"a": 0,
		"b": 1,
		"c": 2,
		"d": 3,
		"e": 4,
		"f": 5,
		"g": 6,
		"h": 7,
		"j": 8,
		"k": 9,
		"l": 10,
		"m": 11,
		"n": 12,
		"p": 13,
		"q": 14,
		"r": 15,
		"s": 16,
		"t": 17,
		"u": 18,
		"v": 19,
		"w": 20,
		"x": 21,
		"y": 22,
		"z": 23,
		"2": 24,
		"3": 25,
		"4": 26,
		"5": 27,
		"6": 28,
		"7": 29,
		"8": 30,
		"9": 31,
	}
)

var currentSeq uint32 = 0

func GenerateMsgId(time int64, channelType int32, targetId string) string {
	seq := getSeq()
	time = time << 12
	time = time | int64(seq)

	time = time << 4
	time = time | (int64(channelType) & 0xf)

	targetHashCode := crc32.ChecksumIEEE([]byte(targetId))
	targetHashCode = targetHashCode & 0x3fffff
	time = time << 6
	time = time | (int64(targetHashCode >> 16))
	low := targetHashCode << 16

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, time)
	binary.Write(buf, binary.BigEndian, low)

	bs := buf.Bytes()
	b1 := bs[0]
	b2 := bs[1]
	b3 := bs[2]
	b4 := bs[3]
	b5 := bs[4]
	b6 := bs[5]
	b7 := bs[6]
	b8 := bs[7]
	b9 := bs[8]
	b10 := bs[9]

	retBs := []byte{}

	retBs = append(retBs, base32EncodeChars[b1>>3])
	retBs = append(retBs, base32EncodeChars[((b1&0x7)<<2)|(b2>>6)])
	retBs = append(retBs, base32EncodeChars[(b2&0x3e)>>1])
	retBs = append(retBs, base32EncodeChars[((b2&0x1)<<4)|(b3>>4)])
	//retBs = append(retBs, '-')
	retBs = append(retBs, base32EncodeChars[((b3&0xf)<<1)|(b4>>7)])
	retBs = append(retBs, base32EncodeChars[(b4&0x7c)>>2])
	retBs = append(retBs, base32EncodeChars[((b4&0x3)<<3)|(b5>>5)])
	retBs = append(retBs, base32EncodeChars[b5&0x1f])
	//retBs = append(retBs, '-')
	retBs = append(retBs, base32EncodeChars[b6>>3])
	retBs = append(retBs, base32EncodeChars[((b6&0x7)<<2)|(b7>>6)])
	retBs = append(retBs, base32EncodeChars[(b7&0x3e)>>1])
	retBs = append(retBs, base32EncodeChars[((b7&0x1)<<4)|(b8>>4)])
	//retBs = append(retBs, '-')
	retBs = append(retBs, base32EncodeChars[((b8&0xf)<<1)|(b9>>7)])
	retBs = append(retBs, base32EncodeChars[(b9&0x7c)>>2])
	retBs = append(retBs, base32EncodeChars[((b9&0x3)<<3)|(b10>>5)])
	retBs = append(retBs, base32EncodeChars[b10&0x1f])

	return string(retBs)
}

func getSeq() uint32 {
	seq := atomic.AddUint32(&currentSeq, 1)
	seq = seq & 0xfff
	return seq
}

func ParseTimeFromMsgId(msgId string) int64 {
	if len(msgId) <= 9 {
		return 0
	}
	var t int64 = 0
	for i := 0; i < 8; i++ {
		c := msgId[i : i+1]
		intVal := base32DecodeChars[c]
		intVal = intVal & 0b11111
		t = t << 5
		t = t | intVal
	}
	c := msgId[8:9]
	intVal := base32DecodeChars[c]
	intVal = intVal >> 3
	intVal = intVal & 0b11
	t = t << 2
	t = t | intVal
	return t
}

func ParseChannelTypeFromMsgId(msgId string) pbobjs.ChannelType {
	if len(msgId) < 12 {
		return pbobjs.ChannelType_Unknown
	}
	var channelType int64 = 0
	a := msgId[10:11]
	intVal := base32DecodeChars[a]
	intVal = intVal & 0b1
	channelType = intVal
	channelType = channelType << 3
	b := msgId[11:12]
	intVal = base32DecodeChars[b]
	intVal = intVal >> 2
	channelType = channelType | intVal
	return pbobjs.ChannelType(channelType)
}

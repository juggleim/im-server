package commonservices

import (
	"im-server/commons/caches"
	"time"
)

var msgClientIdCache *caches.LruCache

func init() {
	msgClientIdCache = caches.NewLruCacheWithReadTimeout("msgclientid_cache", 10000, nil, 10*time.Minute)
}

type MsgAck struct {
	MsgId   string
	MsgTime int64
	MsgSeq  int64
}

func FilterDuplicateMsg(clientMsgId string, ack MsgAck) (*MsgAck, bool) {
	if old, succ := msgClientIdCache.AddIfAbsent(clientMsgId, &MsgAck{
		MsgId:   ack.MsgId,
		MsgTime: ack.MsgTime,
		MsgSeq:  ack.MsgSeq,
	}); succ {
		return nil, false
	} else {
		return old.(*MsgAck), true
	}
}

func RecordMsg(clientMsgId string, ack MsgAck) {
	msgClientIdCache.AddIfAbsent(clientMsgId, &MsgAck{
		MsgId:   ack.MsgId,
		MsgTime: ack.MsgTime,
		MsgSeq:  ack.MsgSeq,
	})
}

package services

import (
	"errors"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"testing"
)

func TestDeleteBuildersPreserveTimeAndSubChannel(t *testing.T) {
	const timestamp int64 = 1790000000123
	msg := &pbobjs.SimpleMsg{MsgId: tools.GenerateMsgId(timestamp, 1, "target")}
	if err := normalizeDelMsgTimes([]*pbobjs.SimpleMsg{nil, msg}); err != nil {
		t.Fatal(err)
	}
	p := newPrivateDelHisMsg("user", "target", "sub", "app", msg)
	g := newGroupDelHisMsg("user", "target", "sub", "app", msg)
	if p.MsgTime != timestamp || g.MsgTime != timestamp || p.SubChannel != "sub" || g.SubChannel != "sub" {
		t.Fatalf("p=%+v g=%+v", p, g)
	}
	msg.MsgTime = timestamp + 1
	if got, err := resolveDelMsgTime(msg); err != nil || got != timestamp+1 {
		t.Fatalf("explicit time lost: %d %v", got, err)
	}
}

func TestInvalidDeleteIDFailsNormalization(t *testing.T) {
	if err := normalizeDelMsgTimes([]*pbobjs.SimpleMsg{{MsgId: "bad"}}); !errors.Is(err, tools.ErrInvalidMsgId) {
		t.Fatalf("err=%v", err)
	}
}

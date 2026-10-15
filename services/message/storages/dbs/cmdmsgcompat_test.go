package dbs

import (
	"testing"

	"im-server/services/message/storages/models"
)

func TestMergeCmdCompatMsgsOrdersAndLimits(t *testing.T) {
	legacy := []*models.Msg{
		{MsgId: "old-1", SendTime: 100},
		{MsgId: "old-3", SendTime: 300},
		{MsgId: "old-4", SendTime: 400},
	}
	split := []*models.Msg{
		{MsgId: "new-2", SendTime: 200},
		{MsgId: "new-3", SendTime: 300},
		{MsgId: "new-5", SendTime: 500},
	}

	got := mergeCmdCompatMsgs(legacy, split, 4)
	want := []string{"old-1", "new-2", "old-3", "new-3"}
	if len(got) != len(want) {
		t.Fatalf("merged length = %d, want %d", len(got), len(want))
	}
	for index, msgID := range want {
		if got[index].MsgId != msgID {
			t.Fatalf("merged[%d] = %q, want %q", index, got[index].MsgId, msgID)
		}
	}
}

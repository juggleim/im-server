package tools

import (
	"errors"
	"testing"
)

func TestParseTimeFromMsgIdStrict(t *testing.T) {
	const timestamp int64 = 1790000000123
	id := GenerateMsgId(timestamp, 1, "target")
	if got, err := ParseTimeFromMsgIdStrict(id); err != nil || got != timestamp {
		t.Fatalf("got=%d err=%v id=%s", got, err, id)
	}
	for _, invalid := range []string{"", "short", id + "x", "!" + id[1:], id[:15] + "!", "0000000000000000"} {
		if _, err := ParseTimeFromMsgIdStrict(invalid); !errors.Is(err, ErrInvalidMsgId) {
			t.Fatalf("accepted %q: %v", invalid, err)
		}
	}
}

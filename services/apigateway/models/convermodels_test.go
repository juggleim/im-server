package models

import (
	"encoding/json"
	"testing"
)

func TestGlobalConverTagsReqPresence(t *testing.T) {
	var absent GlobalConverTagsReq
	if err := json.Unmarshal([]byte(`{"conver_id":"conver","channel_type":2}`), &absent); err != nil {
		t.Fatalf("unmarshal absent field: %v", err)
	}
	if absent.GlobalConverTags != nil {
		t.Fatalf("GlobalConverTags = %v, want nil when field is absent", absent.GlobalConverTags)
	}

	var empty GlobalConverTagsReq
	if err := json.Unmarshal([]byte(`{"conver_id":"conver","channel_type":2,"global_conver_tags":[]}`), &empty); err != nil {
		t.Fatalf("unmarshal empty field: %v", err)
	}
	if empty.GlobalConverTags == nil || len(*empty.GlobalConverTags) != 0 {
		t.Fatalf("GlobalConverTags = %v, want non-nil empty slice", empty.GlobalConverTags)
	}
}

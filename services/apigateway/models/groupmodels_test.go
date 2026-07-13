package models

import (
	"encoding/json"
	"testing"
)

func TestGroupInfoGlobalConverTagsPresence(t *testing.T) {
	var absent GroupInfo
	if err := json.Unmarshal([]byte(`{"group_id":"group"}`), &absent); err != nil {
		t.Fatalf("unmarshal absent field: %v", err)
	}
	if absent.GlobalConverTags != nil {
		t.Fatalf("GlobalConverTags = %v, want nil when field is absent", absent.GlobalConverTags)
	}

	var empty GroupInfo
	if err := json.Unmarshal([]byte(`{"group_id":"group","global_conver_tags":[]}`), &empty); err != nil {
		t.Fatalf("unmarshal empty field: %v", err)
	}
	if empty.GlobalConverTags == nil || len(*empty.GlobalConverTags) != 0 {
		t.Fatalf("GlobalConverTags = %v, want non-nil empty slice", empty.GlobalConverTags)
	}
}

func TestGroupMembersReqGlobalConverTagsPresence(t *testing.T) {
	var absent GroupMembersReq
	if err := json.Unmarshal([]byte(`{"group_id":"group"}`), &absent); err != nil {
		t.Fatalf("unmarshal absent field: %v", err)
	}
	if absent.GlobalConverTags != nil {
		t.Fatalf("GlobalConverTags = %v, want nil when field is absent", absent.GlobalConverTags)
	}

	var empty GroupMembersReq
	if err := json.Unmarshal([]byte(`{"group_id":"group","global_conver_tags":[]}`), &empty); err != nil {
		t.Fatalf("unmarshal empty field: %v", err)
	}
	if empty.GlobalConverTags == nil || len(*empty.GlobalConverTags) != 0 {
		t.Fatalf("GlobalConverTags = %v, want non-nil empty slice", empty.GlobalConverTags)
	}
}

package services

import (
	"errors"
	"reflect"
	"testing"
)

func TestPrepareColdGroupMemberAddQueriesOnlyRequestedMembers(t *testing.T) {
	oldCount := countGroupMembersForAdd
	oldQuery := queryExistingGroupMemberIDsForAdd
	defer func() {
		countGroupMembersForAdd = oldCount
		queryExistingGroupMemberIDsForAdd = oldQuery
	}()

	countGroupMembersForAdd = func(appKey, groupID string) (int, error) {
		if appKey != "app" || groupID != "group" {
			t.Fatalf("unexpected key: %s/%s", appKey, groupID)
		}
		return 32000, nil
	}
	queryExistingGroupMemberIDsForAdd = func(appKey, groupID string, memberIDs []string) ([]string, error) {
		if !reflect.DeepEqual(memberIDs, []string{"u1", "u2", "u3"}) {
			t.Fatalf("queried member IDs = %v", memberIDs)
		}
		return []string{"u2"}, nil
	}

	missing, count, _, err := prepareColdGroupMemberAdd("app", "group", []string{"u1", "u2", "u1", "u3"})
	if err != nil {
		t.Fatal(err)
	}
	if count != 32000 {
		t.Fatalf("member count = %d, want 32000", count)
	}
	if !reflect.DeepEqual(missing, []string{"u1", "u3"}) {
		t.Fatalf("missing IDs = %v, want [u1 u3]", missing)
	}
}

func TestPrepareColdGroupMemberAddReturnsQueryError(t *testing.T) {
	oldCount := countGroupMembersForAdd
	oldQuery := queryExistingGroupMemberIDsForAdd
	defer func() {
		countGroupMembersForAdd = oldCount
		queryExistingGroupMemberIDsForAdd = oldQuery
	}()

	wantErr := errors.New("count failed")
	countGroupMembersForAdd = func(string, string) (int, error) { return 0, wantErr }
	queryExistingGroupMemberIDsForAdd = func(string, string, []string) ([]string, error) {
		t.Fatal("existing-member query should not run after count failure")
		return nil, nil
	}

	_, _, _, err := prepareColdGroupMemberAdd("app", "group", []string{"u1"})
	if !errors.Is(err, wantErr) {
		t.Fatalf("error = %v, want %v", err, wantErr)
	}
}

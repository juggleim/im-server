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

	queryCompleted := false
	countGroupMembersForAdd = func(appKey, groupID string) (int, error) {
		if !queryCompleted {
			t.Fatal("count queried before existing members")
		}
		if appKey != "app" || groupID != "group" {
			t.Fatalf("unexpected key: %s/%s", appKey, groupID)
		}
		return 32000, nil
	}
	queryExistingGroupMemberIDsForAdd = func(appKey, groupID string, memberIDs []string) ([]string, error) {
		queryCompleted = true
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
		return nil, nil
	}

	_, _, _, err := prepareColdGroupMemberAdd("app", "group", []string{"u1"})
	if !errors.Is(err, wantErr) {
		t.Fatalf("error = %v, want %v", err, wantErr)
	}
}

func TestPrepareColdGroupMemberAddSkipsCountWhenAllMembersExist(t *testing.T) {
	oldCount, oldQuery := countGroupMembersForAdd, queryExistingGroupMemberIDsForAdd
	t.Cleanup(func() { countGroupMembersForAdd, queryExistingGroupMemberIDsForAdd = oldCount, oldQuery })
	countGroupMembersForAdd = func(string, string) (int, error) {
		t.Fatal("idempotent cold add must skip count")
		return 0, nil
	}
	queryExistingGroupMemberIDsForAdd = func(string, string, []string) ([]string, error) {
		return []string{"u1", "u2"}, nil
	}
	missing, count, trace, err := prepareColdGroupMemberAdd("app", "group", []string{"u1", "u2", "u1"})
	if err != nil || len(missing) != 0 || count != 0 || !trace.countSkipped {
		t.Fatalf("missing=%v count=%d trace=%+v err=%v", missing, count, trace, err)
	}
}

func TestPrepareColdGroupMemberAddStopsOnExistingQueryError(t *testing.T) {
	oldCount, oldQuery := countGroupMembersForAdd, queryExistingGroupMemberIDsForAdd
	t.Cleanup(func() { countGroupMembersForAdd, queryExistingGroupMemberIDsForAdd = oldCount, oldQuery })
	wantErr := errors.New("query failed")
	countGroupMembersForAdd = func(string, string) (int, error) {
		t.Fatal("count must not run after query error")
		return 0, nil
	}
	queryExistingGroupMemberIDsForAdd = func(string, string, []string) ([]string, error) { return nil, wantErr }
	if _, _, _, err := prepareColdGroupMemberAdd("app", "group", []string{"u1"}); !errors.Is(err, wantErr) {
		t.Fatalf("err=%v", err)
	}
}

package services

import (
	"time"

	"im-server/commons/tools"
	"im-server/services/group/dbs"
)

var countGroupMembersForAdd = func(appKey, groupID string) (int, error) {
	return (dbs.GroupMemberDao{}).CountByGroupId(appKey, groupID)
}

var queryExistingGroupMemberIDsForAdd = func(appKey, groupID string, memberIDs []string) ([]string, error) {
	return (dbs.GroupMemberDao{}).QueryExistingMemberIds(appKey, groupID, memberIDs)
}

type groupMemberColdPathTrace struct {
	countDB    time.Duration
	existingDB time.Duration
}

func prepareColdGroupMemberAdd(appKey, groupID string, requestedMemberIDs []string) ([]string, int, groupMemberColdPathTrace, error) {
	trace := groupMemberColdPathTrace{}
	requestedMemberIDs = tools.FilterDuplicated(requestedMemberIDs)
	if len(requestedMemberIDs) == 0 {
		return []string{}, 0, trace, nil
	}

	started := time.Now()
	currentCount, err := countGroupMembersForAdd(appKey, groupID)
	trace.countDB = time.Since(started)
	if err != nil {
		return nil, 0, trace, err
	}

	started = time.Now()
	existingMemberIDs, err := queryExistingGroupMemberIDsForAdd(appKey, groupID, requestedMemberIDs)
	trace.existingDB = time.Since(started)
	if err != nil {
		return nil, 0, trace, err
	}

	existing := make(map[string]struct{}, len(existingMemberIDs))
	for _, memberID := range existingMemberIDs {
		existing[memberID] = struct{}{}
	}
	missingCapacity := len(requestedMemberIDs) - len(existing)
	if missingCapacity < 0 {
		missingCapacity = 0
	}
	missing := make([]string, 0, missingCapacity)
	for _, memberID := range requestedMemberIDs {
		if _, ok := existing[memberID]; !ok {
			missing = append(missing, memberID)
		}
	}
	return missing, currentCount, trace, nil
}

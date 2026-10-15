package services

import (
	"sync"
	"time"
)

type groupMemberStateLockEntry struct {
	mu   sync.Mutex
	refs int
}

// groupMemberStateLockSet provides an exact-key lock without retaining one
// mutex forever for every group ever observed by the process.
type groupMemberStateLockSet struct {
	mu      sync.Mutex
	entries map[string]*groupMemberStateLockEntry
}

func newGroupMemberStateLockSet() *groupMemberStateLockSet {
	return &groupMemberStateLockSet{entries: make(map[string]*groupMemberStateLockEntry)}
}

func (locks *groupMemberStateLockSet) lock(key string) (func(), time.Duration) {
	waitStarted := time.Now()
	locks.mu.Lock()
	entry := locks.entries[key]
	if entry == nil {
		entry = &groupMemberStateLockEntry{}
		locks.entries[key] = entry
	}
	entry.refs++
	locks.mu.Unlock()

	entry.mu.Lock()
	wait := time.Since(waitStarted)
	return func() {
		entry.mu.Unlock()
		locks.mu.Lock()
		entry.refs--
		if entry.refs == 0 {
			delete(locks.entries, key)
		}
		locks.mu.Unlock()
	}, wait
}

func (locks *groupMemberStateLockSet) size() int {
	locks.mu.Lock()
	defer locks.mu.Unlock()
	return len(locks.entries)
}

var groupMemberStateLocks = newGroupMemberStateLockSet()

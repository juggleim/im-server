package services

import (
	"sync"
	"testing"
	"time"
)

func TestGroupMemberStateLockSerializesSameKeyAndReleasesEntry(t *testing.T) {
	locks := newGroupMemberStateLockSet()
	unlockFirst, _ := locks.lock("app\x00group")

	acquiredSecond := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		unlockSecond, _ := locks.lock("app\x00group")
		close(acquiredSecond)
		unlockSecond()
	}()

	select {
	case <-acquiredSecond:
		t.Fatal("second caller acquired the same key before it was released")
	case <-time.After(20 * time.Millisecond):
	}

	unlockFirst()
	select {
	case <-acquiredSecond:
	case <-time.After(time.Second):
		t.Fatal("second caller did not acquire the released key")
	}
	wg.Wait()
	if got := locks.size(); got != 0 {
		t.Fatalf("retained lock entries = %d, want 0", got)
	}
}

func TestGroupMemberStateLockAllowsDifferentKeys(t *testing.T) {
	locks := newGroupMemberStateLockSet()
	unlockFirst, _ := locks.lock("group-a")
	defer unlockFirst()

	done := make(chan struct{})
	go func() {
		unlockSecond, _ := locks.lock("group-b")
		unlockSecond()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("different keys blocked each other")
	}
}

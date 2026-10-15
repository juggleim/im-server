package caches

import (
	"testing"
	"time"
)

func TestSnapshotForKeysDoesNotRefreshOrEvict(t *testing.T) {
	c := NewLruCache("snapshot_test", 2, nil)
	c.Add("first", 1)
	c.Add("second", 2)
	before, ok := c.lru.Peek("first")
	if !ok {
		t.Fatal("missing first cache entry")
	}
	lastRead := before.(*lruCacheItem).lastReadTime
	c.MaxLifeCycle = time.Nanosecond
	entries := c.SnapshotForKeys([]interface{}{"first", "missing"})
	if len(entries) != 2 || !entries[0].Found || entries[0].Value != 1 || entries[1].Found {
		t.Fatalf("unexpected snapshot entries: %+v", entries)
	}
	if c.Len() != 2 || before.(*lruCacheItem).lastReadTime != lastRead {
		t.Fatal("snapshot changed cache size or last read time")
	}
	c.Add("third", 3)
	if _, ok := c.lru.Peek("first"); ok {
		t.Fatal("snapshot changed LRU eviction order")
	}
}

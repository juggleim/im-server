package dbs

import "testing"

func TestAlignCmdMsgSplitPartitionStart(t *testing.T) {
	dayStart := int64(1790726400000)
	if got := alignCmdMsgSplitPartitionStart(dayStart + cmdMsgSplitPartitionIntervalMillis - 1); got != dayStart {
		t.Fatalf("partition start = %d, want %d", got, dayStart)
	}
}

package commonservices

import (
	"im-server/services/commonservices/dbs"
	"testing"
	"time"
)

func TestNextUserActivityCleanupTime(t *testing.T) {
	loc := time.FixedZone("test", 8*60*60)
	tests := []struct {
		name string
		now  time.Time
		want time.Time
	}{
		{
			name: "before schedule",
			now:  time.Date(2026, 6, 18, 4, 59, 0, 0, loc),
			want: time.Date(2026, 6, 18, 5, 0, 0, 0, loc),
		},
		{
			name: "at schedule",
			now:  time.Date(2026, 6, 18, 5, 0, 0, 0, loc),
			want: time.Date(2026, 6, 19, 5, 0, 0, 0, loc),
		},
		{
			name: "after schedule",
			now:  time.Date(2026, 6, 18, 8, 0, 0, 0, loc),
			want: time.Date(2026, 6, 19, 5, 0, 0, 0, loc),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := nextUserActivityCleanupTime(tt.now); !got.Equal(tt.want) {
				t.Fatalf("nextUserActivityCleanupTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserActivityRetentionCutoff(t *testing.T) {
	loc := time.FixedZone("test", 8*60*60)
	now := time.Date(2026, 6, 18, 5, 0, 0, 0, loc)
	want := dayTimeMark(now) - 2*oneDay

	if got := userActivityRetentionCutoff(now); got != want {
		t.Fatalf("userActivityRetentionCutoff() = %d, want %d", got, want)
	}
}

func TestUserActivityPreviousDayTimeMark(t *testing.T) {
	loc := time.FixedZone("test", 8*60*60)
	now := time.Date(2026, 6, 18, 5, 0, 0, 0, loc)
	want := dayTimeMark(now) - oneDay

	if got := userActivityPreviousDayTimeMark(now); got != want {
		t.Fatalf("userActivityPreviousDayTimeMark() = %d, want %d", got, want)
	}
}

func TestMsgRealtimeStatRetentionCutoff(t *testing.T) {
	loc := time.FixedZone("test", 8*60*60)
	now := time.Date(2026, 6, 18, 5, 0, 30, 0, loc)
	want := dayTimeMark(now) - int64(msgRealtimeStatCleanupRetainDays-1)*oneDay

	if got := msgRealtimeStatRetentionCutoff(now); got != want {
		t.Fatalf("msgRealtimeStatRetentionCutoff() = %d, want %d", got, want)
	}
}

func TestConnectCountCompactDayEnd(t *testing.T) {
	dayMark := int64(1000)
	want := dayMark + oneDay

	if got := connectCountCompactDayEnd(dayMark); got != want {
		t.Fatalf("connectCountCompactDayEnd() = %d, want %d", got, want)
	}
}

func TestConnectCountCompactDayRange(t *testing.T) {
	loc := time.FixedZone("test", 8*60*60)
	now := time.Date(2026, 6, 18, 5, 0, 0, 0, loc)
	cutoff := userActivityRetentionCutoff(now)
	wantStart := cutoff - oneDay
	wantEnd := cutoff

	gotStart, gotEnd := connectCountCompactDayRange(now)
	if gotStart != wantStart || gotEnd != wantEnd {
		t.Fatalf("connectCountCompactDayRange() = (%d, %d), want (%d, %d)", gotStart, gotEnd, wantStart, wantEnd)
	}
}

func TestConnectCountKeeperIDList(t *testing.T) {
	keepers := map[connectCountKeeperKey]int64{
		{appKey: "app1", connectType: 1}: 10,
		{appKey: "app2", connectType: 2}: 20,
	}
	got := connectCountKeeperIDList(keepers)
	if len(got) != 2 {
		t.Fatalf("len(keeperIDs) = %d, want 2", len(got))
	}
	seen := map[int64]struct{}{}
	for _, id := range got {
		seen[id] = struct{}{}
	}
	if _, ok := seen[10]; !ok {
		t.Fatalf("keeperIDs missing 10: %v", got)
	}
	if _, ok := seen[20]; !ok {
		t.Fatalf("keeperIDs missing 20: %v", got)
	}
}

func TestIsBetterConnectCountKeeper(t *testing.T) {
	tests := []struct {
		name      string
		candidate *dbs.ConnectCountDao
		current   *dbs.ConnectCountDao
		want      bool
	}{
		{
			name:      "nil candidate",
			candidate: nil,
			current:   &dbs.ConnectCountDao{ID: 1, Count: 10},
			want:      false,
		},
		{
			name:      "nil current",
			candidate: &dbs.ConnectCountDao{ID: 1, Count: 10},
			current:   nil,
			want:      true,
		},
		{
			name:      "higher count wins",
			candidate: &dbs.ConnectCountDao{ID: 2, Count: 11},
			current:   &dbs.ConnectCountDao{ID: 1, Count: 10},
			want:      true,
		},
		{
			name:      "lower count loses",
			candidate: &dbs.ConnectCountDao{ID: 2, Count: 9},
			current:   &dbs.ConnectCountDao{ID: 1, Count: 10},
			want:      false,
		},
		{
			name:      "tie keeps lower id",
			candidate: &dbs.ConnectCountDao{ID: 1, Count: 10},
			current:   &dbs.ConnectCountDao{ID: 2, Count: 10},
			want:      true,
		},
		{
			name:      "tie rejects higher id",
			candidate: &dbs.ConnectCountDao{ID: 3, Count: 10},
			current:   &dbs.ConnectCountDao{ID: 2, Count: 10},
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isBetterConnectCountKeeper(tt.candidate, tt.current); got != tt.want {
				t.Fatalf("isBetterConnectCountKeeper() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAggregateUserActivityScanRows(t *testing.T) {
	rows := []dbs.UserActivityScanRow{
		{ID: 1, AppKey: "app1", TimeMark: 1000},
		{ID: 2, AppKey: "app1", TimeMark: 1000},
		{ID: 3, AppKey: "app2", TimeMark: 1000},
		{ID: 4, AppKey: "", TimeMark: 1000},
	}

	got := aggregateUserActivityScanRows(rows)
	if len(got) != 2 {
		t.Fatalf("len(aggregates) = %d, want 2", len(got))
	}
	if got[dailyActivityBucketKey{appKey: "app1", timeMark: 1000}] != 2 {
		t.Fatalf("app1@1000 = %d, want 2", got[dailyActivityBucketKey{appKey: "app1", timeMark: 1000}])
	}
	if got[dailyActivityBucketKey{appKey: "app2", timeMark: 1000}] != 1 {
		t.Fatalf("app2@1000 = %d, want 1", got[dailyActivityBucketKey{appKey: "app2", timeMark: 1000}])
	}
}

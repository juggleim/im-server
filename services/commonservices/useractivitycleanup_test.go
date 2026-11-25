package commonservices

import (
	"errors"
	"im-server/services/commonservices/dbs"
	"reflect"
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

type activityArchiveCall struct {
	day    int64
	remove bool
}
type activityArchiveFake struct {
	days  []int64
	calls []activityArchiveCall
	err   error
	today int64
}

func (f *activityArchiveFake) DaysBefore(today int64) ([]int64, error) {
	f.today = today
	return f.days, nil
}
func (f *activityArchiveFake) ArchiveDay(day int64, remove bool) error {
	f.calls = append(f.calls, activityArchiveCall{day, remove})
	return f.err
}

func TestActivityArchiveKeepsTodayAndRetentionBoundary(t *testing.T) {
	now := time.Date(2026, 11, 9, 4, 0, 0, 0, time.FixedZone("UTC+8", 8*3600))
	today := dayTimeMark(now)
	fake := &activityArchiveFake{days: []int64{today - 4*oneDay, today - 2*oneDay, today - oneDay, today, today + oneDay, 1, -1}}
	if err := archiveUserActivityDays(fake, now); err != nil {
		t.Fatal(err)
	}
	want := []activityArchiveCall{{today - 4*oneDay, true}, {today - 2*oneDay, false}, {today - oneDay, false}}
	if fake.today != today || !reflect.DeepEqual(fake.calls, want) {
		t.Fatalf("today=%d calls=%v", fake.today, fake.calls)
	}
}

func TestActivityArchiveStopsOnFailure(t *testing.T) {
	now := time.Now()
	today := dayTimeMark(now)
	failure := errors.New("summary failed")
	fake := &activityArchiveFake{days: []int64{today - 4*oneDay, today - 3*oneDay}, err: failure}
	if err := archiveUserActivityDays(fake, now); !errors.Is(err, failure) || len(fake.calls) != 1 {
		t.Fatalf("calls=%v err=%v", fake.calls, err)
	}
}

func TestActivityDailyWriteUsesMillisecondsWithoutChangingMessageStats(t *testing.T) {
	now := time.Now()
	day := getUserActivityTimeMark()
	if day != dayTimeMark(now) {
		t.Fatalf("activity writer is not using UTC milliseconds: %d", day)
	}
	if day%oneDay != 0 || day/1000 != getDbTimeMark() {
		t.Fatalf("activity day=%d msgstat day=%d", day, getDbTimeMark())
	}
}

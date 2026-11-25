package main

import (
	"bytes"
	"context"
	"errors"
	"io"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestStreamDailyActivitiesUTCAndIncrementalOutput(t *testing.T) {
	var output bytes.Buffer
	var dates []string
	now := time.Date(2026, 10, 12, 2, 0, 0, 0, time.FixedZone("CST", 8*60*60))
	err := streamDailyActivities(context.Background(), &output, "app-a", now, 3, func(_ context.Context, date time.Time) (int64, error) {
		if len(dates) > 0 && !strings.Contains(output.String(), dates[len(dates)-1]) {
			t.Fatal("previous day must be printed before next query")
		}
		dates = append(dates, date.Format(time.DateOnly))
		return 7, nil
	})
	if err != nil || !reflect.DeepEqual(dates, []string{"2026-10-11", "2026-10-10", "2026-10-09"}) {
		t.Fatalf("dates=%v err=%v", dates, err)
	}
}

func TestStreamDailyActivitiesStopsOnError(t *testing.T) {
	var output bytes.Buffer
	want := errors.New("query failed")
	queries := 0
	err := streamDailyActivities(context.Background(), &output, "app", time.Now(), 3, func(context.Context, time.Time) (int64, error) {
		queries++
		if queries == 2 {
			return 0, want
		}
		return 1, nil
	})
	if !errors.Is(err, want) || queries != 2 || strings.Count(output.String(), "\n") != 2 {
		t.Fatalf("queries=%d output=%q err=%v", queries, output.String(), err)
	}
}

func TestCLIRejectsInvalidInputBeforeConnecting(t *testing.T) {
	for _, args := range [][]string{{}, {"--app-key", " "}, {"--app-key", "app", "--days", "0"}, {"--app-key", "app", "--days", "3661"}} {
		if err := run(context.Background(), args, io.Discard, io.Discard); err == nil {
			t.Fatalf("args=%v unexpectedly accepted", args)
		}
	}
}

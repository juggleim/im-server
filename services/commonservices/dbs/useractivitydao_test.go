package dbs

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type activityTestStep struct {
	query   string
	args    []driver.Value
	columns []string
	rows    [][]driver.Value
	err     error
}
type activityTestDriver struct {
	t     *testing.T
	steps []activityTestStep
}
type activityTestConn struct{ script *activityTestDriver }
type activityTestRows struct {
	columns []string
	rows    [][]driver.Value
}

var activityDriverID atomic.Uint64

func (d *activityTestDriver) Open(string) (driver.Conn, error) { return &activityTestConn{d}, nil }
func (c *activityTestConn) Prepare(string) (driver.Stmt, error) {
	return nil, errors.New("unexpected prepare")
}
func (c *activityTestConn) Close() error { return nil }
func (c *activityTestConn) next(query string, args []driver.NamedValue) activityTestStep {
	s := c.script
	if len(s.steps) == 0 {
		s.t.Fatalf("unexpected SQL: %s", query)
	}
	step := s.steps[0]
	s.steps = s.steps[1:]
	values := make([]driver.Value, len(args))
	for i := range args {
		values[i] = args[i].Value
	}
	if !strings.Contains(query, step.query) || !reflect.DeepEqual(values, step.args) && !(len(values) == 0 && len(step.args) == 0) {
		s.t.Fatalf("SQL=%s args=%v; want %s %v", query, values, step.query, step.args)
	}
	return step
}
func (c *activityTestConn) Begin() (driver.Tx, error) {
	step := c.next("BEGIN", nil)
	return c, step.err
}
func (c *activityTestConn) Commit() error   { return c.next("COMMIT", nil).err }
func (c *activityTestConn) Rollback() error { return c.next("ROLLBACK", nil).err }
func (c *activityTestConn) ExecContext(_ context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	step := c.next(query, args)
	return driver.RowsAffected(1), step.err
}
func (c *activityTestConn) QueryContext(_ context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	step := c.next(query, args)
	return &activityTestRows{step.columns, step.rows}, step.err
}
func (r *activityTestRows) Columns() []string { return r.columns }
func (r *activityTestRows) Close() error      { return nil }
func (r *activityTestRows) Next(values []driver.Value) error {
	if len(r.rows) == 0 {
		return io.EOF
	}
	copy(values, r.rows[0])
	r.rows = r.rows[1:]
	return nil
}

func activityTestDB(t *testing.T, steps []activityTestStep) *gorm.DB {
	t.Helper()
	script := &activityTestDriver{t: t, steps: steps}
	name := fmt.Sprintf("activity-test-%d", activityDriverID.Add(1))
	sql.Register(name, script)
	pool, err := sql.Open(name, "")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		pool.Close()
		if len(script.steps) > 0 {
			t.Errorf("unexecuted steps=%v", script.steps)
		}
	})
	db, err := gorm.Open(mysql.New(mysql.Config{Conn: pool, SkipInitializeWithVersion: true}), &gorm.Config{DisableAutomaticPing: true, SkipDefaultTransaction: true, Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func TestArchiveActivityDayDeduplicatesUnitsAndReplacesSummaryOnRetry(t *testing.T) {
	day := time.Date(2026, 11, 6, 0, 0, 0, 0, time.UTC).UnixMilli()
	for _, remove := range []bool{false, true} {
		var steps []activityTestStep
		for range 2 {
			steps = append(steps,
				activityTestStep{query: "BEGIN"},
				activityTestStep{query: "COUNT(DISTINCT user_id)", args: []driver.Value{day, day / 1000}, columns: []string{"app_key", "count"}, rows: [][]driver.Value{{"app", int64(2)}}},
				activityTestStep{query: "ON DUPLICATE KEY UPDATE count=VALUES(count)", args: []driver.Value{"app", day, int64(2)}},
				activityTestStep{query: "DELETE FROM dailyactivities", args: []driver.Value{"app", day / 1000}},
			)
			if remove {
				steps = append(steps, activityTestStep{query: "DELETE FROM `useractivities` WHERE time_mark IN (?, ?)", args: []driver.Value{day, day / 1000}})
			}
			steps = append(steps, activityTestStep{query: "COMMIT"})
		}
		db := activityTestDB(t, steps)
		for range 2 {
			if err := archiveActivityDay(db, day, remove); err != nil {
				t.Fatal(err)
			}
		}
	}
}

func TestArchiveActivityDayRollsBackBeforeDeletingOnSummaryFailure(t *testing.T) {
	day := time.Date(2026, 11, 6, 0, 0, 0, 0, time.UTC).UnixMilli()
	failure := errors.New("write failed")
	for _, stage := range []string{"read", "summary", "legacy", "delete"} {
		t.Run(stage, func(t *testing.T) {
			steps := []activityTestStep{{query: "BEGIN"}, {query: "COUNT(DISTINCT user_id)", args: []driver.Value{day, day / 1000}, columns: []string{"app_key", "count"}, rows: [][]driver.Value{{"app", int64(2)}}}}
			if stage == "read" {
				steps[len(steps)-1].err = failure
			} else {
				steps = append(steps, activityTestStep{query: "INSERT INTO dailyactivities", args: []driver.Value{"app", day, int64(2)}})
				if stage == "summary" {
					steps[len(steps)-1].err = failure
				} else {
					steps = append(steps, activityTestStep{query: "DELETE FROM dailyactivities", args: []driver.Value{"app", day / 1000}})
					if stage == "legacy" {
						steps[len(steps)-1].err = failure
					} else {
						steps = append(steps, activityTestStep{query: "DELETE FROM `useractivities`", args: []driver.Value{day, day / 1000}, err: failure})
					}
				}
			}
			steps = append(steps, activityTestStep{query: "ROLLBACK"})
			if err := archiveActivityDay(activityTestDB(t, steps), day, true); !errors.Is(err, failure) {
				t.Fatalf("err=%v", err)
			}
		})
	}
}

func TestActivityDaysBeforeUsesSeparateSecondAndMillisecondCutoffs(t *testing.T) {
	today := time.Date(2026, 11, 9, 0, 0, 0, 0, time.UTC).UnixMilli()
	db := activityTestDB(t, []activityTestStep{{
		query:   "SELECT DISTINCT CASE WHEN time_mark < ? THEN time_mark * 1000 ELSE time_mark END AS day_mark FROM `useractivities` WHERE (time_mark > 0 AND time_mark < ?) OR (time_mark >= ? AND time_mark < ?)",
		args:    []driver.Value{activityMillisBoundary, today / 1000, activityMillisBoundary, today},
		columns: []string{"day_mark"}, rows: [][]driver.Value{{today - 86400000}, {today - 3*86400000}},
	}})
	days, err := activityDaysBefore(db, today)
	if err != nil || !reflect.DeepEqual(days, []int64{today - 86400000, today - 3*86400000}) {
		t.Fatalf("days=%v err=%v", days, err)
	}
}

package dbcommons

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync/atomic"
	"testing"

	mysqlDriver "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// A scripted in-process SQL driver exercises the real pinned-connection runner
// without adding a dependency or requiring a database server.
type migrationStep struct {
	query string
	value driver.Value
	err   error
	exec  bool
}

type migrationScript struct {
	t     *testing.T
	steps []migrationStep
	opens int
}

type migrationConn struct{ script *migrationScript }
type migrationRows struct {
	value driver.Value
	done  bool
}

var migrationDriverID atomic.Uint64

func (s *migrationScript) Open(string) (driver.Conn, error) {
	s.opens++
	if s.opens != 1 {
		s.t.Error("migration opened more than one physical connection")
	}
	return &migrationConn{s}, nil
}
func (c *migrationConn) Prepare(string) (driver.Stmt, error) {
	return nil, errors.New("unexpected prepare")
}
func (c *migrationConn) Close() error              { return nil }
func (c *migrationConn) Begin() (driver.Tx, error) { return nil, errors.New("unexpected transaction") }
func (c *migrationConn) next(query string, exec bool) migrationStep {
	s := c.script
	if len(s.steps) == 0 {
		s.t.Fatalf("unexpected query %s", query)
	}
	step := s.steps[0]
	s.steps = s.steps[1:]
	if !strings.Contains(query, step.query) || step.exec != exec {
		s.t.Fatalf("query=%s exec=%t, want %+v", query, exec, step)
	}
	return step
}
func (c *migrationConn) QueryContext(_ context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	step := c.next(query, false)
	if strings.Contains(query, "GET_LOCK") && (len(args) != 2 || args[0].Value != migrationLockName("test") || args[1].Value != int64(300)) {
		c.script.t.Fatalf("lock arguments=%v", args)
	}
	if step.err != nil {
		return nil, step.err
	}
	return &migrationRows{value: step.value}, nil
}
func (c *migrationConn) ExecContext(_ context.Context, query string, _ []driver.NamedValue) (driver.Result, error) {
	step := c.next(query, true)
	return driver.RowsAffected(1), step.err
}
func (r *migrationRows) Columns() []string { return []string{"value"} }
func (r *migrationRows) Close() error      { return nil }
func (r *migrationRows) Next(values []driver.Value) error {
	if r.done {
		return io.EOF
	}
	r.done = true
	values[0] = r.value
	return nil
}

func runMigrationScript(t *testing.T, steps []migrationStep) error {
	t.Helper()
	script := &migrationScript{t: t, steps: steps}
	name := fmt.Sprintf("migration-test-%d", migrationDriverID.Add(1))
	sql.Register(name, script)
	pool, err := sql.Open(name, "")
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	testDB, err := gorm.Open(mysql.New(mysql.Config{Conn: pool, SkipInitializeWithVersion: true}), &gorm.Config{DisableAutomaticPing: true})
	if err != nil {
		t.Fatal(err)
	}
	previous := db
	db = testDB
	defer func() { db = previous }()
	err = Upgrade()
	if len(script.steps) != 0 {
		t.Fatalf("unexecuted steps=%+v", script.steps)
	}
	return err
}

func TestUpgradePinnedConnectionAndFailureStops(t *testing.T) {
	migrations, err := listMigrationFiles()
	if err != nil {
		t.Fatal(err)
	}
	latest := migrations[len(migrations)-1]
	previous := fmt.Sprint(migrations[len(migrations)-2].version)
	var statements []string
	content, err := sqlFs.ReadFile(latest.path)
	if err != nil {
		t.Fatal(err)
	}
	if err := executeSQLReader(context.Background(), strings.NewReader(string(content)), latest.path, func(_ context.Context, query string) error { statements = append(statements, query); return nil }); err != nil {
		t.Fatal(err)
	}
	failure := errors.New("database rejected operation")
	for _, stage := range []string{"lock", "null_lock", "read", "invalid", "zero", "negative", "newer", "sql", "version", "release", "success"} {
		t.Run(stage, func(t *testing.T) {
			steps := []migrationStep{{query: "SELECT DATABASE()", value: "test"}, {query: "SELECT GET_LOCK", value: int64(1)}}
			if stage == "lock" || stage == "null_lock" {
				steps[1].value = int64(0)
				if stage == "null_lock" {
					steps[1].value = nil
				}
			} else {
				read := migrationStep{query: "SELECT conf_value", value: previous}
				switch stage {
				case "read":
					read.err = failure
				case "invalid":
					read.value = "invalid"
				case "zero":
					read.value = "0"
				case "negative":
					read.value = "-1"
				case "newer":
					read.value = "99999999"
				}
				steps = append(steps, read)
				if stage != "read" && stage != "invalid" && stage != "zero" && stage != "negative" && stage != "newer" {
					for i, query := range statements {
						step := migrationStep{query: query, exec: true}
						if stage == "sql" && i == 0 {
							step.err = failure
						}
						steps = append(steps, step)
						if step.err != nil {
							break
						}
					}
					if stage != "sql" {
						step := migrationStep{query: "INSERT INTO globalconfs", exec: true}
						if stage == "version" {
							step.err = failure
						}
						steps = append(steps, step)
					}
				}
				release := migrationStep{query: "SELECT RELEASE_LOCK", value: int64(1)}
				if stage == "release" {
					release.err = failure
				}
				steps = append(steps, release)
			}
			err := runMigrationScript(t, steps)
			if stage == "success" {
				if err != nil {
					t.Fatal(err)
				}
				return
			}
			if err == nil {
				t.Fatal("failure accepted")
			}
			if stage == "read" || stage == "sql" || stage == "version" || stage == "release" {
				if !errors.Is(err, failure) {
					t.Fatalf("lost error: %v", err)
				}
			}
		})
	}
}

func TestUpgradeBootstrapOnlyOnMissingVersion(t *testing.T) {
	for _, readErr := range []error{sql.ErrNoRows, &mysqlDriver.MySQLError{Number: 1146}, &mysqlDriver.MySQLError{Number: 1142}} {
		t.Run(readErr.Error(), func(t *testing.T) {
			steps := []migrationStep{{query: "SELECT DATABASE()", value: "test"}, {query: "SELECT GET_LOCK", value: int64(1)}, {query: "SELECT conf_value", err: readErr}}
			failure := errors.New("stop bootstrap")
			missing := errors.Is(readErr, sql.ErrNoRows)
			if e, ok := readErr.(*mysqlDriver.MySQLError); ok && e.Number == 1146 {
				missing = true
			}
			if missing {
				steps = append(steps, migrationStep{query: "CREATE TABLE", exec: true, err: failure})
			} else {
				failure = readErr
			}
			steps = append(steps, migrationStep{query: "SELECT RELEASE_LOCK", value: int64(1)})
			if err := runMigrationScript(t, steps); !errors.Is(err, failure) {
				t.Fatalf("err=%v", err)
			}
		})
	}
}

func TestUpgradeEmptyDatabaseRunsCompleteMigrationChainOnOneConnection(t *testing.T) {
	steps := []migrationStep{
		{query: "SELECT DATABASE()", value: "test"},
		{query: "SELECT GET_LOCK", value: int64(1)},
		{query: "SELECT conf_value", err: &mysqlDriver.MySQLError{Number: 1146}},
	}
	migrations, err := listMigrationFiles()
	if err != nil {
		t.Fatal(err)
	}
	for _, migration := range migrations {
		content, err := sqlFs.ReadFile(migration.path)
		if err != nil {
			t.Fatal(err)
		}
		if err := executeSQLReader(context.Background(), strings.NewReader(string(content)), migration.path, func(_ context.Context, query string) error {
			steps = append(steps, migrationStep{query: query, exec: true})
			return nil
		}); err != nil {
			t.Fatal(err)
		}
		steps = append(steps, migrationStep{query: "INSERT INTO globalconfs", exec: true})
	}
	steps = append(steps, migrationStep{query: "SELECT RELEASE_LOCK", value: int64(1)})
	if err := runMigrationScript(t, steps); err != nil {
		t.Fatal(err)
	}
}

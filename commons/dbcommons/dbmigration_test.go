package dbcommons

import (
	"context"
	"errors"
	"strings"
	"testing"
)

func Test20260908MigrationIsSingleAlter(t *testing.T) {
	content, err := sqlFs.ReadFile("sqls/20260908.sql")
	if err != nil {
		t.Fatal(err)
	}
	sql := string(content)
	if strings.Count(sql, ";") != 1 || strings.Count(strings.ToUpper(sql), "ALTER TABLE") != 1 {
		t.Fatal("20260908 migration must use one ALTER TABLE statement")
	}
	for _, column := range []string{"auth_type", "p8_key_id", "p8_team_id", "p8_private_key", "p8_key_name", "config_version"} {
		if strings.Count(sql, "`"+column+"`") != 1 {
			t.Fatalf("migration must add %s exactly once", column)
		}
	}
}

func TestExecuteSqlFileReturnsExecError(t *testing.T) {
	want := errors.New("statement failed")
	calls := 0
	err := executeSQLReader(context.Background(), strings.NewReader("SELECT 1;\nSELECT 2;\n"), "test.sql", func(context.Context, string) error {
		calls++
		return want
	})
	if !errors.Is(err, want) || calls != 1 || !strings.Contains(err.Error(), "test.sql statement 1") {
		t.Fatalf("calls=%d err=%v", calls, err)
	}
}

func TestMigrationReaderExecutesTrailingStatementAndRejectsOversizeLine(t *testing.T) {
	var statements []string
	err := executeSQLReader(context.Background(), strings.NewReader("-- comment\nSELECT\n 1;\nSELECT 2"), "test.sql", func(_ context.Context, query string) error {
		statements = append(statements, query)
		return nil
	})
	if err != nil || len(statements) != 2 || statements[1] != "SELECT 2" {
		t.Fatalf("statements=%v err=%v", statements, err)
	}
	err = executeSQLReader(context.Background(), strings.NewReader(strings.Repeat("x", migrationMaxLineSize+1)), "test.sql", func(context.Context, string) error { t.Fatal("oversize statement executed"); return nil })
	if err == nil {
		t.Fatal("expected scanner error")
	}
}

func TestMigrationFilesDoNotAdvanceVersionThemselves(t *testing.T) {
	migrations, err := listMigrationFiles()
	if err != nil {
		t.Fatal(err)
	}
	for i, migration := range migrations {
		if i > 0 && migrations[i-1].version >= migration.version {
			t.Fatal("migrations not ordered")
		}
		content, err := sqlFs.ReadFile(migration.path)
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(string(content), JimDbVersionKey) {
			t.Fatalf("%s advances version outside runner", migration.path)
		}
	}
}

func TestMigrationLockNameBoundedAndDatabaseScoped(t *testing.T) {
	if migrationLockName("a") == migrationLockName("b") {
		t.Fatal("database locks collide")
	}
	if len(migrationLockName(strings.Repeat("x", 100))) > 64 {
		t.Fatal("lock name exceeds MySQL limit")
	}
}

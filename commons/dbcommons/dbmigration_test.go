package dbcommons

import (
	"os"
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
	content, err := os.ReadFile("dbmigration.go")
	if err != nil {
		t.Fatal(err)
	}
	source := string(content)
	want := "if err := GetDb().Exec(query).Error; err != nil {\n\t\t\t\t\tfmt.Println(\"[DbMigration_Err]Execute sql error:\", err, query)\n\t\t\t\t\treturn err"
	if !strings.Contains(source, want) {
		t.Fatal("executeSqlFile must immediately return an Exec error")
	}
}

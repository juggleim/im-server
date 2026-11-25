package dbcommons

import (
	"strings"
	"testing"
)

func TestBanUserRuleMigrationOnlyCreatesRuleTable(t *testing.T) {
	content, err := sqlFs.ReadFile("sqls/20261103.sql")
	if err != nil {
		t.Fatal(err)
	}
	sql := string(content)
	if strings.Count(sql, ";") != 1 {
		t.Fatalf("ban-user-rule migration has %d statements, want only the table creation", strings.Count(sql, ";"))
	}
	for _, check := range []string{
		"CREATE TABLE IF NOT EXISTS `banuserrules`",
		"UNIQUE KEY `uniq_ban_rule` (`app_key`,`user_id`,`scope_key`,`scope_value_hash`)",
	} {
		if !strings.Contains(sql, check) {
			t.Errorf("migration missing %q", check)
		}
	}
	if strings.Contains(strings.ToUpper(sql), "INSERT") || strings.Contains(sql, "FROM `banusers`") {
		t.Fatal("migration must not copy legacy banusers data into banuserrules")
	}
}

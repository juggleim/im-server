package dbcommons

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"testing"
)

// Replay the declared table/column/index inventory to catch duplicate ADDs on
// a fresh bootstrap. This checks the embedded migration chain, not MySQL syntax.
func TestMigrationChainHasNoUnguardedDuplicateDDL(t *testing.T) {
	createRE := regexp.MustCompile("(?i)^CREATE TABLE( IF NOT EXISTS)? `([^`]+)`")
	columnRE := regexp.MustCompile("(?m)^\\s*`([^`]+)`\\s+\\w+")
	indexRE := regexp.MustCompile("(?i)(?:UNIQUE )?(?:KEY|INDEX) `([^`]+)`")
	alterRE := regexp.MustCompile("(?i)ALTER TABLE `([^`]+)`")
	addColumnRE := regexp.MustCompile("(?i)ADD COLUMN `([^`]+)`")
	addIndexRE := regexp.MustCompile("(?i)ADD (?:UNIQUE )?(?:KEY|INDEX) `([^`]+)`")
	dropIndexRE := regexp.MustCompile("(?i)DROP INDEX `([^`]+)`")
	type table struct{ columns, indexes map[string]bool }
	tables := map[string]table{}
	files, err := listMigrationFiles()
	if err != nil {
		t.Fatal(err)
	}
	guardedDuplicates := map[string]bool{}
	for _, file := range files {
		content, err := sqlFs.ReadFile(file.path)
		if err != nil {
			t.Fatal(err)
		}
		err = executeSQLReader(context.Background(), strings.NewReader(string(content)), file.path, func(_ context.Context, query string) error {
			if match := createRE.FindStringSubmatch(query); match != nil {
				name := match[2]
				if _, exists := tables[name]; exists {
					if match[1] == "" {
						return fmt.Errorf("duplicate CREATE %s", name)
					}
					return nil
				}
				item := table{map[string]bool{}, map[string]bool{}}
				// The file runner joins lines. Reinsert field boundaries for the
				// simple declaration matcher without splitting commas in types.
				for _, field := range columnRE.FindAllStringSubmatch(strings.ReplaceAll(query, ", `", ",\n `"), -1) {
					item.columns[field[1]] = true
				}
				for _, index := range indexRE.FindAllStringSubmatch(query, -1) {
					item.indexes[index[1]] = true
				}
				tables[name] = item
				return nil
			}
			match := alterRE.FindStringSubmatch(query)
			if match == nil {
				return nil
			}
			name := match[1]
			item, exists := tables[name]
			if !exists {
				return fmt.Errorf("ALTER before CREATE %s", name)
			}
			guarded := strings.HasPrefix(query, "SET @") && strings.Contains(query, "TABLE_SCHEMA = DATABASE()") && strings.Contains(query, "TABLE_NAME = '"+name+"'") && strings.Contains(query, "'SELECT 1'")
			for _, index := range dropIndexRE.FindAllStringSubmatch(query, -1) {
				if !item.indexes[index[1]] {
					return fmt.Errorf("DROP missing index %s.%s", name, index[1])
				}
				delete(item.indexes, index[1])
			}
			for _, field := range addColumnRE.FindAllStringSubmatch(query, -1) {
				if item.columns[field[1]] {
					if !guarded || !strings.Contains(query, "COLUMN_NAME = '"+field[1]+"'") {
						return fmt.Errorf("duplicate column %s.%s", name, field[1])
					}
					guardedDuplicates[name+"."+field[1]] = true
				}
				item.columns[field[1]] = true
			}
			for _, index := range addIndexRE.FindAllStringSubmatch(query, -1) {
				if item.indexes[index[1]] {
					if !guarded || !strings.Contains(query, "INDEX_NAME = '"+index[1]+"'") {
						return fmt.Errorf("duplicate index %s.%s", name, index[1])
					}
					guardedDuplicates[name+"."+index[1]] = true
				}
				item.indexes[index[1]] = true
			}
			return nil
		})
		if err != nil {
			t.Fatalf("%s: %v", file.path, err)
		}
	}
	for _, key := range []string{"appnavs.admin_url", "useractivities.idx_time_mark_id"} {
		if !guardedDuplicates[key] {
			t.Errorf("known duplicate %s was not exercised", key)
		}
	}
}

func TestMigrationDDLGuardsUsePinnedPreparedStatements(t *testing.T) {
	for _, name := range []string{"20260720", "20260813"} {
		content, err := sqlFs.ReadFile("sqls/" + name + ".sql")
		if err != nil {
			t.Fatal(err)
		}
		var statements []string
		if err := executeSQLReader(context.Background(), strings.NewReader(string(content)), name, func(_ context.Context, query string) error { statements = append(statements, query); return nil }); err != nil {
			t.Fatal(err)
		}
		for i, query := range statements {
			if !strings.HasPrefix(query, "SET @") {
				continue
			}
			if i+3 >= len(statements) || !strings.HasPrefix(statements[i+1], "PREPARE ") || !strings.HasPrefix(statements[i+2], "EXECUTE ") || !strings.HasPrefix(statements[i+3], "DEALLOCATE PREPARE ") {
				t.Fatalf("incomplete guard lifecycle in %s", name)
			}
			if strings.Contains(query, "information_schema.STATISTICS") {
				for _, check := range []string{"NON_UNIQUE = 1", "INDEX_TYPE = 'BTREE'", "SUB_PART IS NULL", "COLLATION = 'A'", "IS_VISIBLE = 'YES'", "SEQ_IN_INDEX"} {
					if !strings.Contains(query, check) {
						t.Errorf("%s missing exact index guard %s", name, check)
					}
				}
			}
		}
	}
}

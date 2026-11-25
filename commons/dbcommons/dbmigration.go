package dbcommons

import (
	"bufio"
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"hash/crc32"
	"im-server/commons/tools"
	"io"
	"sort"
	"strings"
	"time"

	mysqlDriver "github.com/go-sql-driver/mysql"
)

//go:embed sqls/*
var sqlFs embed.FS

const (
	JimDbVersionKey      = "jimdb_version"
	migrationLockWait    = 5 * time.Minute
	migrationRunTimeout  = 60 * time.Minute
	migrationReleaseWait = 3 * time.Second
	migrationMaxLineSize = 1024 * 1024
)

type migrationFile struct {
	version int64
	path    string
}

// Upgrade holds a database-scoped advisory lock on the same connection used
// for every statement and version write. A failed migration stops startup.
func Upgrade() (retErr error) {
	if GetDb() == nil {
		return errors.New("mysql is not initialized")
	}
	sqlDB, err := GetDb().DB()
	if err != nil {
		return fmt.Errorf("get mysql connection pool: %w", err)
	}
	acquireCtx, acquireCancel := context.WithTimeout(context.Background(), migrationLockWait+time.Minute)
	defer acquireCancel()
	conn, err := sqlDB.Conn(acquireCtx)
	if err != nil {
		return fmt.Errorf("get migration connection: %w", err)
	}
	defer conn.Close()
	var dbName string
	if err := conn.QueryRowContext(acquireCtx, "SELECT DATABASE()").Scan(&dbName); err != nil {
		return fmt.Errorf("read migration database name: %w", err)
	}
	lockName := migrationLockName(dbName)
	locked, err := acquireMigrationLock(acquireCtx, conn, lockName, migrationLockWait)
	if err != nil {
		return err
	}
	if !locked {
		return fmt.Errorf("timed out waiting %s for migration lock %q", migrationLockWait, lockName)
	}
	acquireCancel()
	defer func() {
		releaseCtx, cancel := context.WithTimeout(context.Background(), migrationReleaseWait)
		defer cancel()
		retErr = errors.Join(retErr, releaseMigrationLock(releaseCtx, conn, lockName))
	}()
	ctx, cancel := context.WithTimeout(context.Background(), migrationRunTimeout)
	defer cancel()
	version, err := readMigrationVersion(ctx, conn)
	if err != nil {
		return err
	}
	fmt.Println("[DbMigration]current version:", version)
	migrations, err := listMigrationFiles()
	if err != nil {
		return err
	}
	if len(migrations) > 0 && version > migrations[len(migrations)-1].version {
		return fmt.Errorf("database version %d is newer than binary version %d", version, migrations[len(migrations)-1].version)
	}
	for _, migration := range migrations {
		if migration.version <= version {
			continue
		}
		fmt.Println("[DbMigration]start to execute sql file:", migration.path)
		if err := executeSQLFile(ctx, conn, migration.path); err != nil {
			return err
		}
		if err := writeMigrationVersion(ctx, conn, migration.version); err != nil {
			return fmt.Errorf("update migration version to %d: %w", migration.version, err)
		}
		version = migration.version
		fmt.Println("[DbMigration]execute sql file success:", migration.path)
	}
	return nil
}

func migrationLockName(dbName string) string {
	const prefix = "imserver:migration:"
	if len(prefix+dbName) <= 64 {
		return prefix + dbName
	}
	return fmt.Sprintf("%s%x", prefix, crc32.ChecksumIEEE([]byte(dbName)))
}

func acquireMigrationLock(ctx context.Context, conn *sql.Conn, name string, wait time.Duration) (bool, error) {
	var result sql.NullInt64
	if err := conn.QueryRowContext(ctx, "SELECT GET_LOCK(?, ?)", name, int(wait/time.Second)).Scan(&result); err != nil {
		return false, fmt.Errorf("acquire migration lock %q: %w", name, err)
	}
	if !result.Valid {
		return false, fmt.Errorf("acquire migration lock %q returned NULL", name)
	}
	return result.Int64 == 1, nil
}

func releaseMigrationLock(ctx context.Context, conn *sql.Conn, name string) error {
	var result sql.NullInt64
	if err := conn.QueryRowContext(ctx, "SELECT RELEASE_LOCK(?)", name).Scan(&result); err != nil {
		return fmt.Errorf("release migration lock %q: %w", name, err)
	}
	if !result.Valid || result.Int64 != 1 {
		return fmt.Errorf("release migration lock %q returned %v", name, result)
	}
	return nil
}

func readMigrationVersion(ctx context.Context, conn *sql.Conn) (int64, error) {
	var value string
	err := conn.QueryRowContext(ctx, "SELECT conf_value FROM globalconfs WHERE conf_key=?", JimDbVersionKey).Scan(&value)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, nil
	}
	if err != nil {
		var mysqlErr *mysqlDriver.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1146 {
			return 0, nil
		}
		return 0, fmt.Errorf("read migration version: %w", err)
	}
	version, err := tools.String2Int64(value)
	if err != nil || version <= 0 {
		return 0, fmt.Errorf("invalid migration version %q", value)
	}
	return version, nil
}

func writeMigrationVersion(ctx context.Context, conn *sql.Conn, version int64) error {
	_, err := conn.ExecContext(ctx, "INSERT INTO globalconfs (conf_key,conf_value) VALUES(?,?) ON DUPLICATE KEY UPDATE conf_value=VALUES(conf_value)", JimDbVersionKey, fmt.Sprintf("%d", version))
	return err
}

func listMigrationFiles() ([]migrationFile, error) {
	entries, err := sqlFs.ReadDir("sqls")
	if err != nil {
		return nil, fmt.Errorf("read migration directory: %w", err)
	}
	migrations := make([]migrationFile, 0, len(entries))
	seen := make(map[int64]bool)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if len(name) != len("20060102.sql") || !strings.HasSuffix(name, ".sql") {
			return nil, fmt.Errorf("invalid migration file name %q", name)
		}
		version, err := tools.String2Int64(strings.TrimSuffix(name, ".sql"))
		if err != nil || version <= 0 || seen[version] {
			return nil, fmt.Errorf("invalid or duplicate migration version in file %q", name)
		}
		seen[version] = true
		migrations = append(migrations, migrationFile{version, "sqls/" + name})
	}
	sort.Slice(migrations, func(i, j int) bool { return migrations[i].version < migrations[j].version })
	return migrations, nil
}

func executeSQLFile(ctx context.Context, conn *sql.Conn, fileName string) error {
	file, err := sqlFs.Open(fileName)
	if err != nil {
		return fmt.Errorf("open migration file %s: %w", fileName, err)
	}
	defer file.Close()
	return executeSQLReader(ctx, file, fileName, func(ctx context.Context, query string) error {
		_, err := conn.ExecContext(ctx, query)
		return err
	})
}

func executeSQLReader(ctx context.Context, reader io.Reader, fileName string, execute func(context.Context, string) error) error {
	scanner := bufio.NewScanner(reader)
	scanner.Buffer(make([]byte, 64*1024), migrationMaxLineSize)
	var query strings.Builder
	statement := 0
	run := func() error {
		text := strings.TrimSpace(query.String())
		query.Reset()
		if text == "" {
			return nil
		}
		statement++
		if err := ctx.Err(); err != nil {
			return err
		}
		if err := execute(ctx, text); err != nil {
			return fmt.Errorf("execute migration %s statement %d: %w", fileName, statement, err)
		}
		return nil
	}
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "--") {
			continue
		}
		query.WriteString(line)
		query.WriteByte(' ')
		if strings.HasSuffix(line, ";") {
			if err := run(); err != nil {
				return err
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scan migration file %s: %w", fileName, err)
	}
	return run()
}

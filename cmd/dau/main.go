// Command dau reports daily active users without starting or migrating the server.
package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v3"
)

type mysqlConfig struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Address  string `yaml:"address"`
	Name     string `yaml:"name"`
}

func main() {
	if err := run(context.Background(), os.Args[1:], os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string, output, stderr io.Writer) error {
	flags := flag.NewFlagSet("dau", flag.ContinueOnError)
	flags.SetOutput(stderr)
	configPath := flags.String("config", "", "配置文件路径（自动查找 conf/config.yml 或 launcher/conf/config.yml）")
	appKey := flags.String("app-key", "", "应用 App Key（必填）")
	days := flags.Int("days", 30, "统计天数，包含今天（1–3660）")
	if err := flags.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return nil
		}
		return err
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("不支持位置参数: %v", flags.Args())
	}
	*appKey = strings.TrimSpace(*appKey)
	if *appKey == "" {
		return fmt.Errorf("--app-key 不能为空")
	}
	if *days < 1 || *days > 3660 {
		return fmt.Errorf("--days 必须在 1 到 3660 之间")
	}
	config, path, err := loadConfig(*configPath)
	if err != nil {
		return err
	}
	dsn := mysql.NewConfig()
	dsn.User, dsn.Passwd = config.User, config.Password
	dsn.Net, dsn.Addr, dsn.DBName = "tcp", config.Address, config.Name
	dsn.ParseTime, dsn.Loc, dsn.Timeout = true, time.UTC, 5*time.Second
	db, err := sql.Open("mysql", dsn.FormatDSN())
	if err != nil {
		return fmt.Errorf("初始化 MySQL 连接失败: %w", err)
	}
	defer db.Close()
	db.SetMaxOpenConns(2)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(5 * time.Minute)
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	err = db.PingContext(pingCtx)
	cancel()
	if err != nil {
		return fmt.Errorf("连接 MySQL 失败（配置: %s）: %w", path, err)
	}
	return streamDailyActivities(ctx, output, *appKey, time.Now(), *days, func(ctx context.Context, date time.Time) (int64, error) {
		queryCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()
		return countDailyActivity(queryCtx, db, *appKey, date)
	})
}

func countDailyActivity(ctx context.Context, db *sql.DB, appKey string, date time.Time) (int64, error) {
	// New writes use milliseconds; legacy seconds may coexist for the same user.
	// Count each user once, and prefer the canonical millisecond summary over
	// its legacy equivalent when per-user detail has been cleaned up.
	var count int64
	if err := db.QueryRowContext(ctx, "SELECT COUNT(DISTINCT user_id) FROM useractivities FORCE INDEX (uniq_userid) WHERE app_key = ? AND time_mark IN (?, ?)", appKey, date.Unix(), date.UnixMilli()).Scan(&count); err != nil {
		return 0, err
	}
	if count > 0 {
		return count, nil
	}
	err := db.QueryRowContext(ctx, "SELECT COALESCE((SELECT count FROM dailyactivities WHERE app_key = ? AND time_mark = ?), (SELECT count FROM dailyactivities WHERE app_key = ? AND time_mark = ?), 0)", appKey, date.UnixMilli(), appKey, date.Unix()).Scan(&count)
	return count, err
}

func loadConfig(path string) (mysqlConfig, string, error) {
	if path == "" {
		for _, candidate := range []string{filepath.Join("conf", "config.yml"), filepath.Join("launcher", "conf", "config.yml")} {
			if _, err := os.Stat(candidate); err == nil {
				path = candidate
				break
			}
		}
	}
	if path == "" {
		return mysqlConfig{}, "", fmt.Errorf("找不到配置文件，请使用 --config 指定")
	}
	content, err := os.ReadFile(path)
	if err != nil {
		return mysqlConfig{}, path, fmt.Errorf("读取配置文件 %q: %w", path, err)
	}
	var config struct {
		MySQL mysqlConfig `yaml:"mysql"`
	}
	if err := yaml.Unmarshal(content, &config); err != nil {
		return mysqlConfig{}, path, fmt.Errorf("解析配置文件 %q: %w", path, err)
	}
	if config.MySQL.User == "" || config.MySQL.Address == "" || config.MySQL.Name == "" {
		return mysqlConfig{}, path, fmt.Errorf("配置文件 %q 必须包含 mysql.user、mysql.address 和 mysql.name", path)
	}
	return config.MySQL, path, nil
}

type dailyActivityQuery func(context.Context, time.Time) (int64, error)

func streamDailyActivities(ctx context.Context, output io.Writer, appKey string, now time.Time, days int, query dailyActivityQuery) error {
	if _, err := fmt.Fprintf(output, "%-10s  %-20s  %s\n", "DATE (UTC)", "APP_KEY", "DAU"); err != nil {
		return err
	}
	today := now.UTC().Truncate(24 * time.Hour)
	for offset := 0; offset < days; offset++ {
		if err := ctx.Err(); err != nil {
			return err
		}
		date := today.AddDate(0, 0, -offset)
		count, err := query(ctx, date)
		if err != nil {
			return fmt.Errorf("查询 %s 日活失败: %w", date.Format(time.DateOnly), err)
		}
		if _, err := fmt.Fprintf(output, "%-10s  %-20s  %d\n", date.Format(time.DateOnly), appKey, count); err != nil {
			return err
		}
	}
	return nil
}

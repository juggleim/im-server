package dbcommons

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"im-server/commons/configures"
	"im-server/commons/logs"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var db *gorm.DB

func GetDb() *gorm.DB {
	return db
}
func InitMysql() error {
	var err error

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		configures.Config.Mysql.User,
		configures.Config.Mysql.Password,
		configures.Config.Mysql.Address,
		configures.Config.Mysql.DbName)

	logLevel := logger.Silent
	if configures.Config.Mysql.Debug {
		logLevel = logger.Info
	}

	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
			NameReplacer:  nil,
			TablePrefix:   "",
		},
		Logger: &dbLogger{logLevel: logLevel},
	})

	if err != nil {
		log.Fatalf("connect mysql err: %v", err)
		return err
	}

	if configures.Config.Mysql.Debug {
		db = db.Debug()
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("get sql db err: %v", err)
		return err
	}
	sqlDB.SetMaxIdleConns(20)
	sqlDB.SetMaxOpenConns(500)
	sqlDB.SetConnMaxLifetime(time.Second * 9) // mysql连接默认10s断开
	return nil
}

// CloseDB closes database connection (unnecessary)
func CloseDB() {
	sqlDB, err := db.DB()
	if err != nil {
		return
	}
	sqlDB.Close()
}

type dbLogger struct {
	logLevel logger.LogLevel
}

func (l *dbLogger) LogMode(level logger.LogLevel) logger.Interface {
	return &dbLogger{logLevel: level}
}

func (l *dbLogger) Info(ctx context.Context, msg string, args ...interface{}) {
	logs.Debugf("GORM: "+msg, args...)
}

func (l *dbLogger) Warn(ctx context.Context, msg string, args ...interface{}) {
	logs.Debugf("GORM: "+msg, args...)
}

func (l *dbLogger) Error(ctx context.Context, msg string, args ...interface{}) {
	logs.Debugf("GORM: "+msg, args...)
}

func (l *dbLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	sql, rows := fc()
	if l.logLevel > logger.Silent {
		logs.Debugf("SQL:%v rows:%v err:%v", sql, rows, err)
	}
}

func Create(t interface{}) error {
	if err := db.Create(t).Error; err != nil {
		return err
	}
	return nil
}

func TxCreate(tx *gorm.DB, t interface{}) error {
	if err := tx.Create(t).Error; err != nil {
		return err
	}
	return nil
}

func UpdModelMapByConds(m interface{}, conds []*Condition, data map[string]interface{}) (int64, error) {
	where, params := GetWhere(conds)
	save := db.Model(m).Where(where, params...).Updates(data)
	return save.RowsAffected, save.Error
}

func TxUpdModelMapByConds(tx *gorm.DB, m interface{}, conds []*Condition, data map[string]interface{}) (int64, error) {
	where, params := GetWhere(conds)
	save := tx.Model(m).Where(where, params...).Updates(data)
	return save.RowsAffected, save.Error
}

type Condition struct {
	K string      `json:"k"`
	V interface{} `json:"v"`
	C string      `json:"cond"`
}

func GetWhere(c []*Condition) (string, []interface{}) {
	var wh []string
	var pa []interface{}
	for _, v := range c {
		re := "?"
		if cu := strings.ToLower(v.C); cu == "in" || cu == "not in" {
			re = "(?)"
		}
		wh = append(wh, fmt.Sprintf(`%s %s %s`, v.K, v.GetCond(), re))
		pa = append(pa, v.V)
	}
	return strings.Join(wh, " AND "), pa
}
func (c *Condition) GetCond() string {
	if c.C == "" {
		return "="
	} else {
		return c.C
	}
}

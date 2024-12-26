package dbcommons

import (
	"fmt"
	"log"
	"strings"
	"time"

	"im-server/commons/configures"
	"im-server/commons/logs"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

func GetDb() *gorm.DB {
	return db
}
func InitMysql() error {
	var err error

	db, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", //&interpolateParams=true
		configures.Config.Mysql.User,
		configures.Config.Mysql.Password,
		configures.Config.Mysql.Address,
		configures.Config.Mysql.DbName))

	if err != nil {
		log.Fatalf("connect mysql err: %v", err)
		return err
	}

	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return "" + defaultTableName
	}

	db.SingularTable(true)
	db.LogMode(configures.Config.Mysql.Debug)
	db.SetLogger(&dbLogger{})
	/*
		db.Callback().Create().Replace("gorm:update_time_stamp", updateTimeStampForCreateCallback)
		db.Callback().Update().Replace("gorm:update_time_stamp", updateTimeStampForUpdateCallback)
		db.Callback().Delete().Replace("gorm:delete", deleteCallback)
	*/
	db.DB().SetMaxIdleConns(20)
	db.DB().SetMaxOpenConns(500)
	db.DB().SetConnMaxLifetime(time.Second * 9) // mysql连接默认10s断开
	return nil
}

// CloseDB closes database connection (unnecessary)
func CloseDB() {
	defer db.Close()
}

type dbLogger struct {
}

func (l *dbLogger) Print(values ...interface{}) {
	logs.Debugf("SQL:%v", values...)
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
	save := db.Model(m).Where(where, params...).Update(data)
	return save.RowsAffected, save.Error
}

func TxUpdModelMapByConds(tx *gorm.DB, m interface{}, conds []*Condition, data map[string]interface{}) (int64, error) {
	where, params := GetWhere(conds)
	save := tx.Model(m).Where(where, params...).Update(data)
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

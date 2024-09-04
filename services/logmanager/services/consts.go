package services

import "fmt"

const (
	TableTmpl = "%s.%s.%s"
	TableMeta = "meta:%s.%s.%s"
	DbMeta    = "meta:%s.%s"
)

func tableName(appKey string, db string, tableName string) string {
	return fmt.Sprintf(TableTmpl, appKey, db, tableName)
}

func tableMetaName(appKey string, db string, tableName string) string {
	return fmt.Sprintf(TableMeta, appKey, db, tableName)
}

func dbMetaName(appKey string, db string) string {
	return fmt.Sprintf(DbMeta, appKey, db)
}

const (
	serverDb = "server"

	userConnectTable = "userconnect"
	connectTable     = "connect"
	businessTable    = "bus"
)

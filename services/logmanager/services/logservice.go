package services

import (
	"fmt"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"strings"
)

func WriteUserConnectLog(data *pbobjs.UserConnectLog) error {
	key := strings.Join([]string{data.AppKey, data.UserId, tools.Int642String(data.Timestamp)}, "_")
	return writeLog(data.AppKey, serverDb, userConnectTable, key, tools.ToJson(data))
}

func QryUserConnectLogs(appkey, userId string, start, count int64) ([]LogEntity, error) {
	prefix := fmt.Sprintf("%s_%s_", appkey, userId)
	startKey := ""
	if start > 0 {
		startKey = fmt.Sprintf("%s_%s_%d", appkey, userId, start)
	}
	return qryLogs(appkey, serverDb, userConnectTable, prefix, startKey, int(count))
}

func WriteConnectLog(data *pbobjs.ConnectionLog) error {
	key := strings.Join([]string{data.AppKey, data.Session, tools.Int642String(data.Timestamp)}, "_")
	return writeLog(data.AppKey, serverDb, connectTable, key, tools.ToJson(data))
}

func QryConnectLogs(appkey, session string, start, count int64) ([]LogEntity, error) {
	prefix := fmt.Sprintf("%s_%s_", appkey, session)
	startKey := ""
	if start > 0 {
		startKey = fmt.Sprintf("%s%d", prefix, start)
	}
	return qryLogs(appkey, serverDb, connectTable, prefix, startKey, int(count))
}

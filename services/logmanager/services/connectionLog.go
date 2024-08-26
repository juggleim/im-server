package services

import (
	"strconv"
)

func WriteConnectionLog(appKey string, userId string, timestamp int64, logContent string) (err error) {
	return writeLog(appKey, serverDb, connectionTable, connectionKey(userId, timestamp), logContent)
}

func connectionKey(userId string, timestamp int64) string {
	return userId + "_" + strconv.FormatInt(timestamp, 10)
}

func FetchConnectionLogs(appKey string, userId string, beginTime int64, count int) (logs []LogEntity, err error) {
	return fetchLogs(appKey, serverDb, connectionTable, appKey+"_"+userId+"_", connectionKey(userId, beginTime), count)
}

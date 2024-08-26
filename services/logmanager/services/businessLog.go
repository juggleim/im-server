package services

import "strconv"

func businessKey(session string, index uint32, timestamp int64) string {
	return session + "_" + strconv.FormatUint(uint64(index), 10) + "_" + strconv.FormatInt(timestamp, 10)
}

func WriteBusinessLog(appKey, session string, index uint32, timestamp int64, logContent string) (err error) {
	return writeLog(appKey, serverDb, businessTable, businessKey(session, index, timestamp), logContent)
}

func FetchBusinessLogs(appKey, session string, index uint32, beginTime int64, count int) (logs []LogEntity, err error) {
	return fetchLogs(appKey, serverDb, businessTable,
		session+"_"+strconv.FormatUint(uint64(index), 10)+"_",
		businessKey(session, index, beginTime), count)
}

package services

import "strconv"

func sdkKey(session string, index uint32, timestamp int64) string {
	return session + "_" + strconv.FormatUint(uint64(index), 10) + "_" + strconv.FormatInt(timestamp, 10)
}

func WriteSdkLog(appKey, session string, index uint32, timestamp int64, logContent string) (err error) {
	return writeLog(appKey, serverDb, sessionTable, sdkKey(session, index, timestamp), logContent)
}

func FetchSdkLogs(appKey, session string, index uint32, beginTime int64, count int) (logs []LogEntity, err error) {
	return fetchLogs(appKey, serverDb, sessionTable, appKey+"_"+session+"_"+strconv.FormatUint(uint64(index), 10)+"_", sdkKey(session, index, beginTime), count)
}

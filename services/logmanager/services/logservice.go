package services

import (
	"fmt"
	"im-server/commons/configures"
	"im-server/commons/kvdbcommons"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"sort"
	"strings"
	"time"
)

var timeFormat = "060102150405.000"

type LogEntity struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func QryUserConnectLogs(appkey, userId string, start, count int64) ([]LogEntity, error) {
	startTime := time.UnixMilli(start)

	candidateFiles := getLogFiles(start, false)
	resultLines := []string{}

	for _, candidateLogFile := range candidateFiles {
		fileHandler := tools.NewFileHandler()
		_, err := fileHandler.GreapWithFile(`"action":"connect"`, fmt.Sprintf("logs/%s", candidateLogFile))
		if err != nil {
			break
		}
		_, err = fileHandler.Greap(fmt.Sprintf(`"appkey":"%s"`, appkey))
		if err != nil {
			break
		}
		_, err = fileHandler.Greap(fmt.Sprintf(`"user_id":"%s"`, userId))
		if err != nil {
			break
		}
		_, err = fileHandler.Awk(fmt.Sprintf("$1 > %s", startTime.Format(timeFormat)))
		if err != nil {
			break
		}
		_, err = fileHandler.Head(int(count))
		if err != nil {
			break
		}
		lines := fileHandler.ResultLines()

		if len(lines) > 0 {
			resultLines = append(resultLines, lines...)
		}
		if len(resultLines) >= int(count) {
			break
		}
	}

	ret := []LogEntity{}
	for _, item := range resultLines {
		if len(item) > 16 {
			ret = append(ret, LogEntity{
				Key:   "",
				Value: appendTimestamp(item),
			})
		}
	}
	return ret, nil
}

func appendTimestamp(log string) string {
	if len(log) > 16 {
		timeStr := log[0:16]
		logStr := strings.TrimSpace(log[16:])
		tmpMap := make(map[string]interface{})
		err := tools.JsonUnMarshal([]byte(logStr), &tmpMap)
		if err == nil {
			t, e := time.ParseInLocation(timeFormat, timeStr, time.Local)
			if e == nil {
				tmpMap["timestamp"] = t.UnixMilli()
			} else {
				fmt.Println(e)
			}
			return tools.ToJson(tmpMap)
		}
	}
	return ""
}

func QryConnectLogs(appkey, session string, start, count int64) ([]LogEntity, error) {
	startTime := time.UnixMilli(start)

	candidateFiles := getLogFiles(start, false)
	resultLines := []string{}

	for _, candidateLogFile := range candidateFiles {
		fileHandler := tools.NewFileHandler()

		_, err := fileHandler.GreapWithFile(fmt.Sprintf(`"session":"%s"`, session), fmt.Sprintf("logs/%s", candidateLogFile))
		if err != nil {
			break
		}
		_, err = fileHandler.Greap(`"service_name":"connectmanager"`) //"action":
		if err != nil {
			break
		}
		_, err = fileHandler.Awk(fmt.Sprintf("$1 > %s", startTime.Format(timeFormat)))
		if err != nil {
			break
		}
		_, err = fileHandler.Head(int(count))
		if err != nil {
			break
		}

		lines := fileHandler.ResultLines()

		if len(lines) > 0 {
			resultLines = append(resultLines, lines...)
		}
		if len(resultLines) >= int(count) {
			break
		}
	}

	ret := []LogEntity{}
	for _, item := range resultLines {
		fmt.Println(item)
		if len(item) > 16 {
			ret = append(ret, LogEntity{
				Key:   "",
				Value: appendTimestamp(item),
			})
		}
	}
	return ret, nil
}

func QryBusinessLogs(appkey, session string, seqIndex int32, start, count int64) ([]LogEntity, error) {
	startTime := time.UnixMilli(start)

	candidateFiles := getLogFiles(start, true)
	resultLines := []string{}

	for _, candidateLogFile := range candidateFiles {
		fileHandler := tools.NewFileHandler()

		_, err := fileHandler.GreapWithFile(fmt.Sprintf(`"session":"%s"`, session), fmt.Sprintf("logs/%s", candidateLogFile))
		if err != nil {
			break
		}
		_, err = fileHandler.Greap(fmt.Sprintf(`"seq_index":%d,`, seqIndex)) //"action":
		if err != nil {
			break
		}
		_, err = fileHandler.Awk(fmt.Sprintf("$1 > %s", startTime.Format(timeFormat)))
		if err != nil {
			break
		}
		_, err = fileHandler.Head(int(count))
		if err != nil {
			break
		}

		lines := fileHandler.ResultLines()

		if len(lines) > 0 {
			resultLines = append(resultLines, lines...)
		}
		if len(resultLines) >= int(count) {
			break
		}
	}
	ret := []LogEntity{}
	for _, item := range resultLines {
		fmt.Println(item)
		if len(item) > 16 {
			ret = append(ret, LogEntity{
				Key:   "",
				Value: appendTimestamp(item),
			})
		}
	}
	return ret, nil
}

func getLogFiles(start int64, isEqual bool) []string {
	logName := configures.Config.Log.LogName
	logFiles := []string{}
	files := tools.ListDir("logs")
	for _, file := range files {
		if strings.HasPrefix(file, logName+".") && len(file) >= len(logName)+15 {
			logFiles = append(logFiles, file)
		}
	}
	if len(logFiles) > 0 {
		logArr := LogFilesArray(logFiles)
		sort.Sort(logArr)
		startLogFile := fmt.Sprintf("%s.%s.log", logName, time.UnixMilli(start).Format("2006010215"))
		retFile := []string{}
		for _, logFile := range logArr {
			compare := compareLogFile(startLogFile, logFile)
			if isEqual && compare == 0 {
				retFile = append(retFile, logFile)
			} else if !isEqual && compare <= 0 {
				retFile = append(retFile, logFile)
			}
		}
		return retFile
	}
	return []string{}
}

type LogFilesArray []string

func (s LogFilesArray) Len() int { return len(s) }

func (s LogFilesArray) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s LogFilesArray) Less(i, j int) bool {
	return compareLogFile(s[i], s[j]) < 0
}

func compareLogFile(l, r string) int {
	if len(l) == len(r) && l == r {
		return 0
	}

	lStr, lNo := parseLogFileName(l)
	rStr, rNo := parseLogFileName(r)

	if lStr != rStr {
		if lStr < rStr {
			return -1
		} else {
			return 1
		}
	} else {
		if lNo < rNo {
			return -1
		} else {
			return 1
		}
	}
}

func parseLogFileName(logFile string) (string, int) {
	index := strings.LastIndex(logFile, ".log")
	if index >= 0 {
		logStr := logFile[:index+4]
		logNoStr := logFile[index+4:]
		logNo := 0
		if logNoStr != "" {
			logNoStr = strings.TrimPrefix(logNoStr, ".")
			c, err := tools.String2Int64(logNoStr)
			if err == nil {
				logNo = int(c)
			}
		}
		return logStr, logNo
	}
	return "", 0
}

func WriteUserConnectLog(data *pbobjs.UserConnectLog) error {
	data.RealTime = data.Timestamp
	key := strings.Join([]string{string(ServerLogType_UserConnect), data.AppKey, data.UserId}, "_")
	_, err := kvdbcommons.TsAppend([]byte(key), []byte(tools.ToJson(data)))
	return err
}

func QryUserConnectLogs_bak(appkey, userId string, start, count int64) ([]LogEntity, error) {
	key := strings.Join([]string{string(ServerLogType_UserConnect), appkey, userId}, "_")
	items, err := kvdbcommons.TsScan([]byte(key), start, int(count))
	if err != nil {
		return []LogEntity{}, err
	}
	ret := []LogEntity{}
	for _, item := range items {
		var data pbobjs.UserConnectLog
		tools.JsonUnMarshal(item.Value, &data)
		data.Timestamp = item.Timestamp
		ret = append(ret, LogEntity{
			Key:   string(item.Key),
			Value: tools.ToJson(&data),
		})
	}
	return ret, nil
}

func WriteConnectLog(data *pbobjs.ConnectionLog) error {
	data.RealTime = data.Timestamp
	key := strings.Join([]string{string(ServerLogType_Connect), data.AppKey, data.Session}, "_")
	_, err := kvdbcommons.TsAppend([]byte(key), []byte(tools.ToJson(data)))
	return err
}

func QryConnectLogs_bak(appkey, session string, start, count int64) ([]LogEntity, error) {
	key := strings.Join([]string{string(ServerLogType_Connect), appkey, session}, "_")
	items, err := kvdbcommons.TsScan([]byte(key), start, int(count))
	if err != nil {
		return []LogEntity{}, err
	}
	ret := []LogEntity{}
	for _, item := range items {
		var data pbobjs.ConnectionLog
		tools.JsonUnMarshal(item.Value, &data)
		data.Timestamp = item.Timestamp
		ret = append(ret, LogEntity{
			Key:   string(item.Key),
			Value: tools.ToJson(&data),
		})
	}
	return ret, nil
}

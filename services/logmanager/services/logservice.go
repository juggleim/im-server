package services

import (
	"im-server/commons/kvdbcommons"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"strings"
)

type LogEntity struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func WriteUserConnectLog(data *pbobjs.UserConnectLog) error {
	data.RealTime = data.Timestamp
	key := strings.Join([]string{string(ServerLogType_UserConnect), data.AppKey, data.UserId}, "_")
	_, err := kvdbcommons.TsAppend([]byte(key), []byte(tools.ToJson(data)))
	return err
}

func QryUserConnectLogs(appkey, userId string, start, count int64) ([]LogEntity, error) {
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

func QryConnectLogs(appkey, session string, start, count int64) ([]LogEntity, error) {
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

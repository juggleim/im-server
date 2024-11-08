package services

import (
	"fmt"
	"im-server/services/commonservices"
	"im-server/services/connectmanager/server/imcontext"
	"strings"
	"time"
)

var connectStatisTimer *time.Ticker

func init() {
	//start concurrent connect statis
	startConnectStatis()
}

func startConnectStatis() {
	if connectStatisTimer != nil {
		connectStatisTimer.Stop()
	}
	connectStatisTimer = time.NewTicker(30 * time.Second)
	go func() {
		for task := range connectStatisTimer.C {
			current := time.Now().UnixMilli()
			if current-task.UnixMilli() > 500 {
				continue
			}
			foreachConnect()
		}
	}()
}

func foreachConnect() {
	connectCountMap := map[string]int64{}
	OnlineUserConnectMap.Range(func(key, value any) bool {
		identifier := key.(string)
		if len(identifier) > 0 {
			index := strings.Index(identifier, "_")
			if index > 0 {
				appkey := identifier[:index]
				ctxMap := value.(map[string]imcontext.WsHandleContext)
				c := len(ctxMap)
				if count, exist := connectCountMap[appkey]; exist {
					connectCountMap[appkey] = count + int64(c)
				} else {
					connectCountMap[appkey] = int64(c)
				}
			}
		}
		return true
	})
	for appkey, count := range connectCountMap {
		fmt.Println(appkey, "\t", count)
		commonservices.ReportConcurrentConnectCount(appkey, count)
	}
}

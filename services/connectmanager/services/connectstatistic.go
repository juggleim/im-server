package services

import (
	"im-server/services/commonservices"
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
	foreachAppConnectCount(func(appkey string, count int64) {
		commonservices.ReportConcurrentConnectCount(appkey, count)
	})
}

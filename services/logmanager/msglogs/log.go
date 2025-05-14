package msglogs

import (
	"fmt"
	"im-server/commons/configures"
	"im-server/commons/tools"
	"im-server/services/logmanager/msglogs/lumberjack"
	"path/filepath"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

func NewMsgLogger(appkey string) *zerolog.Logger {
	dir := configures.Config.MsgLogs.LogPath
	if dir == "" {
		dir = filepath.Join(configures.Config.Log.LogPath, "msglogs", appkey)
	} else {
		dir = filepath.Join(dir, "msglogs", appkey)
	}
	err := tools.CreateDirs(dir)
	if err != nil {
		fmt.Println("create msglogs path failed:", err)
		return nil
	}
	maxBackups := 24
	if configures.Config.MsgLogs.MaxBackups > 0 {
		maxBackups = configures.Config.MsgLogs.MaxBackups
	}
	logFile := &lumberjack.Logger{
		Filename:         filepath.Join(dir, "messages.log"),
		MaxSize:          0,
		MaxBackups:       maxBackups,
		Compress:         configures.Config.MsgLogs.IsCompress,
		LocalTime:        true,
		BackupTimeFormat: "2006010215",
	}
	hook := &HourlyRotationHook{
		TimeFormat: "060102150405.000",
		Logger:     logFile,
		LastHour:   time.Now().Hour(),
	}
	logger := zerolog.New(logFile).Hook(hook).With().Logger()
	return &logger
}

// HourlyRotationHook 按小时切割日志
type HourlyRotationHook struct {
	TimeFormat string
	Logger     *lumberjack.Logger
	LastHour   int
	Mutex      sync.Mutex
}

// Run 实现 Hook 接口
func (h *HourlyRotationHook) Run(e *zerolog.Event, level zerolog.Level, message string) {
	now := time.Now()
	currentHour := now.Hour()

	h.Mutex.Lock()
	defer h.Mutex.Unlock()

	// 如果小时数变化，切割日志
	if currentHour != h.LastHour {
		h.Logger.Rotate()
		h.LastHour = currentHour
	}
	timestamp := now.Local().Format(h.TimeFormat)
	h.Logger.Write([]byte(timestamp))
	h.Logger.Write([]byte{'\t'})
	h.Logger.Write([]byte(message))
}

package msglogs

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/caches"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/commonservices"
	"im-server/services/commonservices/msgdefines"
	"time"

	"github.com/rs/zerolog"
)

var msgLoggerCache *caches.LruCache
var msgLoggerLocks *tools.SegmentatedLocks

func init() {
	msgLoggerCache = caches.NewLruCacheWithReadTimeout("msglogger_cache", 1000, nil, time.Hour)
	msgLoggerLocks = tools.NewSegmentatedLocks(32)
}

func getMsgLogger(appkey string) *zerolog.Logger {
	if val, exist := msgLoggerCache.Get(appkey); exist {
		return val.(*zerolog.Logger)
	} else {
		l := msgLoggerLocks.GetLocks(appkey)
		l.Lock()
		defer l.Unlock()
		if val, exist := msgLoggerCache.Get(appkey); exist {
			return val.(*zerolog.Logger)
		} else {
			msgLogger := NewMsgLogger(appkey)
			if msgLogger != nil {
				msgLoggerCache.Add(appkey, msgLogger)
			}
			return msgLogger
		}
	}
}

func LogMsg(ctx context.Context, downMsg *pbobjs.DownMsg) {
	isState := msgdefines.IsStateMsg(downMsg.Flags)
	isCmd := msgdefines.IsCmdMsg(downMsg.Flags)
	if isState || isCmd {
		return
	}
	appkey := bases.GetAppKeyFromCtx(ctx)
	appinfo, exist := commonservices.GetAppInfo(appkey)
	if exist && appinfo != nil && appinfo.RecordMsgLogs {
		msgLogger := getMsgLogger(appkey)
		if msgLogger != nil {
			platform := bases.GetPlatformFromCtx(ctx)
			msgLogger.Info().Str("platform", platform).
				Str("sender", downMsg.SenderId).
				Str("receiver", downMsg.TargetId).
				Int("channel_type", int(downMsg.ChannelType)).
				Str("msg_id", downMsg.MsgId).
				Int64("msg_time", downMsg.MsgTime).
				Str("msg_type", downMsg.MsgType).
				Str("msg_content", string(downMsg.MsgContent)).Msg("")
		}
	}
}

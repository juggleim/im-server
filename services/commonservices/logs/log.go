package logs

import (
	"context"
	"fmt"
	"im-server/commons/bases"
	"im-server/commons/logs"
	"im-server/commons/tools"
)

var (
	openLog bool = true
)

type LogEntity struct {
	kvPairs map[string]interface{}
}

func NewLogEntity() *LogEntity {
	return &LogEntity{
		kvPairs: make(map[string]interface{}),
	}
}

func WithContext(ctx context.Context) *LogEntity {
	log := &LogEntity{
		kvPairs: make(map[string]interface{}),
	}
	if openLog {
		//handle service tags
		tags := bases.GetTagsFromCtx(ctx)
		if len(tags) > 0 {
			for k, v := range tags {
				log.kvPairs[k] = v
			}
		}
		//handle ctx
		log.kvPairs["session"] = bases.GetSessionFromCtx(ctx)
		log.kvPairs["method"] = bases.GetMethodFromCtx(ctx)
		log.kvPairs["expend"] = bases.GetExpendFromCtx(ctx)
		log.kvPairs["seq_index"] = bases.GetSeqIndexFromCtx(ctx)
	}
	return log
}

func (log *LogEntity) WithField(key string, value interface{}) *LogEntity {
	if openLog {
		log.kvPairs[key] = value
	}
	return log
}

func (log *LogEntity) WithFields(kvPairs map[string]interface{}) *LogEntity {
	if openLog {
		for k, v := range kvPairs {
			log.kvPairs[k] = v
		}
	}
	return log
}

func (log *LogEntity) Errorf(format string, v ...interface{}) {
	if openLog {
		if format != "" {
			log.kvPairs["msg"] = fmt.Sprintf(format, v...)
		}
		logs.Error(tools.ToJson(log.kvPairs))
	}
}

func (log *LogEntity) Error(errMsg string) {
	log.Errorf(errMsg)
}

func (log *LogEntity) Warnf(format string, v ...interface{}) {
	if openLog {
		if format != "" {
			log.kvPairs["msg"] = fmt.Sprintf(format, v...)
		}
		logs.Warn(tools.ToJson(log.kvPairs))
	}
}

func (log *LogEntity) Warn(warnMsg string) {
	log.Warnf(warnMsg)
}

func (log *LogEntity) Infof(format string, v ...interface{}) {
	if openLog {
		if format != "" {
			log.kvPairs["msg"] = fmt.Sprintf(format, v...)
		}
		logs.Info(tools.ToJson(log.kvPairs))
	}
}

func (log *LogEntity) Info(infoMsg string) {
	log.Infof(infoMsg)
}

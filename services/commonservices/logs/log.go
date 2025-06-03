package logs

import (
	"context"
	"fmt"
	"im-server/commons/bases"
	"im-server/commons/logs"
	"strings"
)

var (
	openLog bool = true
)

type LogEntity struct {
	fields []string
}

func NewLogEntity() *LogEntity {
	return &LogEntity{
		fields: []string{},
	}
}

func WithContext(ctx context.Context) *LogEntity {
	log := &LogEntity{
		fields: []string{},
	}
	if openLog {
		//handle service tags
		tags := bases.GetTagsFromCtx(ctx)
		if len(tags) > 0 {
			for k, v := range tags {
				log.fields = append(log.fields, fmt.Sprintf("%s:%s", k, v))
			}
		}
		//handle ctx
		log.fields = append(log.fields, fmt.Sprintf("session:%v", bases.GetSessionFromCtx(ctx)))
		log.fields = append(log.fields, fmt.Sprintf("method:%v", bases.GetMethodFromCtx(ctx)))
		log.fields = append(log.fields, fmt.Sprintf("expend:%v", bases.GetExpendFromCtx(ctx)))
		log.fields = append(log.fields, fmt.Sprintf("seq_index:%v", bases.GetSeqIndexFromCtx(ctx)))
	}
	return log
}

func (log *LogEntity) WithField(key string, value interface{}) *LogEntity {
	if openLog {
		log.fields = append(log.fields, fmt.Sprintf("%s:%v", key, value))
	}
	return log
}

func (log *LogEntity) Errorf(format string, v ...interface{}) {
	if openLog {
		initFormat := strings.Join(log.fields, "\t")
		logs.Errorf(initFormat+"\t"+format, v...)
	}
}

func (log *LogEntity) Error(errMsg string) {
	log.Errorf(errMsg)
}

func (log *LogEntity) Warnf(format string, v ...interface{}) {
	if openLog {
		initFormat := strings.Join(log.fields, "\t")
		logs.Warnf(initFormat+"\t"+format, v...)
	}
}

func (log *LogEntity) Warn(warnMsg string) {
	log.Warnf(warnMsg)
}

func (log *LogEntity) Infof(format string, v ...interface{}) {
	if openLog {
		initFormat := strings.Join(log.fields, "\t")
		logs.Infof(initFormat+"\t"+format, v...)
	}
}

func (log *LogEntity) Info(infoMsg string) {
	log.Infof(infoMsg)
}

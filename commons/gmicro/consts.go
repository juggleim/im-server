package gmicro

import (
	"im-server/commons/gmicro/logs"

	"github.com/sirupsen/logrus"
)

func SetLogger(logger *logrus.Logger) {
	if logger != nil {
		logs.Logger = logger
	}
}

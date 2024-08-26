package logs

import "github.com/sirupsen/logrus"

var Logger *logrus.Logger

func getLogger() *logrus.Logger {
	if Logger == nil {
		Logger = logrus.StandardLogger()
	}
	return Logger
}

func Panic(f interface{}, v ...interface{}) {
	getLogger().Panic(f, v)
}

func Fata(f interface{}, v ...interface{}) {
	getLogger().Fatal(f, v)
}

func Error(f interface{}, v ...interface{}) {
	getLogger().Error(f, v)
}
func Warn(f interface{}, v ...interface{}) {
	getLogger().Warn(f, v)
}

func Info(f interface{}, v ...interface{}) {
	getLogger().Info(f, v)
}

func Debug(f interface{}, v ...interface{}) {
	getLogger().Debug(f, v)
}

func Trace(f interface{}, v ...interface{}) {
	getLogger().Trace(f, v)
}

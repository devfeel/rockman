package logger

import "github.com/devfeel/dotlog"

const (
	defaultDateFormatForFileName = "2006_01_02"
	defaultFullTimeLayout        = "2006-01-02 15:04:05.999999"
)

const (
	LoggerName_Service    = "ServiceLogger"
	LoggerName_Repository = "RepositoryLogger"
	LoggerName_Node       = "NodeLogger"
	LoggerName_Runtime    = "RuntimeLogger"
	LoggerName_Default    = "DefaultLogger"
)

type Logger interface {
	dotlog.Logger
}

func StartLogService(confPath string) error {
	return dotlog.StartLogService(confPath + "/dotlog.conf")
}

func GetLogger(loggerName string) Logger {
	return dotlog.GetLogger(loggerName)
}

func Default() Logger {
	return GetLogger(LoggerName_Default)
}

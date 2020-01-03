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
)

type Logger interface {
	dotlog.Logger
}

func StartLogService(confPath string) error {
	return dotlog.StartLogService(confPath + "/log.conf")
}

func GetLogger(loggerName string) dotlog.Logger {
	return dotlog.GetLogger(loggerName)
}

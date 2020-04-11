package service

import (
	"github.com/devfeel/cache"
	"github.com/devfeel/dotlog"
	"github.com/devfeel/rockman/logger"
)

type BaseService struct {
	RedisCache cache.RedisCache
}

func GetLogger() dotlog.Logger {
	return logger.GetLogger(logger.LoggerName_Service)
}

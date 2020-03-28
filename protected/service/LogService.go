package service

import (
	"github.com/devfeel/rockman/logger"
	"github.com/devfeel/rockman/protected/model"
)

type LogService struct {
	BaseService
}

func NewLogService() *LogService {
	service := &LogService{}
	return service
}

func (service *LogService) WriteExecLog(log *model.TaskExecLog) error {
	logger.Service().DebugS("CreateExecLog", *log)
	return nil
}

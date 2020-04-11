package service

import (
	"github.com/devfeel/rockman/protected/model"
	"github.com/devfeel/rockman/protected/repository"
	"time"
)

type LogService struct {
	BaseService
	executorRepo *repository.ExecutorRepo
}

func NewLogService() *LogService {
	service := &LogService{
		executorRepo: repository.GetExecutorRepo(),
	}
	return service
}

// WriteExecLog
func (service *LogService) WriteExecLog(log *model.TaskExecLog) error {
	log.CreateTime = time.Now()
	_, err := service.executorRepo.WriteExecLog(log)
	return err
}

// QueryExecLogs
func (service *LogService) QueryExecLogs(taskId string, pageReq *model.PageRequest) (*model.PageResult, error) {
	result, err := service.executorRepo.QueryExecLogs(taskId, pageReq)
	return result, err
}

// WriteNodeTraceLog
func (service *LogService) WriteNodeTraceLog(log *model.NodeTraceLog) error {
	log.CreateTime = time.Now()
	_, err := service.executorRepo.WriteNodeTraceLog(log)
	return err
}

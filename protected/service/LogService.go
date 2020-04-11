package service

import (
	"time"

	"github.com/devfeel/rockman/protected/model"
	"github.com/devfeel/rockman/protected/repository"
	"github.com/devfeel/rockman/protected/viewmodel"
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
func (service *LogService) QueryExecLogs(qc *viewmodel.TaskExecLogQC) (*model.PageResult, error) {
	result, err := service.executorRepo.QueryExecLogs(qc)
	return result, err
}

// WriteNodeTraceLog
func (service *LogService) WriteNodeTraceLog(log *model.NodeTraceLog) error {
	log.CreateTime = time.Now()
	_, err := service.executorRepo.WriteNodeTraceLog(log)
	return err
}

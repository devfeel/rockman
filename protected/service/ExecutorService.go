package service

import (
	"github.com/devfeel/rockman/logger"
	"github.com/devfeel/rockman/protected/model"
	"github.com/devfeel/rockman/protected/repository"
	"time"
)

var (
	defaultLogger logger.Logger
)

type ExecutorService struct {
	BaseService
	executorRepository *repository.ExecutorRepository
}

func init() {
	defaultLogger = logger.GetLogger(logger.LoggerName_Service)
}

func NewExecutorService() *ExecutorService {
	service := &ExecutorService{
		executorRepository: repository.GetTaskRepository(),
	}
	return service
}

// QueryExecutors
func (service *ExecutorService) QueryExecutors(nodeId string, pageReq *model.PageRequest) (*model.PageResult, error) {
	result, err := service.executorRepository.QueryExecutors(nodeId, pageReq)
	return result, err
}

// WriteExecLog
func (service *ExecutorService) WriteExecLog(log *model.TaskExecLog) error {
	log.CreateTime = time.Now()
	_, err := service.executorRepository.WriteExecLog(log)
	return err
}

// QueryExecLogs
func (service *ExecutorService) QueryExecLogs(taskId string, pageReq *model.PageRequest) (*model.PageResult, error) {
	result, err := service.executorRepository.QueryExecLogs(taskId, pageReq)
	return result, err
}

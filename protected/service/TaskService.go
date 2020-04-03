package service

import (
	"github.com/devfeel/rockman/logger"
	"github.com/devfeel/rockman/protected/model"
	repository2 "github.com/devfeel/rockman/protected/repository"
	"time"
)

var (
	defaultLogger logger.Logger
)

type TaskService struct {
	BaseService
	taskRepository *repository2.TaskRepository
}

func init() {
	defaultLogger = logger.GetLogger(logger.LoggerName_Service)
}

func NewTaskService() *TaskService {
	service := &TaskService{
		taskRepository: repository2.GetTaskRepository(),
	}
	return service
}

// QueryTasks
func (service *TaskService) QueryTasks() ([]*model.TaskInfo, error) {
	var results []*model.TaskInfo
	var err error
	err = service.taskRepository.QueryTasks(&results)
	return results, err
}

// WriteExecLog
func (service *TaskService) WriteExecLog(log *model.TaskExecLog) error {
	log.CreateTime = time.Now()
	_, err := service.taskRepository.WriteExecLog(log)
	return err
}

func (service *TaskService) QueryLogs() ([]*model.TaskExecLog, error) {
	var results []*model.TaskExecLog
	var err error
	err = service.taskRepository.QueryLogs(&results)
	return results, err
}

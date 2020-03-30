package service

import (
	"errors"
	"github.com/devfeel/rockman/logger"
	"github.com/devfeel/rockman/protected/model"
	repository2 "github.com/devfeel/rockman/protected/repository"
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

// QueryTasksByNodeID
func (service *TaskService) QueryTasks() ([]*model.TaskInfo, error) {
	var results []*model.TaskInfo
	var err error
	err = service.taskRepository.QueryTasks(&results)
	if err == nil {
		if len(results) <= 0 {
			results = nil
			err = errors.New("not exists task info")
		}
	}
	return results, err
}

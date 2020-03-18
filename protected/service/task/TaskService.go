package task

import (
	"errors"
	"github.com/devfeel/rockman/logger"
	"github.com/devfeel/rockman/protected/model"
	task2 "github.com/devfeel/rockman/protected/repository/task"
	"github.com/devfeel/rockman/protected/service"
)

var (
	defaultLogger logger.Logger
)

type TaskService struct {
	service.BaseService
	taskRepository *task2.TaskRepository
}

func init() {
	defaultLogger = logger.GetLogger(logger.LoggerName_Service)
}

func NewTaskService() *TaskService {
	service := &TaskService{
		taskRepository: task2.GetTaskRepository(),
	}
	return service
}

// QueryTasksByNodeID
func (service *TaskService) QueryTasksByNodeID(nodeID string) ([]*model.TaskInfo, error) {
	if nodeID == "" {
		return nil, errors.New("must set NodeID")
	}
	var results []*model.TaskInfo
	var err error
	err = service.taskRepository.QueryTasksByNodeID(&results, nodeID)
	if err == nil {
		if len(results) <= 0 {
			results = nil
			err = errors.New("not exists task info")
		}
	}
	return results, err
}

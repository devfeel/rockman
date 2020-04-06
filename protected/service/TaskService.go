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
	result, err := service.taskRepository.QueryTasks()
	return result, err
}

// WriteExecLog
func (service *TaskService) WriteExecLog(log *model.TaskExecLog) error {
	log.CreateTime = time.Now()
	_, err := service.taskRepository.WriteExecLog(log)
	return err
}

// QueryExecLogs
func (service *TaskService) QueryExecLogs(taskId string, pageReq *model.PageRequest) (*model.PageResult, error) {
	result, err := service.taskRepository.QueryExecLogs(taskId, pageReq)
	return result, err
}

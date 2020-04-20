package service

import (
	"github.com/devfeel/rockman/protected/model"
	"github.com/devfeel/rockman/protected/repository"
	"time"
)

type LogService struct {
	BaseService
	repo        *repository.LogRepo
	execService *ExecutorService
}

func NewLogService() *LogService {
	service := &LogService{
		repo:        repository.NewLogRepo(),
		execService: NewExecutorService(),
	}
	return service
}

// WriteExecLog
func (service *LogService) WriteExecLog(log *model.TaskExecLog) error {
	log.CreateTime = time.Now()
	_, err := service.repo.WriteExecLog(log)
	return err
}

// WriteNodeTraceLog
func (service *LogService) WriteNodeTraceLog(log *model.NodeTraceLog) error {
	log.CreateTime = time.Now()
	_, err := service.repo.WriteNodeTraceLog(log)
	return err
}

// WriteSubmitLog
func (service *LogService) WriteSubmitLog(log *model.TaskSubmitLog) error {
	log.CreateTime = time.Now()
	_, err := service.repo.WriteSubmitLog(log)
	if err != nil {
		return err
	}
	if log.IsSuccess {
		runInfo := &model.ExecutorRunInfo{
			TaskID:       log.TaskID,
			NodeID:       log.NodeID,
			NodeEndPoint: log.NodeEndPoint,
		}
		service.execService.SetExecutorRunInfo(runInfo)
	}
	return nil
}

// QueryExecLogs
func (service *LogService) QueryExecLogs(taskId string, pageReq *model.PageRequest) (*model.PageResult, error) {
	result, err := service.repo.QueryExecLogs(taskId, pageReq)
	return result, err
}

// QueryStateLog
func (service *LogService) QueryStateLog(taskId string, pageReq *model.PageRequest) (*model.PageResult, error) {
	result, err := service.repo.QueryStateLog(taskId, pageReq)
	return result, err
}

// QueryTaskSubmitLog
func (service *LogService) QueryTaskSubmitLog(taskId string, pageReq *model.PageRequest) (*model.PageResult, error) {
	result, err := service.repo.QueryTaskSubmitLog(taskId, pageReq)
	return result, err
}

// QueryNodeTraceLog
func (service *LogService) QueryNodeTraceLog(nodeId string, pageReq *model.PageRequest) (*model.PageResult, error) {
	result, err := service.repo.QueryNodeTraceLog(nodeId, pageReq)
	return result, err
}

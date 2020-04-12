package service

import (
	"time"

	"github.com/devfeel/rockman/protected/model"
	"github.com/devfeel/rockman/protected/repository"
	"github.com/devfeel/rockman/protected/viewmodel"
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

// QueryExecLogs
func (service *LogService) QueryExecLogs(qc *viewmodel.TaskExecLogQC) (*model.PageResult, error) {
	result, err := service.repo.QueryExecLogs(qc)
	return result, err
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
			NodeID:       log.TaskID,
			NodeEndPoint: log.NodeEndPoint,
		}
		service.execService.SetExecutorRunInfo(runInfo)
	}
	return nil
}

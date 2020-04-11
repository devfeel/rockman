package service

import (
	"github.com/devfeel/rockman/protected/model"
	"github.com/devfeel/rockman/protected/repository"
	"time"
)

type LogService struct {
	BaseService
	repo *repository.LogRepo
}

func NewLogService() *LogService {
	service := &LogService{
		repo: repository.NewLogRepo(),
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
func (service *LogService) QueryExecLogs(taskId string, pageReq *model.PageRequest) (*model.PageResult, error) {
	result, err := service.repo.QueryExecLogs(taskId, pageReq)
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
	return err
}

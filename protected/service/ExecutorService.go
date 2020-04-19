package service

import (
	"database/sql"
	"github.com/devfeel/rockman/core"
	"github.com/devfeel/rockman/logger"
	"github.com/devfeel/rockman/protected/model"
	"github.com/devfeel/rockman/protected/repository"
	runtime "github.com/devfeel/rockman/runtime/executor"
	"strings"
	"sync"
	"time"
)

type ExecutorService struct {
	BaseService
	repo         *repository.ExecutorRepo
	updateLocker *sync.Mutex
}

func NewExecutorService() *ExecutorService {
	service := &ExecutorService{
		repo:         repository.NewExecutorRepo(),
		updateLocker: new(sync.Mutex),
	}
	return service
}

// AddExecutor
func (service *ExecutorService) AddExecutor(model *model.ExecutorInfo) *core.Result {
	result := validateExecutorInfo(model)
	if !result.IsSuccess() {
		return result
	}
	data, err := service.QueryExecutorByTaskId(model.TaskID)
	if err != nil {
		return core.FailedResult(-3001, "QueryExecutor error: "+err.Error())
	}
	if data != nil {
		return core.FailedResult(-2101, "already exists this TaskID["+model.TaskID+"]")
	}

	err = service.repo.InsertOnce(model)
	if err != nil {
		return core.FailedResult(-3002, "InsertOnce error: "+err.Error())
	} else {
		return core.SuccessResult()
	}
}

// UpdateExecutor
func (service *ExecutorService) UpdateExecutor(model *model.ExecutorInfo) *core.Result {
	result := validateExecutorInfo(model)
	if !result.IsSuccess() {
		return result
	}
	data, err := service.QueryExecutorByTaskId(model.TaskID)
	if err != nil {
		return core.FailedResult(-3001, "QueryExecutor error: "+err.Error())
	}
	if data == nil {
		return core.FailedResult(-2101, "not exists this TaskID["+model.TaskID+"]")
	}
	err = service.repo.UpdateOnce(model)
	if err != nil {
		return core.FailedResult(-3001, "UpdateOnce error: "+err.Error())
	} else {
		//TODO remove executor to leader node
		//TODO submit executor to leader node
		return core.SuccessResult()
	}
}

// SetExecutorRunInfo
// if exists, update it
// if not exists, insert it
func (service *ExecutorService) SetExecutorRunInfo(model *model.ExecutorRunInfo) error {
	service.updateLocker.Lock()
	defer service.updateLocker.Unlock()
	_, err := service.repo.QueryRunInfo(model.TaskID)
	if err != nil && err != sql.ErrNoRows {
		logger.Service().Error(err, "SetExecutorRunInfo["+model.TaskID+"] error")
		return err
	}
	if err == sql.ErrNoRows {
		model.LastUpdateTime = time.Now()
		model.CreateTime = time.Now()
		return service.repo.InsertRunInfo(model)
	} else {
		model.LastUpdateTime = time.Now()
		return service.repo.UpdateRunInfo(model)
	}
}

// RemoveExecutor
func (service *ExecutorService) RemoveExecutor(id int64) error {
	// TODO check data
	// TODO remove executor to leader node
	// TODO remove log?
	return service.repo.DeleteOnce(id)
}

// QueryExecutorById
func (service *ExecutorService) QueryExecutorById(id int64) (*model.ExecutorInfo, error) {
	return service.repo.GetExecutorById(id)
}

// QueryExecutorByTaskId
func (service *ExecutorService) QueryExecutorByTaskId(taskId string) (*model.ExecutorInfo, error) {
	return service.repo.GetExecutorByTaskId(taskId)
}

// QueryExecutors
func (service *ExecutorService) QueryExecutors(pageReq *model.PageRequest) (*model.PageResult, error) {
	result, err := service.repo.QueryExecutors(pageReq)
	return result, err
}

// QueryAllExecutors
func (service *ExecutorService) QueryAllExecutors() ([]*model.ExecutorInfo, error) {
	result, err := service.repo.QueryAllExecutors()
	return result, err
}

// validateExecutorInfo
func validateExecutorInfo(model *model.ExecutorInfo) *core.Result {
	if model == nil {
		return core.FailedResult(-2000, "executor info is nil")
	}
	if model.TaskID == "" {
		return core.FailedResult(-2001, "TaskID is empty")
	}
	if len(model.TaskID) > 64 {
		return core.FailedResult(-2002, "TaskID is more than 60 characters")
	}
	if len(model.TargetType) == 0 {
		return core.FailedResult(-2003, "TargetType is empty")
	}
	if model.TaskType == "" {
		return core.FailedResult(-2004, "TaskType is empty")
	}
	model.TaskType = strings.ToLower(model.TaskType)
	if model.TaskType != "cron" && model.TaskType != "loop" {
		return core.FailedResult(-2005, "TaskType is not match")
	}
	if model.TaskType == "cron" && model.Express == "" {
		return core.FailedResult(-2006, "Express is empty")
	}
	if model.TargetType == "" {
		return core.FailedResult(-2007, "TargetType is empty")
	}
	model.TargetType = strings.ToLower(model.TargetType)
	if !runtime.ValidateTargetType(model.TargetType) {
		return core.FailedResult(-2008, "TargetType is not match")
	}
	if model.TargetConfig == "" {
		return core.FailedResult(-2009, "TargetConfig is empty")
	}
	if len(model.Remark) > 100 {
		return core.FailedResult(-2010, "Remark is more than 100 characters")
	}

	if model.TargetType == runtime.TargetType_Http {
		conf := new(runtime.HttpConfig)
		err := conf.LoadFromJson(model.TargetConfig)
		if err != nil {
			return core.FailedResult(-2011, "convert http config failed: "+err.Error())
		} else {
			model.RealTargetConfig = conf
		}
	}
	if model.TargetType == runtime.TargetType_GoSo {
		conf := new(runtime.GoConfig)
		err := conf.LoadFromJson(model.TargetConfig)
		if err != nil {
			return core.FailedResult(-2011, "convert go so config failed: "+err.Error())
		} else {
			model.RealTargetConfig = conf
		}
	}
	if model.TargetType == runtime.TargetType_Shell {
		conf := new(runtime.ShellConfig)
		err := conf.LoadFromJson(model.TargetConfig)
		if err != nil {
			return core.FailedResult(-2011, "convert shell config failed: "+err.Error())
		} else {
			model.RealTargetConfig = conf
		}
	}

	return core.SuccessResult()
}

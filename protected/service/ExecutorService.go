package service

import (
	"github.com/devfeel/rockman/core"
	"github.com/devfeel/rockman/logger"
	"github.com/devfeel/rockman/protected/model"
	"github.com/devfeel/rockman/protected/repository/executor"
	runtime "github.com/devfeel/rockman/runtime/executor"
	"strings"
	"time"
)

var (
	defaultLogger logger.Logger
)

type ExecutorService struct {
	BaseService
	executorRepository *executor.ExecutorRepository
}

func init() {
	defaultLogger = logger.GetLogger(logger.LoggerName_Service)
}

func NewExecutorService() *ExecutorService {
	service := &ExecutorService{
		executorRepository: executor.GetRepository(),
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

	err = service.executorRepository.InsertOnce(model)
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
	err = service.executorRepository.UpdateOnce(model)
	if err != nil {
		return core.FailedResult(-3001, "UpdateOnce error: "+err.Error())
	} else {
		//TODO remove executor to leader node
		//TODO submit executor to leader node
		return core.SuccessResult()
	}
}

// RemoveExecutor
func (service *ExecutorService) RemoveExecutor(id int64) error {
	// TODO check data
	// TODO remove executor to leader node
	// TODO remove log?
	return service.executorRepository.DeleteOnce(id)
}

// QueryExecutorById
func (service *ExecutorService) QueryExecutorById(id int64) (*model.ExecutorInfo, error) {
	return service.executorRepository.GetExecutorById(id)
}

// QueryExecutorByTaskId
func (service *ExecutorService) QueryExecutorByTaskId(taskId string) (*model.ExecutorInfo, error) {
	return service.executorRepository.GetExecutorByTaskId(taskId)
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

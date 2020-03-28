package runtime

import (
	"errors"
	"github.com/devfeel/dottask"
	"github.com/devfeel/rockman/logger"
	"github.com/devfeel/rockman/packets"
	"github.com/devfeel/rockman/runtime/executor"
)

const (
	RuntimeStatus_Init = 0
	RuntimeStatus_Run  = 1
	RuntimeStatus_Stop = 2
)

type (
	Runtime struct {
		TaskService *task.TaskService
		Executors   map[string]executor.Executor
		Status      int
	}
)

func NewRuntime() *Runtime {
	r := &Runtime{Status: RuntimeStatus_Init, Executors: make(map[string]executor.Executor)}
	r.TaskService = task.StartNewService()
	r.TaskService.SetLogger(logger.GetLogger(logger.LoggerName_Runtime))
	logger.Default().Debug("Runtime init success.")
	return r
}

func (r *Runtime) Start() {
	logger.Default().Debug("Runtime start...")
	r.TaskService.StartAllTask()
	r.Status = RuntimeStatus_Run
}

// CreateCronExecutor create new cron executor and register to task service
// now support http\shell\go.so
func (r *Runtime) CreateExecutor(taskConf *packets.TaskConfig) (executor.Executor, error) {
	var exec executor.Executor
	if taskConf.TargetType == executor.TargetType_Http {
		exec = executor.NewHttpExecutor(taskConf)
	} else if taskConf.TargetType == executor.TargetType_Shell {
		exec = executor.NewShellExecutor(taskConf)
	} else if taskConf.TargetType == executor.TargetType_GoSo {
		exec = executor.NewGoExecutor(taskConf)
	}

	err := r.RegisterExecutor(exec)
	return exec, err
}

func (r *Runtime) RegisterExecutor(exec executor.Executor) error {
	task, err := r.TaskService.CreateTask(convertToDotTaskConfig(exec.GetTaskConfig()))
	if err != nil {
		return err
	}
	exec.SetTask(task)
	r.Executors[exec.GetTaskID()] = exec
	return nil
}

func (r *Runtime) StartExecutor(taskId string) error {
	task, exists := r.TaskService.GetTask(taskId)
	if !exists {
		return errors.New("not exists this task")
	}
	task.Start()
	return nil
}

func (r *Runtime) StopExecutor(taskId string) error {
	task, exists := r.TaskService.GetTask(taskId)
	if !exists {
		return errors.New("not exists this task")
	}
	task.Stop()
	return nil
}
func (r *Runtime) RemoveExecutor(taskId string) error {
	task, exists := r.TaskService.GetTask(taskId)
	if !exists {
		return errors.New("not exists this task")
	}
	task.Stop()
	r.TaskService.RemoveTask(taskId)
	return nil
}

func (r *Runtime) QueryAllExecutorConfig() map[string]packets.TaskConfig {
	confs := make(map[string]packets.TaskConfig)
	for key, value := range r.Executors {
		confs[key] = *value.GetTaskConfig()
	}
	return confs
}

func convertToDotTaskConfig(conf *packets.TaskConfig) task.TaskConfig {
	return task.TaskConfig{
		TaskID:   conf.TaskID,
		TaskType: conf.TaskType,
		IsRun:    conf.IsRun,
		DueTime:  conf.DueTime,
		Interval: conf.Interval,
		Express:  conf.Express,
		TaskData: conf.TaskData,
		Handler:  conf.Handler,
	}
}

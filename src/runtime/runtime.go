package runtime

import (
	"github.com/devfeel/dottask"
	"github.com/devfeel/rockman/src/logger"
	"github.com/devfeel/rockman/src/runtime/executor"
)

const (
	Status_Init = 0
	Status_Run  = 1
	Status_Stop = 2
)

type Runtime struct {
	TaskService *task.TaskService
	Executors   map[string]executor.Executor
	Status      int
}

func NewRuntime() *Runtime {
	r := &Runtime{Status: Status_Init, Executors: make(map[string]executor.Executor)}
	r.TaskService = task.StartNewService()
	r.TaskService.SetLogger(logger.GetLogger(logger.LoggerName_Runtime))
	logger.Default().Debug("Runtime Init Success!")
	return r
}

func (r *Runtime) Start() {
	logger.Default().Debug("Runtime Start...")
	r.TaskService.StartAllTask()
	r.Status = Status_Run
}

// CreateCronExecutor create new cron executor and register to task service
// now support http\shell\go.so
func (r *Runtime) CreateExecutor(target string, targetType string, taskConf executor.TaskConfig) (executor.Executor, error) {
	var exec executor.Executor
	if targetType == executor.HttpType {
		exec = executor.NewHttpExecutor(taskConf)
	} else if targetType == executor.ShellType {
		exec = executor.NewShellExecutor(taskConf)
	} else if targetType == executor.ShellType {
		exec = executor.NewGoExecutor(taskConf)
	}

	err := r.RegisterExecutor(exec)
	return exec, err
}

func (r *Runtime) RegisterExecutor(exec executor.Executor) error {
	_, err := r.TaskService.CreateTask(exec.GetDotTaskConfig())
	if err != nil {
		return err
	}
	r.Executors[exec.GetTaskID()] = exec
	return nil
}

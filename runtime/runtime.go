package runtime

import (
	"errors"
	"github.com/devfeel/dottask"
	"github.com/devfeel/rockman/config"
	"github.com/devfeel/rockman/core"
	"github.com/devfeel/rockman/logger"
	"github.com/devfeel/rockman/protected/model"
	"github.com/devfeel/rockman/protected/service"
	"github.com/devfeel/rockman/runtime/executor"
	"sync"
	"time"
)

const (
	RuntimeStatus_Init   = 0
	RuntimeStatus_Run    = 1
	RuntimeStatus_Stop   = 2
	TaskHeader_StartTime = "Rockman.Runtime.StartTime"
)

type (
	Runtime struct {
		TaskService     *task.TaskService
		Executors       map[string]executor.Executor
		executorsLocker *sync.RWMutex
		Status          int
		logLogic        *service.LogService
		nodeInfo        *core.NodeInfo
		config          *config.Profile
	}
)

func NewRuntime(nodeInfo *core.NodeInfo, profile *config.Profile) *Runtime {
	r := &Runtime{Status: RuntimeStatus_Init}
	r.Executors = make(map[string]executor.Executor)
	r.executorsLocker = new(sync.RWMutex)
	r.nodeInfo = nodeInfo
	r.config = profile
	r.TaskService = task.StartNewService()
	r.logLogic = service.NewLogService()
	r.TaskService.SetLogger(logger.Task())
	r.TaskService.SetOnBeforeHandler(func(ctx *task.TaskContext) error {
		ctx.Header[TaskHeader_StartTime] = time.Now()
		return nil
	})
	r.TaskService.SetOnEndHandler(func(ctx *task.TaskContext) error {
		err := r.writeExecLog(ctx)
		if err != nil {
			logger.Runtime().Error(err, "Write ExecLog error")
		}
		return nil
	})
	logger.Default().Debug("Runtime init success.")
	return r
}

func (r *Runtime) Start() error {
	logger.Default().Debug("Runtime start...")
	r.TaskService.StartAllTask()
	r.Status = RuntimeStatus_Run
	return nil
}

func (r *Runtime) Stop() error {
	logger.Default().Debug("Runtime stop.")
	r.TaskService.StopAllTask()
	r.Status = RuntimeStatus_Stop
	return nil
}

// CreateExecutor create new executor and register to task service
// now support http\shell\go.so
func (r *Runtime) CreateExecutor(taskInfo *core.TaskConfig) (executor.Executor, error) {
	var exec executor.Executor
	var err error
	if taskInfo.TargetType == executor.TargetType_Http {
		exec, err = executor.NewHttpExecutor(taskInfo)
	} else if taskInfo.TargetType == executor.TargetType_Shell {
		exec, err = executor.NewShellExecutor(taskInfo)
	} else if taskInfo.TargetType == executor.TargetType_GoSo {
		exec, err = executor.NewGoExecutor(taskInfo)
	}
	if err != nil {
		return nil, err
	}

	err = r.registerExecutor(exec)
	if err == nil {
		if exec.GetTaskConfig().IsRun {
			exec.GetTask().Start()
		}
	}
	return exec, err
}

func (r *Runtime) registerExecutor(exec executor.Executor) error {
	task, err := r.TaskService.CreateTask(convertToDotTaskConfig(exec.GetTaskConfig()))
	if err != nil {
		return err
	}
	exec.SetTask(task)
	r.executorsLocker.Lock()
	r.Executors[exec.GetTaskID()] = exec
	r.executorsLocker.Unlock()

	// reg info to registry
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

	r.executorsLocker.Lock()
	delete(r.Executors, taskId)
	r.executorsLocker.Unlock()

	return nil
}

func (r *Runtime) QueryAllExecutorConfig() map[string]core.TaskConfig {
	r.executorsLocker.RLock()
	defer r.executorsLocker.RUnlock()
	confs := make(map[string]core.TaskConfig)
	for key, value := range r.Executors {
		confs[key] = *value.GetTaskConfig()
	}
	return confs
}

func (r *Runtime) GetTaskIDs() []string {
	var ids []string
	for _, exec := range r.Executors {
		ids = append(ids, exec.GetTaskID())
	}
	return ids
}

func convertToDotTaskConfig(conf *core.TaskConfig) task.TaskConfig {
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

func (r *Runtime) writeExecLog(ctx *task.TaskContext) error {
	var startTime time.Time
	var isSuccess bool
	var failureType, failureCause string
	value, isExists := ctx.Header[TaskHeader_StartTime]
	if isExists {
		startTime = value.(time.Time)
	}
	if ctx.Error != nil {
		isSuccess = false
		failureType = "error"
		failureCause = ctx.Error.Error()
	} else {
		isSuccess = true
	}
	endTime := time.Now()
	execLog := &model.TaskExecLog{
		TaskID:       ctx.TaskID,
		NodeID:       r.nodeInfo.NodeID,
		NodeEndPoint: r.nodeInfo.EndPoint(),
		StartTime:    startTime,
		EndTime:      endTime,
		IsSuccess:    isSuccess,
		FailureType:  failureType,
		FailureCause: failureCause,
	}
	err := r.logLogic.WriteExecLog(execLog)
	return err
}

func registerDemoExecutors(r *Runtime) {
	logger.Node().Debug("Register Demo Executors Begin")
	goExec := executor.NewDebugGoExecutor(("go"))
	err := r.registerExecutor(goExec)
	if err != nil {
		logger.Node().Error(err, "service.CreateCronTask {go.exec} error!")
	}

	httpExec := executor.NewDebugHttpExecutor("http")
	err = r.registerExecutor(httpExec)
	if err != nil {
		logger.Node().Error(err, "service.CreateCronTask {http.exec} error!")
	}

	shellExec := executor.NewDebugShellExecutor("shell")
	err = r.registerExecutor(shellExec)
	if err != nil {
		logger.Node().Error(err, "service.CreateCronTask {shell.exec} error!")
	}
	logger.Node().Debug("Register Demo Executors Success!")
}

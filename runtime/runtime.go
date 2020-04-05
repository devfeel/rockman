package runtime

import (
	"errors"
	"github.com/devfeel/dottask"
	"github.com/devfeel/rockman/core"
	"github.com/devfeel/rockman/logger"
	"github.com/devfeel/rockman/protected/model"
	"github.com/devfeel/rockman/protected/service"
	"github.com/devfeel/rockman/runtime/executor"
	"time"
)

const (
	RuntimeStatus_Init = 0
	RuntimeStatus_Run  = 1
	RuntimeStatus_Stop = 2

	TaskHeader_StartTime = "Rockman.Runtime.StartTime"
)

type (
	Runtime struct {
		TaskService *task.TaskService
		Executors   map[string]executor.Executor
		Status      int
		taskLog     *service.TaskService
		nodeInfo    *core.NodeInfo
	}
)

func NewRuntime(nodeInfo *core.NodeInfo) *Runtime {
	r := &Runtime{Status: RuntimeStatus_Init, Executors: make(map[string]executor.Executor)}
	r.nodeInfo = nodeInfo
	r.TaskService = task.StartNewService()
	r.taskLog = service.NewTaskService()
	r.TaskService.SetLogger(logger.GetLogger(logger.LoggerName_Runtime))
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

// CreateCronExecutor create new cron executor and register to task service
// now support http\shell\go.so
func (r *Runtime) CreateExecutor(taskConf *core.TaskConfig) (executor.Executor, error) {
	var exec executor.Executor
	var err error
	if taskConf.TargetType == executor.TargetType_Http {
		exec, err = executor.NewHttpExecutor(taskConf)
	} else if taskConf.TargetType == executor.TargetType_Shell {
		exec, err = executor.NewShellExecutor(taskConf)
	} else if taskConf.TargetType == executor.TargetType_GoSo {
		exec, err = executor.NewGoExecutor(taskConf)
	}
	if err != nil {
		return nil, err
	}

	err = r.RegisterExecutor(exec)
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

func (r *Runtime) QueryAllExecutorConfig() map[string]core.TaskConfig {
	confs := make(map[string]core.TaskConfig)
	for key, value := range r.Executors {
		confs[key] = *value.GetTaskConfig()
	}
	return confs
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
	err := r.taskLog.WriteExecLog(execLog)
	return err
}

func registerDemoExecutors(r *Runtime) {
	logger.Node().Debug("Register Demo Executors Begin")
	goExec := executor.NewDebugGoExecutor(("go"))
	err := r.RegisterExecutor(goExec)
	if err != nil {
		logger.Node().Error(err, "service.CreateCronTask {go.exec} error!")
	}

	httpExec := executor.NewDebugHttpExecutor("http")
	err = r.RegisterExecutor(httpExec)
	if err != nil {
		logger.Node().Error(err, "service.CreateCronTask {http.exec} error!")
	}

	shellExec := executor.NewDebugShellExecutor("shell")
	err = r.RegisterExecutor(shellExec)
	if err != nil {
		logger.Node().Error(err, "service.CreateCronTask {shell.exec} error!")
	}
	logger.Node().Debug("Register Demo Executors Success!")
}

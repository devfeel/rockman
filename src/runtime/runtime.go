package runtime

import (
	"fmt"
	"github.com/devfeel/dottask"
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
	return &Runtime{Status: Status_Init, Executors: make(map[string]executor.Executor)}
}

func (r *Runtime) Start() {
	r.TaskService = task.StartNewService()
	registerDemoExecutors(r)
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

	err := r.registerExecutor(exec)
	return exec, err
}

func (r *Runtime) registerExecutor(exec executor.Executor) error {
	fmt.Println(exec.GetDotTaskConfig())
	_, err := r.TaskService.CreateTask(exec.GetDotTaskConfig())
	if err != nil {
		return err
	}
	r.Executors[exec.GetTaskID()] = exec
	return nil
}

func registerDemoExecutors(r *Runtime) {
	goExec := executor.NewDebugGoExecutor(("go"))
	err := r.registerExecutor(goExec)
	if err != nil {
		fmt.Println("service.CreateCronTask {go.exec} error! => ", err.Error())
	}

	httpExec := executor.NewDebugHttpExecutor("http")
	err = r.registerExecutor(httpExec)
	if err != nil {
		fmt.Println("service.CreateCronTask {http.exec} error! => ", err.Error())
	}

	shellExec := executor.NewDebugShellExecutor("shell")
	err = r.registerExecutor(shellExec)
	if err != nil {
		fmt.Println("service.CreateCronTask {shell.exec} error! => ", err.Error())
	}
}

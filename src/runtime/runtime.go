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

func registerDemoExecutors(r *Runtime) {
	goExec := executor.NewGoExecutor("go")
	httpExec := executor.NewHttpExecutor("http")
	shellExec := executor.NewShellExecutor("shell")
	_, err := r.TaskService.CreateCronTask("go.exec", true, "0 * * * * *", goExec.Exec, nil)
	if err != nil {
		fmt.Println("service.CreateCronTask {go.exec} error! => ", err.Error())
	}
	r.Executors["go"] = goExec

	_, err = r.TaskService.CreateCronTask("http.exec", true, "0 * * * * *", httpExec.Exec, nil)
	if err != nil {
		fmt.Println("service.CreateCronTask {http.exec} error! => ", err.Error())
	}
	r.Executors["http"] = httpExec

	_, err = r.TaskService.CreateCronTask("shell.exec", true, "0 * * * * *", shellExec.Exec, nil)
	if err != nil {
		fmt.Println("service.CreateCronTask {shell.exec} error! => ", err.Error())
	}
	r.Executors["shell"] = shellExec
}

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
	Service   *task.TaskService
	Executors map[string]executor.Executor
	Status    int
}

func NewRuntime() *Runtime {
	return &Runtime{Status: Status_Init}
}

func (r *Runtime) Start() {
	r.Service = task.StartNewService()
	registerDemoExecutors(r.Service)
	r.Service.StartAllTask()
	r.Status = Status_Run
}

func registerDemoExecutors(service *task.TaskService) {
	goExec := executor.NewGoExecutor("go.so file")
	httpExec := executor.NewHttpExecutor("dotweb.cn")
	shellExec := executor.NewShellExecutor("rock.sh")
	_, err := service.CreateCronTask("go.exec", true, "0 * * * * *", goExec.Exec, nil)
	if err != nil {
		fmt.Println("service.CreateCronTask {go.exec} error! => ", err.Error())
	}
	_, err = service.CreateCronTask("http.exec", true, "0 * * * * *", httpExec.Exec, nil)
	if err != nil {
		fmt.Println("service.CreateCronTask {http.exec} error! => ", err.Error())
	}
	_, err = service.CreateCronTask("shell.exec", true, "0 * * * * *", shellExec.Exec, nil)
	if err != nil {
		fmt.Println("service.CreateCronTask {shell.exec} error! => ", err.Error())
	}
}

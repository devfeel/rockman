package executor

import (
	"fmt"
	"github.com/devfeel/dottask"
)

type HttpExecutor struct {
	baseExecutor
}

func NewDebugHttpExecutor(taskID string) Executor {
	conf := TaskConfig{}
	conf.TaskID = taskID + "-debug"
	conf.TaskType = "cron"
	conf.IsRun = true
	conf.DueTime = 0
	conf.Interval = 0
	conf.Express = "0 * * * * *"
	conf.TaskData = "http-url"
	return NewHttpExecutor(conf)
}

func NewHttpExecutor(conf TaskConfig) *HttpExecutor {
	exec := new(HttpExecutor)
	exec.TaskID = conf.TaskID
	exec.TaskType = conf.TaskType
	exec.IsRun = conf.IsRun
	exec.DueTime = conf.DueTime
	exec.Interval = conf.Interval
	exec.Express = conf.Express
	exec.Handler = exec.Exec
	exec.TaskData = conf.TaskData

	exec.Target = conf.TaskData.(string)
	exec.TargetType = HttpType
	return exec
}

func (exec *HttpExecutor) GetTaskID() string {
	return exec.TaskID
}

func (exec *HttpExecutor) GetTargetType() string {
	return exec.TargetType
}

func (exec *HttpExecutor) Exec(ctx *task.TaskContext) error {
	fmt.Println("HttpExecutor exec", exec.TaskID)
	return nil
}

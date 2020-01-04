package executor

import (
	"fmt"
	"github.com/devfeel/dottask"
)

type GoExecutor struct {
	baseExecutor
}

func NewDebugGoExecutor(taskID string) Executor {
	conf := TaskConfig{}
	conf.TaskID = taskID + "-debug"
	conf.TaskType = "cron"
	conf.IsRun = true
	conf.DueTime = 0
	conf.Interval = 0
	conf.Express = "0 * * * * *"
	conf.TaskData = "go.so"
	return NewGoExecutor(conf)
}

func NewGoExecutor(conf TaskConfig) *GoExecutor {
	exec := new(GoExecutor)
	exec.TaskID = conf.TaskID
	exec.TaskType = conf.TaskType
	exec.IsRun = conf.IsRun
	exec.DueTime = conf.DueTime
	exec.Interval = conf.Interval
	exec.Express = conf.Express
	exec.Handler = exec.Exec
	exec.TaskData = conf.TaskData

	exec.Target = conf.TaskData.(string)
	exec.TargetType = GoSoType
	return exec
}

func (exec *GoExecutor) GetTaskID() string {
	return exec.TaskID
}

func (exec *GoExecutor) GetTargetType() string {
	return exec.TargetType
}

func (exec *GoExecutor) Exec(ctx *task.TaskContext) error {
	fmt.Println("GoExceutor exec", exec.TaskID)
	return nil
}

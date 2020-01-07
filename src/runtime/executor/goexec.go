package executor

import (
	"fmt"
	"github.com/devfeel/dottask"
)

type GoTaskConfig struct {
	TaskConfig
	GoSoFile string
}

type GoExecutor struct {
	baseExecutor
	TaskConfig *GoTaskConfig
}

func NewDebugGoExecutor(taskID string) Executor {
	conf := &GoTaskConfig{}
	conf.TaskID = taskID + "-debug"
	conf.TaskType = "cron"
	conf.IsRun = true
	conf.DueTime = 0
	conf.Interval = 0
	conf.Express = "0 * * * * *"
	conf.TaskData = "go.so"
	return NewGoExecutor(conf)
}

func NewGoExecutor(conf *GoTaskConfig) *GoExecutor {
	exec := new(GoExecutor)
	exec.TargetType = GoSoType
	exec.TaskConfig = conf
	exec.TaskConfig.Handler = exec.Exec
	return exec
}

func (exec *GoExecutor) GetTaskID() string {
	return exec.TaskConfig.TaskID
}

func (exec *GoExecutor) GetTaskConfig() TaskConfig {
	return exec.TaskConfig.TaskConfig
}

func (exec *GoExecutor) Exec(ctx *task.TaskContext) error {
	fmt.Println("GoExceutor exec", exec.TaskConfig.TaskID)
	return nil
}

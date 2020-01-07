package executor

import (
	"fmt"
	"github.com/devfeel/dottask"
)

type HttpTaskConfig struct {
	TaskConfig
	Url         string
	Method      string
	ContentType string
	PostBody    string
	Timeout     int //单位为秒
}

type HttpExecutor struct {
	baseExecutor
	TaskConfig *HttpTaskConfig
}

func NewDebugHttpExecutor(taskID string) Executor {
	conf := &HttpTaskConfig{}
	conf.TaskID = taskID + "-debug"
	conf.TaskType = "cron"
	conf.IsRun = true
	conf.DueTime = 0
	conf.Interval = 0
	conf.Express = "0 * * * * *"
	conf.TaskData = "http-url"
	return NewHttpExecutor(conf)
}

func NewHttpExecutor(conf *HttpTaskConfig) *HttpExecutor {
	exec := new(HttpExecutor)
	exec.TargetType = HttpType
	exec.TaskConfig = conf
	exec.TaskConfig.Handler = exec.Exec
	return exec
}

func (exec *HttpExecutor) GetTaskID() string {
	return exec.TaskConfig.TaskID
}

func (exec *HttpExecutor) GetTaskConfig() TaskConfig {
	return exec.TaskConfig.TaskConfig
}

func (exec *HttpExecutor) Exec(ctx *task.TaskContext) error {
	fmt.Println("HttpExecutor exec", exec.TaskConfig.TaskID)
	return nil
}

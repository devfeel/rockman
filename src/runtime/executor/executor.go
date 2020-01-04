package executor

import (
	"github.com/devfeel/dottask"
)

const (
	HttpType  = "http"
	ShellType = "shell"
	GoSoType  = "goso"
)

type Executor interface {
	GetTaskID() string
	GetTargetType() string
	GetDotTaskConfig() task.TaskConfig
	Exec(ctx *task.TaskContext) error
}

type TaskConfig struct {
	task.TaskConfig
	TargetType string
	Target     string
}

type baseExecutor struct {
	TaskConfig
}

func (exec *baseExecutor) GetDotTaskConfig() task.TaskConfig {
	return exec.TaskConfig.TaskConfig
}

// ValidateExecType validate the execType is supported
func ValidateExecType(execType string) bool {
	if execType != HttpType && execType != GoSoType && execType != ShellType {
		return false
	}
	return true
}

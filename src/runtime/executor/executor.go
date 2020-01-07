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
	GetTaskConfig() TaskConfig
	Exec(ctx *task.TaskContext) error
}

type TaskConfig struct {
	task.TaskConfig
}

type baseExecutor struct {
	TargetType string
}

func (exec *baseExecutor) GetTargetType() string {
	return exec.TargetType
}

// ValidateExecType validate the execType is supported
func ValidateExecType(execType string) bool {
	if execType != HttpType && execType != GoSoType && execType != ShellType {
		return false
	}
	return true
}

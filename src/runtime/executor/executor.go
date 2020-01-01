package executor

import "github.com/devfeel/dottask"

const (
	HttpType  = "http"
	ShellType = "shell"
	GoSoType  = "goso"
)

type Executor interface {
	GetName() string
	GetType() string
	Exec(ctx *task.TaskContext) error
}

// ValidateExecType validate the execType is supported
func ValidateExecType(execType string) bool {
	if execType != HttpType && execType != GoSoType && execType != ShellType {
		return false
	}
	return true
}

package executor

import "github.com/devfeel/dottask"

type Executor interface {
	GetName() string
	Exec(ctx *task.TaskContext) error
}

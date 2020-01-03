package executor

import (
	"fmt"
	"github.com/devfeel/dottask"
)

type GoExecutor struct {
	Name   string
	Type   string
	Target string
}

func NewGoExecutor(name string) Executor {
	exec := &GoExecutor{Name: name, Type: GoSoType}
	return exec
}

func (exec *GoExecutor) GetName() string {
	return exec.Name
}

func (exec *GoExecutor) GetType() string {
	return exec.Type
}

func (exec *GoExecutor) Exec(ctx *task.TaskContext) error {
	fmt.Println("GoExceutor exec", exec.Name)
	return nil
}

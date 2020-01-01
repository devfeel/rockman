package executor

import (
	"fmt"
	"github.com/devfeel/dottask"
)

type GoExecutor struct {
	Name   string
	Target string
}

func NewGoExecutor(name string) Executor {
	return &GoExecutor{Name: name}
}

func (exec *GoExecutor) GetName() string {
	return exec.Name
}

func (exec *GoExecutor) Exec(ctx *task.TaskContext) error {
	fmt.Println("GoExceutor exec", exec.Name)
	return nil
}

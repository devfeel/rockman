package executor

import (
	"fmt"
	"github.com/devfeel/dottask"
)

type HttpExecutor struct {
	Name   string
	Target string
}

func NewHttpExecutor(name string) Executor {
	return &HttpExecutor{Name: name}
}

func (exec *HttpExecutor) GetName() string {
	return exec.Name
}

func (exec *HttpExecutor) Exec(ctx *task.TaskContext) error {
	fmt.Println("HttpExecutor exec", exec.Name)
	return nil
}

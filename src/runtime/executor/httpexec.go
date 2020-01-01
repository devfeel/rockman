package executor

import (
	"fmt"
	"github.com/devfeel/dottask"
)

type HttpExecutor struct {
	Name   string
	Type   string
	Target string
}

func NewHttpExecutor(name string) Executor {
	return &HttpExecutor{Name: name, Type: HttpType}
}

func (exec *HttpExecutor) GetName() string {
	return exec.Name
}

func (exec *HttpExecutor) GetType() string {
	return exec.Type
}

func (exec *HttpExecutor) Exec(ctx *task.TaskContext) error {
	fmt.Println("HttpExecutor exec", exec.Name)
	return nil
}

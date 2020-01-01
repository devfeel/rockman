package executor

import (
	"fmt"
	"github.com/devfeel/dottask"
)

type ShellExecutor struct {
	Name   string
	Target string
}

func NewShellExecutor(name string) Executor {
	return &ShellExecutor{Name: name}
}

func (exec *ShellExecutor) GetName() string {
	return exec.Name
}

func (exec *ShellExecutor) Exec(ctx *task.TaskContext) error {
	fmt.Println("ShellExecutor exec", exec.Name)
	return nil
}

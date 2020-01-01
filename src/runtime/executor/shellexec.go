package executor

import (
	"bytes"
	"fmt"
	"github.com/devfeel/dottask"
	"os/exec"
)

type ShellExecutor struct {
	Name   string
	Type   string
	Target string
}

func NewShellExecutor(name string) Executor {
	return &ShellExecutor{Name: name, Type: ShellType}
}

func (exec *ShellExecutor) GetName() string {
	return exec.Name
}

func (exec *ShellExecutor) GetType() string {
	return exec.Type
}

func (exec *ShellExecutor) Exec(ctx *task.TaskContext) error {
	fmt.Println("ShellExecutor exec", exec.Name)
	result, err := execShell(exec.Target)
	//TODO log exec result
	//1.file
	//2.mysql
	fmt.Print(result)
	return err
}

func execShell(s string) (string, error) {
	cmd := exec.Command("/bin/bash", "-c", s)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return out.String(), err
}

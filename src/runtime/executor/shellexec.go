package executor

import (
	"bytes"
	"fmt"
	"github.com/devfeel/dottask"
	"os/exec"
)

type ShellExecutor struct {
	baseExecutor
}

func NewDebugShellExecutor(taskID string) Executor {
	conf := TaskConfig{}
	conf.TaskID = taskID + "-debug"
	conf.TaskType = "cron"
	conf.IsRun = true
	conf.DueTime = 0
	conf.Interval = 0
	conf.Express = "0 * * * * *"
	conf.TaskData = "shell.sh"
	return NewShellExecutor(conf)
}

func NewShellExecutor(conf TaskConfig) *ShellExecutor {
	exec := new(ShellExecutor)
	exec.TaskID = conf.TaskID
	exec.TaskType = conf.TaskType
	exec.IsRun = conf.IsRun
	exec.DueTime = conf.DueTime
	exec.Interval = conf.Interval
	exec.Express = conf.Express
	exec.Handler = exec.Exec
	exec.TaskData = conf.TaskData

	exec.Target = conf.TaskData.(string)
	exec.TargetType = ShellType
	return exec
}

func (exec *ShellExecutor) GetTaskID() string {
	return exec.TaskID
}

func (exec *ShellExecutor) GetTargetType() string {
	return exec.TargetType
}

func (exec *ShellExecutor) Exec(ctx *task.TaskContext) error {
	fmt.Println("ShellExecutor exec", exec.TaskID)
	result, err := execShell(exec.Target)
	//TODO log exec result
	//1.file
	//2.mysql
	fmt.Print(result)
	return err
}

func execShell(s string) (string, error) {
	return "", nil
	cmd := exec.Command("/bin/bash", "-c", s)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return out.String(), err
}

package executor

import (
	"bytes"
	"fmt"
	"github.com/devfeel/dottask"
	"os/exec"
)

type ShellTaskConfig struct {
	TaskConfig
	ShellFile string
}

type ShellExecutor struct {
	baseExecutor
	TaskConfig *ShellTaskConfig
}

func NewDebugShellExecutor(taskID string) Executor {
	conf := &ShellTaskConfig{}
	conf.TaskID = taskID + "-debug"
	conf.TaskType = "cron"
	conf.IsRun = true
	conf.DueTime = 0
	conf.Interval = 0
	conf.Express = "0 * * * * *"
	conf.TaskData = "shell.sh"
	return NewShellExecutor(conf)
}

func NewShellExecutor(conf *ShellTaskConfig) *ShellExecutor {
	exec := new(ShellExecutor)
	conf.TargetType = ShellType
	exec.TaskConfig = conf
	exec.TaskConfig.Handler = exec.Exec
	exec.baseTaskConfig = &conf.TaskConfig
	return exec
}

func (exec *ShellExecutor) Exec(ctx *task.TaskContext) error {
	fmt.Println("ShellExecutor exec", exec.TaskConfig.TaskID)
	result, err := execShell(exec.TaskConfig.ShellFile)
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

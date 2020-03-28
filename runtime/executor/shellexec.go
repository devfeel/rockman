package executor

import (
	"bytes"
	"fmt"
	"github.com/devfeel/dottask"
	"github.com/devfeel/rockman/logger"
	"github.com/devfeel/rockman/packets"
	"os/exec"
)

type (
	ShellConfig struct {
		FileName string
	}

	ShellExecutor struct {
		baseExecutor
		shellConfig *ShellConfig
	}
)

func NewDebugShellExecutor(taskID string) Executor {
	conf := &packets.TaskConfig{}
	conf.TaskID = taskID + "-debug"
	conf.TaskType = "cron"
	conf.IsRun = true
	conf.DueTime = 0
	conf.Interval = 0
	conf.Express = "0 * * * * *"
	conf.TaskData = "shell.sh"
	conf.TargetType = TargetType_Shell
	conf.TargetConfig = &ShellConfig{
		FileName: "demo.sh",
	}
	return NewShellExecutor(conf)
}

func NewShellExecutor(conf *packets.TaskConfig) *ShellExecutor {
	exec := new(ShellExecutor)
	exec.TaskConfig = conf
	exec.TaskConfig.Handler = exec.Exec
	return exec
}

func (exec *ShellExecutor) Exec(ctx *task.TaskContext) error {
	logTitle := "ShellExecutor [" + exec.GetTaskID() + "] "
	conf, isOk := exec.TaskConfig.TargetConfig.(*ShellConfig)
	if !isOk {
		logger.Runtime().Error(ErrorNotMatchConfigType, logTitle+"convert config error")
		return ErrorNotMatchConfigType
	}
	fmt.Println("ShellExecutor exec", *conf)
	return nil
	result, err := execShell(conf.FileName)
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

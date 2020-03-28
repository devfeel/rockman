package executor

import (
	"bytes"
	"fmt"
	"github.com/devfeel/dottask"
	"github.com/devfeel/mapper"
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
	err := mapper.MapperMap(exec.TaskConfig.TargetConfig.(map[string]interface{}), exec.shellConfig)
	if err != nil {
		logger.Runtime().Error(err, "convert config error")
	}
	return exec
}

func (exec *ShellExecutor) Exec(ctx *task.TaskContext) error {
	logTitle := "ShellExecutor [" + exec.GetTaskID() + "] "
	fmt.Println(logTitle+" success", *exec.shellConfig)
	return nil
	result, err := execShell(exec.shellConfig.FileName)
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

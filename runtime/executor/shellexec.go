package executor

import (
	"bytes"
	"errors"
	"github.com/devfeel/dottask"
	"github.com/devfeel/mapper"
	"github.com/devfeel/rockman/logger"
	"github.com/devfeel/rockman/packets"
	"os/exec"
	"strings"
)

const (
	ShellType_Script = "SCRIPT"
	ShellType_File   = "FILE"
	CorrectResult    = "OK"
)

type (
	ShellConfig struct {
		Type     string //default will be ShellType_Script
		FileName string
		Script   string
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
	exec.shellConfig = new(ShellConfig)
	err := mapper.MapperMap(exec.TaskConfig.TargetConfig.(map[string]interface{}), exec.shellConfig)
	if err != nil {
		logger.Runtime().Error(err, "convert config error")
	}
	if exec.shellConfig.Type == "" {
		exec.shellConfig.Type = ShellType_Script
	}
	exec.shellConfig.Type = strings.ToUpper(exec.shellConfig.Type)
	return exec
}

func (exec *ShellExecutor) Exec(ctx *task.TaskContext) error {
	logTitle := "ShellExecutor [" + exec.GetTaskID() + "] [" + exec.shellConfig.Type + "] "
	if exec.shellConfig.Type == ShellType_Script {
		result, err := execScript(exec.shellConfig.Script)
		logger.Runtime().DebugS(logTitle+"result= "+result, "error=", err)
		if err != nil {
			ctx.Error = err
			return nil
		}
		if result != CorrectResult {
			ctx.Error = errors.New("shell response not " + CorrectResult + ", is " + result)
		}
		return nil
	}
	logger.Runtime().Debug(logTitle + "not support shell type")
	ctx.Error = errors.New("not support shell type [" + exec.shellConfig.Type + "]")
	return nil
}

func execScript(s string) (string, error) {
	cmd := exec.Command("/bin/bash", "-c", s)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	str := strings.Replace(out.String(), " ", "", -1)
	str = strings.Replace(out.String(), "\n", "", -1)
	return str, err
}

package executor

import (
	"bytes"
	"errors"
	"github.com/devfeel/dottask"
	"github.com/devfeel/mapper"
	"github.com/devfeel/rockman/config"
	"github.com/devfeel/rockman/core"
	"github.com/devfeel/rockman/logger"
	_file "github.com/devfeel/rockman/util/file"
	_json "github.com/devfeel/rockman/util/json"
	"os/exec"
	"strings"
)

const (
	ShellType_Script = "SCRIPT"
	ShellType_File   = "FILE"
	ShellFilePath    = "shells/"
)

var (
	ErrorNotSupportShellType       = errors.New("not support shell type")
	ErrorShellFileNotInSpecifyPath = errors.New("shell file not in specify path")
	ErrorNotEnabledShellScriptMode = errors.New("shell script mode is not enabled")
)

type (
	ShellConfig struct {
		Type   string //default will be ShellType_Script
		Script string
	}

	ShellExecutor struct {
		baseExecutor
		shellConfig *ShellConfig
	}
)

func NewDebugShellExecutor(taskID string) Executor {
	conf := &core.TaskConfig{}
	conf.TaskID = taskID + "-debug"
	conf.TaskType = "cron"
	conf.IsRun = true
	conf.DueTime = 0
	conf.Interval = 0
	conf.Express = "0 * * * * *"
	conf.TaskData = "shell.sh"
	conf.TargetType = TargetType_Shell
	conf.TargetConfig = &ShellConfig{
		Script: "demo.sh",
		Type:   ShellType_File,
	}
	exec, _ := NewShellExecutor(conf)
	return exec
}

func NewShellExecutor(conf *core.TaskConfig) (*ShellExecutor, error) {
	lt := "NewShellExecutor[" + conf.TaskID + "] "
	exec := new(ShellExecutor)
	exec.TaskConfig = conf
	exec.TaskConfig.Handler = exec.Exec
	exec.shellConfig = new(ShellConfig)
	err := mapper.MapperMap(exec.TaskConfig.TargetConfig.(map[string]interface{}), exec.shellConfig)
	if err != nil {
		logger.Runtime().Error(err, lt+"convert config error")
		return nil, err
	}
	if exec.shellConfig.Type == "" {
		exec.shellConfig.Type = ShellType_File
	}
	exec.shellConfig.Type = strings.ToUpper(exec.shellConfig.Type)
	if !config.CurrentProfile.Runtime.EnableShellScript {
		if exec.shellConfig.Type == ShellType_Script {
			logger.Runtime().Debug(lt + "error: " + ErrorNotEnabledShellScriptMode.Error())
			return nil, ErrorNotEnabledShellScriptMode
		}
	}

	if exec.shellConfig.Type != ShellType_File && exec.shellConfig.Type != ShellType_Script {
		logger.Runtime().Debug(lt + "error: " + ErrorNotSupportShellType.Error())
		return nil, ErrorNotSupportShellType
	}

	if exec.shellConfig.Type == ShellType_File {
		exec.shellConfig.Script = ShellFilePath + exec.shellConfig.Script
		if !_file.ExistsInPath(ShellFilePath, exec.shellConfig.Script) {
			logger.Runtime().Debug(lt + "error: " + ErrorShellFileNotInSpecifyPath.Error())
			return nil, ErrorShellFileNotInSpecifyPath
		}
	}
	return exec, nil
}

func (exec *ShellExecutor) Exec(ctx *task.TaskContext) error {
	logTitle := "ShellExecutor [" + exec.GetTaskID() + "] [" + exec.shellConfig.Type + "] "

	var result string
	var err error
	if exec.shellConfig.Type == ShellType_Script {
		result, err = execShellScript(exec.shellConfig.Script)
	}
	if exec.shellConfig.Type == ShellType_File {
		result, err = execShellFile(exec.shellConfig.Script)
	}

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

func execShellScript(s string) (string, error) {
	if !config.CurrentProfile.Runtime.EnableShellScript {
		return "", ErrorNotEnabledShellScriptMode
	}
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

func execShellFile(f string) (string, error) {
	cmd := exec.Command("/bin/sh", f)
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

func (c *ShellConfig) LoadFromJson(json string) error {
	return _json.Unmarshal(json, c)
}

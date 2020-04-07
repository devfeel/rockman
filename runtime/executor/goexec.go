package executor

import (
	"errors"
	"github.com/devfeel/dottask"
	"github.com/devfeel/mapper"
	"github.com/devfeel/rockman/core"
	"github.com/devfeel/rockman/logger"
	_file "github.com/devfeel/rockman/util/file"
	_json "github.com/devfeel/rockman/util/json"
	"plugin"
)

/*
	build script on linux: go build --buildmode=plugin -o plugin.so plugin.go
*/

const GoFilePath = "plugins/"

var (
	ErrorGoSoFileNotInSpecifyPath = errors.New("go.so file not in specify path")
)

type (
	GoConfig struct {
		FileName string
	}

	GoExecutor struct {
		baseExecutor
		goConfig *GoConfig
	}
)

func NewDebugGoExecutor(taskID string) Executor {
	conf := &core.TaskConfig{}
	conf.TaskID = taskID + "-debug"
	conf.TaskType = "cron"
	conf.IsRun = true
	conf.DueTime = 0
	conf.Interval = 0
	conf.Express = "0 * * * * *"
	conf.TaskData = "go.so"
	conf.TargetType = TargetType_GoSo
	conf.TargetConfig = &GoConfig{
		FileName: "demo.so",
	}
	exec, _ := NewGoExecutor(conf)
	return exec
}

func NewGoExecutor(conf *core.TaskConfig) (*GoExecutor, error) {
	exec := new(GoExecutor)
	exec.TaskConfig = conf
	exec.TaskConfig.Handler = exec.Exec
	exec.goConfig = new(GoConfig)
	err := mapper.MapperMap(exec.TaskConfig.TargetConfig.(map[string]interface{}), exec.goConfig)
	if err != nil {
		logger.Runtime().Error(err, "convert config error")
		return nil, err
	}

	exec.goConfig.FileName = GoFilePath + exec.goConfig.FileName
	if !_file.ExistsInPath(GoFilePath, exec.goConfig.FileName) {
		logger.Runtime().Debug("NewGoExecutor error: " + ErrorGoSoFileNotInSpecifyPath.Error())
		return nil, ErrorGoSoFileNotInSpecifyPath
	}
	return exec, nil
}

func (exec *GoExecutor) Exec(ctx *task.TaskContext) error {
	logTitle := "GoExecutor [" + exec.GetTaskID() + "] [" + exec.goConfig.FileName + "] "
	p, err := plugin.Open(exec.goConfig.FileName)
	if err != nil {
		logger.Runtime().Error(err, logTitle+"open plugin error: "+err.Error())
		ctx.Error = err
		return nil
	}
	s, err := p.Lookup("Exec")
	if err != nil {
		logger.Runtime().Error(err, logTitle+"lookup Exec error: "+err.Error())
		ctx.Error = err
		return nil
	}
	if execFunc, ok := s.(func(ctx *task.TaskContext) error); ok {
		err := execFunc(ctx)
		if err != nil {
			logger.Runtime().Error(err, logTitle+"exec err:"+err.Error())
		} else {
			logger.Runtime().DebugS(logTitle + "exec success")
		}
		ctx.Error = err
		return nil
	} else {
		err := errors.New("not match Exec function")
		logger.Runtime().Error(err, logTitle+"not match Exec function")
		ctx.Error = err
		return nil
	}
}

func (c *GoConfig) LoadFromJson(json string) error {
	return _json.Unmarshal(json, c)
}

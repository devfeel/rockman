package executor

import (
	"errors"
	"github.com/devfeel/dottask"
	"github.com/devfeel/rockman/logger"
	"github.com/devfeel/rockman/packets"
	"plugin"
)

type (
	GoConfig struct {
		FileName string
	}

	GoExecutor struct {
		baseExecutor
	}
)

func NewDebugGoExecutor(taskID string) Executor {
	conf := &packets.TaskConfig{}
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
	return NewGoExecutor(conf)
}

func NewGoExecutor(conf *packets.TaskConfig) *GoExecutor {
	exec := new(GoExecutor)
	exec.TaskConfig = conf
	exec.TaskConfig.Handler = exec.Exec
	return exec
}

// Exec TODO:log to mysql log
func (exec *GoExecutor) Exec(ctx *task.TaskContext) error {
	logTitle := "GoExecutor [" + exec.GetTaskID() + "] "
	conf, isOk := exec.TaskConfig.TargetConfig.(*GoConfig)
	if !isOk {
		logger.Runtime().Error(ErrorNotMatchConfigType, logTitle+"convert config error")
		return ErrorNotMatchConfigType
	}
	p, err := plugin.Open(conf.FileName)
	if err != nil {
		logger.Runtime().Error(err, logTitle+"open plugin error: "+err.Error())
		return err
	}
	s, err := p.Lookup("Exec")
	if err != nil {
		logger.Runtime().Error(err, logTitle+"lookup Exec error: "+err.Error())
		return err
	}
	if execFunc, ok := s.(Exec); ok {
		err := execFunc(ctx)
		if err != nil {
			logger.Runtime().DebugS(logTitle + "exec success")
		} else {
			logger.Runtime().Error(err, logTitle+"exec err:"+err.Error())
		}
		return err
	} else {
		err := errors.New("not match Exec function")
		logger.Runtime().Error(err, logTitle+"not match Exec function")
		return err
	}
}

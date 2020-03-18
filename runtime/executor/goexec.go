package executor

import (
	"errors"
	"github.com/devfeel/dottask"
	"github.com/devfeel/rockman/logger"
	"plugin"
)

type GoTaskConfig struct {
	TaskConfig
	GoSoFile string
}

type GoExecutor struct {
	baseExecutor
	TaskConfig *GoTaskConfig
}

func NewDebugGoExecutor(taskID string) Executor {
	conf := &GoTaskConfig{}
	conf.TaskID = taskID + "-debug"
	conf.TaskType = "cron"
	conf.IsRun = true
	conf.DueTime = 0
	conf.Interval = 0
	conf.Express = "0 * * * * *"
	conf.TaskData = "go.so"
	return NewGoExecutor(conf)
}

func NewGoExecutor(conf *GoTaskConfig) *GoExecutor {
	exec := new(GoExecutor)
	conf.TargetType = GoSoType
	exec.TaskConfig = conf
	exec.TaskConfig.Handler = exec.Exec
	exec.baseTaskConfig = &conf.TaskConfig
	return exec
}

// Exec TODO:log to mysql log
func (exec *GoExecutor) Exec(ctx *task.TaskContext) error {
	logTtitle := "GoExceutor [" + exec.TaskConfig.TaskID + "] "
	p, err := plugin.Open(exec.TaskConfig.GoSoFile)
	if err != nil {
		logger.Runtime().Error(err, logTtitle+"error open plugin: "+err.Error())
		return err
	}
	s, err := p.Lookup("Exec")
	if err != nil {
		logger.Runtime().Error(err, logTtitle+"error lookup Exec: "+err.Error())
		return err
	}
	if execFunc, ok := s.(Exec); ok {
		err := execFunc(ctx)
		if err != nil {
			logger.Runtime().DebugS(logTtitle + "exec success")
		} else {
			logger.Runtime().Error(err, logTtitle+"exec err:"+err.Error())
		}
		return err
	} else {
		err := errors.New("not match Exec function")
		logger.Runtime().Error(err, logTtitle+"not match Exec function")
		return err
	}
}

package executor

import (
	"fmt"
	"github.com/devfeel/dottask"
	"github.com/devfeel/rockman/logger"
	"github.com/devfeel/rockman/packets"
)

type (
	HttpConfig struct {
		Url         string
		Method      string
		ContentType string
		PostBody    string
		Timeout     int //单位为秒
	}

	HttpExecutor struct {
		baseExecutor
	}
)

func NewDebugHttpExecutor(taskID string) Executor {
	conf := &packets.TaskConfig{}
	conf.TaskID = taskID + "-debug"
	conf.TaskType = "cron"
	conf.IsRun = true
	conf.DueTime = 0
	conf.Interval = 0
	conf.Express = "0 * * * * *"
	conf.TaskData = "http-url"
	conf.TargetType = TargetType_Http
	conf.TargetConfig = &HttpConfig{
		Url:    "http://www.dotweb.cn",
		Method: "GET",
	}
	return NewHttpExecutor(conf)
}

func NewHttpExecutor(conf *packets.TaskConfig) *HttpExecutor {
	exec := new(HttpExecutor)
	exec.TaskConfig = conf
	exec.TaskConfig.Handler = exec.Exec
	return exec
}

func (exec *HttpExecutor) Exec(ctx *task.TaskContext) error {
	logTitle := "HttpExecutor [" + exec.GetTaskID() + "] "
	conf, isOk := exec.TaskConfig.TargetConfig.(*HttpConfig)
	if !isOk {
		logger.Runtime().Error(ErrorNotMatchConfigType, logTitle+"convert config error")
		return ErrorNotMatchConfigType
	}
	fmt.Println("HttpExecutor exec", *conf)
	return nil
}

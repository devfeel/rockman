package executor

import (
	"fmt"
	"github.com/devfeel/dottask"
	"github.com/devfeel/mapper"
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
		httpConfig *HttpConfig
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
	exec.httpConfig = new(HttpConfig)
	err := mapper.MapperMap(exec.TaskConfig.TargetConfig.(map[string]interface{}), exec.httpConfig)
	if err != nil {
		logger.Runtime().Error(err, "convert config error")
	}
	return exec
}

func (exec *HttpExecutor) Exec(ctx *task.TaskContext) error {
	logTitle := "HttpExecutor [" + exec.GetTaskID() + "] "
	fmt.Println(logTitle+"exec", exec.httpConfig)
	return nil
}

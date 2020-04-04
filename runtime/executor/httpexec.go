package executor

import (
	"errors"
	"github.com/devfeel/dottask"
	"github.com/devfeel/mapper"
	"github.com/devfeel/rockman/core"
	"github.com/devfeel/rockman/logger"
	_http "github.com/devfeel/rockman/util/http"
	"strings"
	"time"
)

const (
	HttpMethod_HEAD = "HEAD"
	HttpMethod_GET  = "GET"
	HttpMethod_POST = "POST"
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
	conf := &core.TaskConfig{}
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
		Method: HttpMethod_GET,
	}
	exec, _ := NewHttpExecutor(conf)
	return exec
}

func NewHttpExecutor(conf *core.TaskConfig) (*HttpExecutor, error) {
	exec := new(HttpExecutor)
	exec.TaskConfig = conf
	exec.TaskConfig.Handler = exec.Exec
	exec.httpConfig = new(HttpConfig)
	err := mapper.MapperMap(exec.TaskConfig.TargetConfig.(map[string]interface{}), exec.httpConfig)
	if err != nil {
		logger.Runtime().Error(err, "convert config error")
		return nil, err
	}
	if exec.httpConfig.Method == "" {
		exec.httpConfig.Method = HttpMethod_GET
	}
	exec.httpConfig.Method = strings.ToUpper(exec.httpConfig.Method)
	return exec, nil
}

func (exec *HttpExecutor) Exec(ctx *task.TaskContext) error {
	logTitle := "HttpExecutor [" + exec.GetTaskID() + "] [" + exec.httpConfig.Method + "] "
	if exec.httpConfig.Method == HttpMethod_HEAD {
		result := _http.HttpHead(exec.httpConfig.Url, time.Second*time.Duration(exec.httpConfig.Timeout))
		logger.Runtime().DebugS(logTitle+"result= "+result.Status, "error=", result.Error)
		if result.Error != nil {
			ctx.Error = result.Error
			return nil
		}
		if result.Status != CorrectStatus {
			ctx.Error = errors.New("http response status not " + CorrectStatus + ", is " + result.Status)
		}
		return nil
	}
	if exec.httpConfig.Method == HttpMethod_GET {
		result := _http.HttpGet(exec.httpConfig.Url, time.Second*time.Duration(exec.httpConfig.Timeout))
		logger.Runtime().DebugS(logTitle+"result= "+result.Status, "error=", result.Error)
		if result.Error != nil {
			ctx.Error = result.Error
			return nil
		}
		if result.Status != CorrectStatus {
			ctx.Error = errors.New("http response status not " + CorrectStatus + ", is " + result.Status)
		}
		return nil
	}
	if exec.httpConfig.Method == HttpMethod_POST {
		result := _http.HttpPost(exec.httpConfig.Url, exec.httpConfig.PostBody, exec.httpConfig.ContentType, time.Second*time.Duration(exec.httpConfig.Timeout))
		logger.Runtime().DebugS(logTitle+"result= "+result.Status, "error=", result.Error)
		if result.Error != nil {
			ctx.Error = result.Error
			return nil
		}
		if result.Status != CorrectStatus {
			ctx.Error = errors.New("http response status not " + CorrectStatus + ", is " + result.Status)
		}
		return nil
	}

	logger.Runtime().Debug(logTitle + "not support http method [" + exec.httpConfig.Method + "]")
	ctx.Error = errors.New(logTitle + "not support http method [" + exec.httpConfig.Method + "]")
	return nil
}

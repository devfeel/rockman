package model

import (
	"errors"
	"fmt"
	"time"

	"github.com/devfeel/rockman/core"
	"github.com/devfeel/rockman/logger"
	"github.com/devfeel/rockman/runtime/executor"
)

type ExecutorInfo struct {
	ID                int64
	TaskID            string
	TaskType          string
	IsRun             bool
	DueTime           int64
	Interval          int64
	Express           string
	TaskData          string
	TargetType        string
	TargetConfig      string
	RealTargetConfig  interface{}
	NodeID            string
	DistributeType    int
	IsSubmitToCluster bool
	Remark            string
	CreateTime        time.Time
}

func (e *ExecutorInfo) TaskConfig() *core.TaskConfig {
	e.InitTargetConfig()

	defer func() {
		if err := recover(); err != nil {
			errInfo := errors.New(fmt.Sprintln(err))
			logger.Default().Error(errInfo, "ExecutorInfo.TaskConfig() throw unhandled error:"+errInfo.Error())
		}
	}()

	conf := &core.TaskConfig{}
	conf.TaskID = e.TaskID
	conf.TaskType = e.TaskType
	conf.TargetType = e.TargetType
	conf.IsRun = e.IsRun
	conf.DueTime = e.DueTime
	conf.Interval = e.Interval
	conf.Express = e.Express
	conf.TaskData = e.TaskData
	conf.HAFlag = true
	if e.TargetType == executor.TargetType_Http {
		conf.TargetConfig = e.RealTargetConfig.(*executor.HttpConfig)
	}
	if e.TargetType == executor.TargetType_GoSo {
		conf.TargetConfig = e.RealTargetConfig.(*executor.GoConfig)
	}
	if e.TargetType == executor.TargetType_Shell {
		conf.TargetConfig = e.RealTargetConfig.(*executor.ShellConfig)
	}
	return conf
}

func (e *ExecutorInfo) InitTargetConfig() {
	if e.RealTargetConfig != nil {
		return
	}
	if e.TargetType == executor.TargetType_Http {
		conf := new(executor.HttpConfig)
		err := conf.LoadFromJson(e.TargetConfig)
		if err != nil {
			e.RealTargetConfig = conf
		}
	}
	if e.TargetType == executor.TargetType_GoSo {
		conf := new(executor.GoConfig)
		err := conf.LoadFromJson(e.TargetConfig)
		if err != nil {
			e.RealTargetConfig = conf
		}
	}
	if e.TargetType == executor.TargetType_Shell {
		conf := new(executor.ShellConfig)
		err := conf.LoadFromJson(e.TargetConfig)
		if err != nil {
			e.RealTargetConfig = conf
		}
	}
}

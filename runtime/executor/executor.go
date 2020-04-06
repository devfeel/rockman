package executor

import (
	"errors"
	"github.com/devfeel/dottask"
	"github.com/devfeel/rockman/core"
)

const (
	TargetType_Http  = "http"
	TargetType_Shell = "shell"
	TargetType_GoSo  = "goso"
	TargetType_Code  = "code"

	CorrectStatus = "200 OK"
	CorrectResult = "OK"
)

type (
	Executor interface {
		GetTask() task.Task
		SetTask(task.Task)
		GetTaskID() string
		GetTargetType() string
		GetTaskConfig() *core.TaskConfig
		Exec(*task.TaskContext) error
	}

	Exec func(ctx *task.TaskContext) error

	baseExecutor struct {
		Task       task.Task
		TaskConfig *core.TaskConfig
	}
)

var ErrorNotMatchConfigType = errors.New("not match config type")

func (exec *baseExecutor) GetTask() task.Task {
	return exec.Task
}

func (exec *baseExecutor) SetTask(task task.Task) {
	exec.Task = task
}

func (exec *baseExecutor) GetTaskID() string {
	return exec.TaskConfig.TaskID
}

func (exec *baseExecutor) GetTaskConfig() *core.TaskConfig {
	return exec.TaskConfig
}

func (exec *baseExecutor) GetTargetType() string {
	return exec.TaskConfig.TargetType
}

// ValidateExecType validate the execType is supported
func ValidateExecType(execType string) bool {
	if execType != TargetType_Http && execType != TargetType_GoSo && execType != TargetType_Shell {
		return false
	}
	return true
}

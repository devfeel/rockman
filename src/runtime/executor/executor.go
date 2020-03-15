package executor

import (
	"encoding/json"
	"errors"
	"github.com/devfeel/dottask"
)

const (
	HttpType  = "http"
	ShellType = "shell"
	GoSoType  = "goso"
	CodeType  = "code"
)

type Executor interface {
	GetTask() task.Task
	SetTask(task task.Task)
	GetTaskID() string
	GetTargetType() string
	GetTaskConfig() *TaskConfig
	Exec(ctx *task.TaskContext) error
}

type Exec func(ctx *task.TaskContext) error

type TaskConfig struct {
	TaskID     string
	TaskType   string
	IsRun      bool
	Handler    task.TaskHandle `json:"-"`
	DueTime    int64
	Interval   int64
	Express    string
	TargetType string
	TaskData   interface{}
}

type baseExecutor struct {
	Task           task.Task
	baseTaskConfig *TaskConfig
}

func (exec *baseExecutor) GetTask() task.Task {
	return exec.Task
}

func (exec *baseExecutor) SetTask(task task.Task) {
	exec.Task = task
}

func (exec *baseExecutor) GetTaskID() string {
	return exec.baseTaskConfig.TaskID
}

func (exec *baseExecutor) GetTaskConfig() *TaskConfig {
	return exec.baseTaskConfig
}

func (exec *baseExecutor) GetTargetType() string {
	return exec.baseTaskConfig.TargetType
}

// ValidateExecType validate the execType is supported
func ValidateExecType(execType string) bool {
	if execType != HttpType && execType != GoSoType && execType != ShellType {
		return false
	}
	return true
}

func ConvertRealTaskConfig(taskConfig *TaskConfig) (interface{}, error) {
	jsonStr, err := json.Marshal(taskConfig)
	if err != nil {
		return nil, err
	}
	if taskConfig.TargetType == HttpType {
		conf := &HttpTaskConfig{}
		err := json.Unmarshal([]byte(jsonStr), conf)
		return conf, err
	} else if taskConfig.TargetType == ShellType {
		conf := &HttpTaskConfig{}
		err := json.Unmarshal([]byte(jsonStr), conf)
		return conf, err
	} else if taskConfig.TargetType == ShellType {
		conf := &HttpTaskConfig{}
		err := json.Unmarshal([]byte(jsonStr), conf)
		return conf, err
	}
	return nil, errors.New("not support target type: " + taskConfig.TaskType)
}

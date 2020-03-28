package packets

import task "github.com/devfeel/dottask"

type TaskConfig struct {
	TaskID       string
	TaskType     string
	IsRun        bool
	Handler      task.TaskHandle `json:"-"`
	DueTime      int64
	Interval     int64
	Express      string
	TaskData     interface{}
	TargetType   string
	TargetConfig interface{}
}

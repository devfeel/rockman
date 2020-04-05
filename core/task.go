package core

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
	HAFlag       bool //HA flag, if set true, leader will watch it, when it offline will resubmit
	TargetType   string
	TargetConfig interface{}
}

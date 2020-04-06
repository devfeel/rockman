package model

type TaskInfo struct {
	ID           int
	TaskID       string
	TaskType     string
	IsRun        bool
	DueTime      int64
	Interval     int64
	Express      string
	TaskData     interface{}
	HAFlag       bool //HA flag, if set true, leader will watch it, when it offline will resubmit
	TargetType   string
	TargetConfig interface{}
}

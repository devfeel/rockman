package model

type TaskInfo struct {
	ID         int
	TargetType string //http/shell/goso
	TaskID     string
	TaskType   string
	IsRun      bool
	DueTime    int64
	Interval   int64
	Express    string
	TaskData   string
}

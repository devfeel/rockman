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

type HttpTaskInfo struct {
	TaskInfo
	Url         string
	Method      string
	ContentType string
	PostBody    string
	Timeout     int //单位为秒
}

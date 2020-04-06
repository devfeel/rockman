package model

type ExecutorInfo struct {
	ID             int64
	TaskID         string
	TaskType       string
	IsRun          bool
	DueTime        int64
	Interval       int64
	Express        string
	TaskData       interface{}
	TargetType     string
	TargetConfig   interface{}
	NodeID         string
	DistributeType int
	Remark         string
}

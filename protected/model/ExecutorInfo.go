package model

type ExecutorInfo struct {
	ID               int64
	TaskID           string
	TaskType         string
	IsRun            bool
	DueTime          int64
	Interval         int64
	Express          string
	TaskData         string
	TargetType       string
	TargetConfig     string
	RealTargetConfig interface{}
	NodeID           string
	DistributeType   int
	Remark           string
}

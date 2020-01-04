package model

type DemoInfo struct {
	ID         int
	NodeID     string
	TargetType string //http/shell/goso
	Target     string
	TaskID     string
	TaskType   string
	IsRun      bool
	DueTime    int64
	Interval   int64
	Express    string
}

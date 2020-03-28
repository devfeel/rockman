package model

import "time"

type TaskExecLog struct {
	LogID        int64
	TaskID       string
	NodeID       string
	NodeEndPoint string
	StartTime    time.Time
	EndTime      time.Time
	IsSuccess    bool
	FailureType  string
	FailureCause string
	CreateTime   time.Time
}

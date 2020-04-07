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

type TaskStateLog struct {
	LogID        int64
	TaskID       string
	NodeID       string
	NodeEndPoint string
	State        bool
	Message      string
	CreateTime   time.Time
}

type NodeTraceLog struct {
	LogID        int64
	NodeID       string
	NodeEndPoint string
	IsLeader     bool
	IsMaster     bool
	IsWorker     bool
	Event        string
	IsSuccess    bool
	FailureType  string
	FailureCause string
	CreateTime   time.Time
}

package contract

import "github.com/devfeel/rockman/protected/model"

type ExecutorQR struct {
	model.PageRequest
	NodeID string
}

type TaskExecLogQR struct {
	model.PageRequest
	TaskID string
}

type TaskStateLogQR struct {
	model.PageRequest
	TaskID string
}

type TaskSubmitLogQR struct {
	model.PageRequest
	TaskID string
}

type NodeTraceLogQR struct {
	model.PageRequest
	NodeID string
}

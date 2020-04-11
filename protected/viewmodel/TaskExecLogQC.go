package viewmodel

import "github.com/devfeel/rockman/protected/model"

type TaskExecLogQC struct {
	model.PageRequest
	TaskID string
}

package core

import _json "github.com/devfeel/rockman/util/json"

type ExecutorInfo struct {
	TaskConfig *TaskConfig
	Worker     *NodeInfo
}

func (e *ExecutorInfo) Json() string {
	return _json.GetJsonString(e)
}

func (e *ExecutorInfo) LoadFromJson(json string) error {
	return _json.Unmarshal(json, e)
}

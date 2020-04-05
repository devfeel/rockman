package core

import _json "github.com/devfeel/rockman/util/json"

type ExecutorInfo struct {
	TaskID         string
	IsOnline       bool
	Node           *NodeInfo
	DistributeType int
}

func (n *ExecutorInfo) Json() string {
	return _json.GetJsonString(n)
}

func (n *ExecutorInfo) LoadFromJson(json string) error {
	return _json.Unmarshal(json, n)
}

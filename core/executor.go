package core

import _json "github.com/devfeel/rockman/util/json"

type ExecutorInfo struct {
	TaskID   string
	IsOnline bool
	Config   *TaskConfig
	Node     *NodeInfo
}

func (n *ExecutorInfo) Json() string {
	return _json.GetJsonString(n)
}

func (n *ExecutorInfo) LoadFromJson(json string) error {
	return _json.Unmarshal(json, n)
}

func GetExecutorKeyPrefix(clusterId string) string {
	return ClusterKeyPrefix + clusterId + "/executors/"
}

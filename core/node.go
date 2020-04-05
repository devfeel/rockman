package core

import "github.com/devfeel/rockman/util/json"

const ClusterKeyPrefix = "devfeel/rockman/"

type NodeInfo struct {
	NodeID    string
	Cluster   string
	Host      string
	Port      string
	OuterHost string
	OuterPort string
	IsMaster  bool
	IsWorker  bool
	IsOnline  bool
}

func (n *NodeInfo) EndPoint() string {
	host := n.Host
	if n.OuterHost != "" {
		host = n.OuterHost
	}
	port := n.Port
	if n.OuterPort != "" {
		host = n.OuterPort
	}
	return host + ":" + port
}

func (n *NodeInfo) Json() string {
	return _json.GetJsonString(n)
}

func (n *NodeInfo) LoadFromJson(json string) error {
	return _json.Unmarshal(json, n)
}

func (n *NodeInfo) GetEmptyResource() *ResourceInfo {
	return &ResourceInfo{EndPoint: n.EndPoint()}
}

func (n *NodeInfo) GetNodeKey(clusterId string) string {
	return GetNodeKeyPrefix(clusterId) + n.EndPoint()
}

func GetNodeKeyPrefix(clusterId string) string {
	return ClusterKeyPrefix + clusterId + "/nodes/"
}

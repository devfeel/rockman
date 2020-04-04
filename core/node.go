package core

import jsonutil "github.com/devfeel/rockman/util/json"

const NodeKeyPrefix = "devfeel/rockman/nodekey/"

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
	return jsonutil.GetJsonString(n)
}

func (n *NodeInfo) LoadFromJson(json string) error {
	return jsonutil.Unmarshal(json, n)
}

func (n *NodeInfo) GetNodeKey(clusterId string) string {
	return NodeKeyPrefix + clusterId + ":" + n.EndPoint()
}

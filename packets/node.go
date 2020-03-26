package packets

import jsonutil "github.com/devfeel/rockman/util/json"

const NodeKeyPrefix = "devfeel/rockman:nodekey:"

type NodeInfo struct {
	NodeID   string
	Cluster  string
	Host     string
	Port     string
	IsMaster bool
	IsWorker bool
	IsOnline bool
}

func (n *NodeInfo) EndPoint() string {
	return n.Host + ":" + n.Port
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

package packets

type NodeInfo struct {
	NodeID string
	Host   string
	Port   string
}

func (n *NodeInfo) EndPoint() string {
	return n.Host + ":" + n.Port
}

func (n *NodeInfo) GetNodeKey(clusterId string) string {
	return "devfeel/rockman:" + clusterId + ":node:" + n.EndPoint()
}

package packets

type NodeInfo struct {
	Host     string
	Port     string
	NodeID   string
	IsWorker bool
	IsMaster bool
}

package rpc

type JsonRequest struct {
	Version string
	Command string
	Message interface{}
}

type JsonResult struct {
	RetCode int
	RetMsg  string
	Message interface{}
}

type NodeInfo struct {
	Host     string
	Port     string
	NodeID   string
	IsWorker bool
	IsMaster bool
}

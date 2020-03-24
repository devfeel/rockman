package packets

type SubmitInfo struct {
	ExecutorConfig interface{}
	Worker         *NodeInfo
	DistributeType int ``
}

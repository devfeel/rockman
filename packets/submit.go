package packets

type SubmitInfo struct {
	TaskConfig     *TaskConfig
	Worker         *NodeInfo
	DistributeType int ``
}

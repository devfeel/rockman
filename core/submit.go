package core

type SubmitInfo struct {
	TaskConfig     *TaskConfig
	Worker         *NodeInfo
	DistributeType int
}

func (s *SubmitInfo) ExecutorInfo() *ExecutorInfo {
	return &ExecutorInfo{
		TaskID:         s.TaskConfig.TaskID,
		Node:           s.Worker,
		DistributeType: s.DistributeType,
	}
}

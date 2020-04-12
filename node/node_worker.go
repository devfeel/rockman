package node

import (
	"github.com/devfeel/rockman/core"
	"github.com/devfeel/rockman/logger"
)

// RegisterExecutor
func (n *Node) RegisterExecutor(taskInfo *core.TaskConfig) *core.Result {
	logTitle := "Node RegisterExecutor [" + taskInfo.TaskID + "] "
	if !n.IsWorker() {
		logger.Node().Debug(logTitle + "failed, current node is not worker.")
		return core.FailedResult(-1001, "current node is not worker")
	}
	_, err := n.Runtime.CreateExecutor(taskInfo)
	if err != nil {
		logger.Node().Warn(logTitle + "CreateExecutor error:" + err.Error())
		return core.FailedResult(-2001, err.Error())
	} else {
		// reg to registry server
		execInfo := new(core.ExecutorInfo)
		execInfo.TaskConfig = taskInfo
		execInfo.Worker = n.NodeInfo()
		_, err := n.Registry.Set(execInfo.GetExecutorKey(n.ClusterId()), execInfo.Json(), nil)
		if err != nil {
			logger.Node().Warn(logTitle + "sync to registry error:" + err.Error())
		}
		return core.SuccessResult()
	}
}

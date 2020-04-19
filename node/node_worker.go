package node

import (
	"github.com/devfeel/rockman/core"
	"github.com/devfeel/rockman/logger"
)

// RegisterExecutor
func (n *Node) RegisterExecutor(taskInfo *core.TaskConfig) *core.Result {
	lt := "Node RegisterExecutor [" + taskInfo.TaskID + "] "
	if !n.IsWorker() {
		logger.Node().Debug(lt + "failed, current node is not worker.")
		return core.FailedResult(-1001, "current node is not worker")
	}
	_, err := n.Runtime.CreateExecutor(taskInfo)
	if err != nil {
		logger.Node().Warn(lt + "CreateExecutor error:" + err.Error())
		return core.FailedResult(-2001, err.Error())
	} else {
		// update node info
		nodeInfo := n.refreshNodeInfo()
		// reg to registry server
		_, err := n.Registry.Set(nodeInfo.GetNodeKey(n.ClusterId()), nodeInfo.Json(), nil)
		if err != nil {
			logger.Node().Warn(lt + "Registry Set error:" + err.Error())
		}

		return core.SuccessResult()
	}
}

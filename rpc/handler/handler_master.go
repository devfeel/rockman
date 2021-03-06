package handler

import (
	"github.com/devfeel/rockman/core"
	"github.com/devfeel/rockman/logger"
	"github.com/devfeel/rockman/rpc/packet"
)

// RegisterWorker register worker node to leader
// it will check cluster id
func (h *RpcHandler) RegisterNode(nodeInfo *core.NodeInfo, reply *packet.RpcReply) error {
	logTitle := "RpcServer.RegisterNode[" + nodeInfo.EndPoint() + "] "
	if !h.getNode().IsLeader() {
		logger.Rpc().Warn(logTitle + "can not register to not leader node")
		*reply = packet.FailedReply(-1001, "can not register to not leader node")
		return nil
	}

	result := h.getNode().Cluster.AddNodeInfo(nodeInfo)
	if result.Error != nil {
		logger.Rpc().DebugS(logTitle+"error:", result.Error.Error())
		*reply = packet.FailedReply(-2001, result.Message())
	} else {
		if !result.IsSuccess() {
			logger.Rpc().DebugS(logTitle+"failed, ", result.Message())
			*reply = packet.FailedReply(-2002, result.Message())
		} else {
			logger.Rpc().DebugS(logTitle+"success,", nodeInfo.Json())
			*reply = packet.SuccessRpcReply(len(h.getNode().Runtime.Executors))
		}
	}
	return nil
}

// QueryNodes query node list from leader
func (h *RpcHandler) QueryNodes(pageInfo *core.PageInfo, reply *packet.RpcReply) error {
	logTitle := "RpcServer.QueryNodes "
	if !h.getNode().IsLeader() {
		logger.Rpc().Warn(logTitle + "can not query nodes from not leader node")
		*reply = packet.FailedReply(-1001, "can not query nodes from not leader node")
		return nil
	}

	logger.Rpc().DebugS(logTitle + "success")
	*reply = packet.SuccessRpcReply(h.getNode().Cluster.Nodes)
	return nil
}

// QueryClusterExecutorInfos return executor info in cluster by taskId
func (h *RpcHandler) QueryClusterExecutors(taskId string, reply *packet.RpcReply) error {
	logTitle := "RpcServer.QueryClusterExecutors [" + taskId + "] "
	if !h.getNode().IsLeader() {
		logger.Rpc().Warn(logTitle + "can not query executor info from not leader node")
		*reply = packet.FailedReply(-1001, "can not query executor info from not leader node")
		return nil
	}
	if taskId != "" {
		execInfo, exists := h.getNode().Cluster.FindExecutor(taskId)
		if !exists {
			logger.Rpc().Debug(logTitle + "not exists this taskId")
			*reply = packet.FailedReply(-2001, "not exists this taskId")
			return nil
		}
		configs := make(map[string]*core.ExecutorInfo)
		configs[taskId] = execInfo
		*reply = packet.SuccessRpcReply(configs)
	} else {
		*reply = packet.SuccessRpcReply(h.getNode().Cluster.ExecutorInfos)
	}

	logger.Rpc().Debug(logTitle + "success")
	return nil

}

// SubmitExecutor submit executor to leader node, then register to worker node
// it will check cluster id
func (h *RpcHandler) SubmitExecutor(execInfo *core.ExecutorInfo, reply *packet.RpcReply) error {
	logTitle := "RpcServer.SubmitExecutor: "
	if !h.getNode().IsLeader() {
		logger.Rpc().Warn("can not submit to not leader node")
		*reply = packet.FailedReply(-1001, "can not submit to not leader node")
		return nil
	}

	//async send executor to worker node
	result := h.getNode().SubmitExecutor(execInfo)
	if result.Error != nil {
		logger.Rpc().DebugS(logTitle+"error:", result.Error.Error())
		*reply = packet.FailedReply(-2001, result.Message())
	} else {
		if !result.IsSuccess() {
			logger.Rpc().DebugS(logTitle + "failed, " + result.Message())
			*reply = packet.FailedReply(-2002, result.Message())
		} else {
			logger.Rpc().DebugS(logTitle+"success", execInfo.TaskConfig.TaskID)
			*reply = packet.SuccessRpcReply(len(h.getNode().Runtime.Executors))
		}
	}
	return nil
}

// SubmitStopExecutor
func (h *RpcHandler) SubmitStopExecutor(taskId string, reply *packet.RpcReply) error {
	logTitle := "RpcServer.SubmitStopExecutor: "
	if !h.getNode().IsLeader() {
		logger.Rpc().Warn("can not submit to not leader node")
		*reply = packet.FailedReply(-1001, "can not submit to not leader node")
		return nil
	}

	//async send executor to worker node
	result := h.getNode().SubmitStopExecutor(taskId)
	if result.Error != nil {
		logger.Rpc().DebugS(logTitle+"error:", result.Error.Error())
		*reply = packet.FailedReply(-2001, result.Message())
	} else {
		if !result.IsSuccess() {
			logger.Rpc().DebugS(logTitle + "failed, " + result.Message())
			*reply = packet.FailedReply(-2002, result.Message())
		} else {
			logger.Rpc().DebugS(logTitle+"success", taskId)
			*reply = packet.SuccessRpcReply(len(h.getNode().Runtime.Executors))
		}
	}
	return nil
}

// SubmitStartExecutor
func (h *RpcHandler) SubmitStartExecutor(taskId string, reply *packet.RpcReply) error {
	lt := "RpcServer.SubmitStartExecutor: "
	if !h.getNode().IsLeader() {
		logger.Rpc().Warn("can not submit to not leader node")
		*reply = packet.FailedReply(-1001, "can not submit to not leader node")
		return nil
	}

	//async send executor to worker node
	result := h.getNode().SubmitStartExecutor(taskId)
	if result.Error != nil {
		logger.Rpc().DebugS(lt+"error:", result.Error.Error())
		*reply = packet.FailedReply(-2001, result.Message())
	} else {
		if !result.IsSuccess() {
			logger.Rpc().DebugS(lt + "failed, " + result.Message())
			*reply = packet.FailedReply(-2002, result.Message())
		} else {
			logger.Rpc().DebugS(lt+"success", taskId)
			*reply = packet.SuccessRpcReply(len(h.getNode().Runtime.Executors))
		}
	}
	return nil
}

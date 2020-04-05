package handler

import (
	"github.com/devfeel/rockman/core"
	"github.com/devfeel/rockman/logger"
	"github.com/devfeel/rockman/rpc/packet"
)

// RegisterWorker register worker node to leader node
// it will check cluster id
func (h *RpcHandler) RegisterNode(nodeInfo *core.NodeInfo, reply *packet.RpcReply) error {
	logTitle := "RpcServer.RegisterNode[" + nodeInfo.EndPoint() + "] "
	if !h.getNode().IsLeader() {
		logger.Default().Warn(logTitle + "can not register to not leader node")
		*reply = packet.CreateFailedReply(-1001, "can not register to not leader node")
		return nil
	}

	result := h.getNode().Cluster.AddNodeInfo(nodeInfo)
	if result.Error != nil {
		logger.Default().DebugS(logTitle+"error:", result.Error.Error())
		*reply = packet.CreateFailedReply(-2001, result.Message())
	} else {
		if !result.IsSuccess() {
			logger.Default().DebugS(logTitle+"failed, ", result.Message())
			*reply = packet.CreateFailedReply(-2002, result.Message())
		} else {
			logger.Default().DebugS(logTitle+"success,", nodeInfo.Json())
			*reply = packet.CreateSuccessRpcReply(len(h.getNode().Runtime.Executors))
		}
	}
	return nil
}

// QueryNodes query node list from leader node
func (h *RpcHandler) QueryNodes(pageInfo *core.PageInfo, reply *packet.RpcReply) error {
	logTitle := "RpcServer.QueryNodes "
	if !h.getNode().IsLeader() {
		logger.Default().Warn(logTitle + "can not query nodes from not leader node")
		*reply = packet.CreateFailedReply(-1001, "can not query nodes from not leader node")
		return nil
	}

	logger.Default().DebugS(logTitle + "success")
	*reply = packet.CreateSuccessRpcReply(h.getNode().Cluster.Nodes)
	return nil
}

// SubmitExecutor submit executor to leader node, then register to worker node
// it will check cluster id
func (h *RpcHandler) SubmitExecutor(submit *core.SubmitInfo, reply *packet.RpcReply) error {
	logTitle := "RpcServer.SubmitExecutor: "
	if !h.getNode().IsLeader() {
		logger.Default().Warn("can not submit to not leader node")
		*reply = packet.CreateFailedReply(-1001, "can not submit to not leader node")
		return nil
	}

	//async send executor to worker node
	result := h.getNode().SubmitExecutor(submit)
	if result.Error != nil {
		logger.Default().DebugS(logTitle+"error:", result.Error.Error())
		*reply = packet.CreateFailedReply(-2001, result.Message())
	} else {
		if !result.IsSuccess() {
			logger.Default().DebugS(logTitle + "failed, " + result.Message())
			*reply = packet.CreateFailedReply(-2002, result.Message())
		} else {
			logger.Default().DebugS(logTitle+"success", submit.TaskConfig.TaskID)
			*reply = packet.CreateSuccessRpcReply(len(h.getNode().Runtime.Executors))
		}
	}
	return nil
}

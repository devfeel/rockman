package handler

import (
	"github.com/devfeel/rockman/core"
	"github.com/devfeel/rockman/logger"
)

// RegisterWorker register worker node to leader node
func (h *RpcHandler) RegisterNode(nodeInfo *core.NodeInfo, result *core.JsonResult) error {
	logTitle := "RpcServer.RegisterNode[" + nodeInfo.EndPoint() + "] "
	if !h.getNode().IsLeader() {
		logger.Default().Warn(logTitle + "can not register to not leader node")
		*result = createResult(-1001, "can not register to not leader node", nil)
		return nil
	}

	h.getNode().Cluster.AddNodeInfo(nodeInfo)
	logger.Default().DebugS(logTitle + "success")
	*result = createResult(0, "ok", h.getNode().Cluster.Nodes)
	return nil
}

// QueryNodes query node list from leader node
func (h *RpcHandler) QueryNodes(pageInfo *core.PageInfo, result *core.JsonResult) error {
	logTitle := "RpcServer.QueryNodes "
	if !h.getNode().IsLeader() {
		logger.Default().Warn(logTitle + "can not query nodes from not leader node")
		*result = createResult(-1001, "can not query nodes from not leader node", nil)
		return nil
	}

	logger.Default().DebugS(logTitle + "success")
	*result = createResult(0, "ok", h.getNode().Cluster.Nodes)
	return nil
}

// SubmitExecutor submit executor to leader node, then register to worker node
func (h *RpcHandler) SubmitExecutor(submit *core.SubmitInfo, result *core.JsonResult) error {
	logTitle := "RpcServer.SubmitExecutor: "
	if !h.getNode().IsLeader() {
		logger.Default().Warn("can not submit to not leader node")
		*result = createResult(-1001, "can not submit to not leader node", nil)
		return nil
	}

	//async send executor to worker node
	err, reply := h.getNode().SubmitExecutor(submit)
	if err != nil {
		logger.Default().DebugS(logTitle+"error:", err.Error())
		*result = createResult(-2001, "submit error", nil)
	} else {
		if reply.RetCode != reply.CorrectCode() {
			logger.Default().DebugS(logTitle+"failed, reply code:", reply.RetCode)
			*result = createResult(reply.RetCode, reply.RetMsg, reply.Message)
		} else {
			logger.Default().DebugS(logTitle+"success", submit)
			*result = createResult(reply.RetCode, reply.RetMsg, h.getNode().Runtime.Executors)
		}
	}
	return nil
}

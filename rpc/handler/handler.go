package handler

import (
	"github.com/devfeel/rockman/core"
	"github.com/devfeel/rockman/logger"
	"github.com/devfeel/rockman/node"
	"strconv"
)

type RpcHandler struct {
	node *node.Node
}

func NewRpcHandler(node *node.Node) *RpcHandler {
	return &RpcHandler{node: node}
}

// Echo
func (h *RpcHandler) Echo(content string, result *string) error {
	logger.Default().Debug("RpcServer.Echo:" + content)
	*result = content
	return nil
}

// RegisterWorker register worker node to leader node
func (h *RpcHandler) RegisterNode(nodeInfo *core.NodeInfo, result *core.JsonResult) error {
	logTitle := "RpcServer.RegisterNode[" + nodeInfo.EndPoint() + "] "
	if !h.getNode().IsLeader {
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
	if !h.getNode().IsLeader {
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
	if !h.getNode().IsLeader {
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

// RegisterExecutor register executor to runtime in worker node
func (h *RpcHandler) RegisterExecutor(config *core.TaskConfig, result *core.JsonResult) error {
	logTitle := "RpcServer.RegisterExecutor: "
	if !h.getNode().Config.IsWorker {
		logger.Default().Warn("unworker node can not register executor")
		*result = core.JsonResult{RetCode: -1001, RetMsg: "unworker node can not register executor"}
		return nil
	}

	exec, err := h.getNode().Runtime.CreateExecutor(config)
	if err != nil {
		logger.Default().Warn(logTitle + "CreateExecutor error:" + err.Error())
		*result = core.JsonResult{RetCode: -9001, RetMsg: "CreateExecutor error:" + err.Error()}
		return nil
	} else {
		if exec.GetTaskConfig().IsRun {
			exec.GetTask().Start()
		}
	}
	logger.Default().DebugS(logTitle+"success", config)
	*result = core.JsonResult{RetCode: 0, RetMsg: "ok", Message: h.getNode().Runtime.Executors}
	return nil
}

// StartExecutor start executor by taskId
func (h *RpcHandler) StartExecutor(taskId string, result *core.JsonResult) error {
	logTitle := "RpcServer.StartExecutor[" + taskId + "] "
	if !h.getNode().Config.IsWorker {
		logger.Default().Warn("unworker node can not start executor")
		*result = core.JsonResult{-1001, "unworker node can not start executor", nil}
		return nil
	}
	err := h.getNode().Runtime.StartExecutor(taskId)
	if err != nil {
		logger.Default().Debug(logTitle + "error:" + err.Error())
		logger.Default().Error(err, logTitle+"error")
		*result = core.JsonResult{-2001, err.Error(), nil}
	}
	logger.Default().Debug(logTitle + "success")
	*result = core.JsonResult{RetCode: 0, RetMsg: "ok", Message: nil}
	return nil
}

// StopExecutor stop executor by taskId
func (h *RpcHandler) StopExecutor(taskId string, result *core.JsonResult) error {
	logTitle := "RpcServer.StopExecutor[" + taskId + "] "
	if !h.getNode().Config.IsWorker {
		logger.Default().Warn(logTitle + "unworker node can not start executor")
		*result = core.JsonResult{-1001, "unworker node can not start executor", nil}
		return nil
	}
	if !h.getNode().Config.IsWorker {
		logger.Default().Warn(logTitle + "unworker node can not start executor")
		*result = core.JsonResult{-1001, "unworker node can not start executor", nil}
		return nil
	}

	err := h.getNode().Runtime.StopExecutor(taskId)
	if err != nil {
		logger.Default().Debug(logTitle + "error:" + err.Error())
		logger.Default().Error(err, logTitle+"error")
		*result = core.JsonResult{-2001, logTitle + "error:" + err.Error(), nil}
	}
	logger.Default().Debug(logTitle + "success")
	*result = core.JsonResult{RetCode: 0, RetMsg: "ok", Message: nil}
	return nil
}

// RemoveExecutor remove executor by taskId
// if task is running, auto stop it first
func (h *RpcHandler) RemoveExecutor(taskId string, result *core.JsonResult) error {
	logTitle := "RpcServer.RemoveExecutor[" + taskId + "] "
	if !h.getNode().Config.IsWorker {
		logger.Default().Warn(logTitle + "unworker node can not start executor")
		*result = core.JsonResult{-1001, "unworker node can not start executor", nil}
		return nil
	}
	err := h.getNode().Runtime.RemoveExecutor(taskId)
	if err != nil {
		logger.Default().Debug(logTitle + "error:" + err.Error())
		logger.Default().Error(err, logTitle+"error")
		*result = core.JsonResult{-2001, logTitle + "error:" + err.Error(), nil}
	}
	logger.Default().Debug(logTitle + "success")
	*result = core.JsonResult{RetCode: 0, RetMsg: "ok", Message: h.getNode().Runtime.Executors}
	return nil
}

// QueryExecutors return executors in runtime by taskId
// if taskId is nil, return all executors
func (h *RpcHandler) QueryExecutorConfig(taskId string, result *core.JsonResult) error {
	logTitle := "RpcServer.QueryExecutors [" + taskId + "] "
	if !h.getNode().Config.IsWorker {
		logger.Default().Warn(logTitle + "unworker node can not start executor")
		*result = core.JsonResult{-1001, "unworker node can not start executor", nil}
		return nil
	}
	configs := h.getNode().Runtime.QueryAllExecutorConfig()
	if taskId != "" {
		exec, isOk := configs[taskId]
		if !isOk {
			*result = core.JsonResult{-1001, "not exists this taskId", nil}
		} else {
			logger.Default().Debug(logTitle + "success")
			configs = make(map[string]core.TaskConfig)
			configs[taskId] = exec
			*result = core.JsonResult{0, "ok", configs}
		}
	} else {
		logger.Default().Debug(logTitle + "success, config count = " + strconv.Itoa(len(configs)))
		*result = core.JsonResult{RetCode: 0, RetMsg: "ok", Message: configs}
	}
	return nil
}

func (h *RpcHandler) getNode() *node.Node {
	return h.node
}

func createResult(retCode int, retMsg string, message interface{}) core.JsonResult {
	return core.JsonResult{
		RetCode: retCode,
		RetMsg:  retMsg,
		Message: message,
	}
}

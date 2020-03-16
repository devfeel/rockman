package handler

import (
	"github.com/devfeel/mapper"
	"github.com/devfeel/rockman/src/logger"
	"github.com/devfeel/rockman/src/node"
	"github.com/devfeel/rockman/src/packets"
	"github.com/devfeel/rockman/src/rpc/packet"
	"github.com/devfeel/rockman/src/runtime/executor"
)

type RpcHandler struct {
	node *node.Node
}

func NewRpcHandler(node *node.Node) *RpcHandler {
	return &RpcHandler{node: node}
}

// Echo
func (h *RpcHandler) Echo(content string, result *packet.JsonResult) error {
	logger.Default().Debug("RpcHandler:Echo:" + content)
	*result = packet.JsonResult{0, "ok", content}
	return nil
}

// RegisterWorker register worker node to leader node
func (h *RpcHandler) RegisterWorker(worker *packets.WorkerInfo, result *packet.JsonResult) error {
	if !h.getNode().IsLeader {
		*result = packet.JsonResult{-1001, "can not register to not leader node", nil}
		return nil
	}

	err := h.getNode().Cluster.AddWorker(worker)
	if err != nil {
		*result = packet.JsonResult{-9001, "can not add node to cluster:" + err.Error(), nil}
		return nil
	}

	*result = packet.JsonResult{0, "ok", h.getNode().Cluster.Workers}
	return nil
}

func (h *RpcHandler) RegisterExecutor(config interface{}, result *packet.JsonResult) error {
	if !h.getNode().Config.IsWorker {
		logger.Default().Warn("unworker node can not register executor")
		*result = packet.JsonResult{-1001, "unworker node can not register executor", nil}
		return nil
	}

	taskConfig := &executor.TaskConfig{}
	err := mapper.MapperMap(config.(map[string]interface{}), taskConfig)
	if err != nil {
		logger.Default().Error(err, "mapper config to TaskConfig error:"+err.Error())
		*result = packet.JsonResult{-2001, "mapper config to TaskConfig error:" + err.Error(), nil}
		return nil
	}
	realTaskConfig, err := executor.ConvertRealTaskConfig(taskConfig)
	if err != nil {
		logger.Default().Error(err, "convert real task config err:"+err.Error())
		*result = packet.JsonResult{-2002, "convert real task config err:" + err.Error(), nil}
		return nil
	}

	exec, err := h.getNode().Runtime.CreateExecutor(taskConfig.TaskID, taskConfig.TargetType, realTaskConfig)
	if err != nil {
		logger.Default().Error(err, "CreateExecutor error:"+err.Error())
		*result = packet.JsonResult{-9001, "CreateExecutor error:" + err.Error(), nil}
		return nil
	} else {
		if exec.GetTaskConfig().IsRun {
			exec.GetTask().Start()
		}
	}

	logger.Default().DebugS("RegisterExecutor success", config)
	*result = packet.JsonResult{0, "ok", nil}
	return nil
}

func (h *RpcHandler) getNode() *node.Node {
	return h.node
}
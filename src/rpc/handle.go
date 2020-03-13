package rpc

import (
	"github.com/devfeel/rockman/src/core/packets"
	"github.com/devfeel/rockman/src/logger"
	"github.com/devfeel/rockman/src/node"
	"github.com/devfeel/rockman/src/runtime/executor"
)

type RpcHandler struct {
	server *RpcServer
}

func NewRpcHandler(server *RpcServer) *RpcHandler {
	return &RpcHandler{server: server}
}

// Echo
func (h *RpcHandler) Echo(content string, result *JsonResult) error {
	logger.Default().Debug("RpcHandler:Echo:" + content)
	*result = JsonResult{0, "ok", content}
	return nil
}

// RegisterNode
func (h *RpcHandler) RegisterNode(nodeInfo packets.NodeInfo, result *JsonResult) error {
	if !h.getNode().Cluster.IsMaster {
		*result = JsonResult{-1001, "can not register to unmaster node", nil}
		return nil
	}
	if h.server.RpcHost == nodeInfo.Host && h.server.RpcPort == nodeInfo.Port {
		*result = JsonResult{-1002, "can not register node to self", nil}
		return nil
	}

	err := h.getNode().Cluster.AddNode(&nodeInfo)
	if err != nil {
		*result = JsonResult{-9001, "can not add node to cluster:" + err.Error(), nil}
		return nil
	}

	*result = JsonResult{0, "ok", h.getNode().Cluster.Workers}
	return nil
}

func (h *RpcHandler) RegisterExecutor(config interface{}, result *JsonResult) error {
	taskConfig, isOk := config.(executor.TaskConfig)
	if !isOk {
		*result = JsonResult{-1001, "Invalid config type", nil}
		return nil
	}
	_, err := h.getNode().Runtime.CreateExecutor(taskConfig.TaskID, taskConfig.TaskType, taskConfig)
	if err != nil {
		*result = JsonResult{-9001, "CreateExecutor error:" + err.Error(), nil}
		return nil
	}
	*result = JsonResult{0, "ok", nil}
	return nil
}

func (h *RpcHandler) getNode() *node.Node {
	return h.server.Node
}

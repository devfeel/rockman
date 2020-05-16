package handler

import (
	"github.com/devfeel/rockman/core"
	"github.com/devfeel/rockman/logger"
	"github.com/devfeel/rockman/node"
	"github.com/devfeel/rockman/rpc/packet"
	"github.com/devfeel/rockman/util/sysx"
	"strconv"
)

type RpcHandler struct {
	node *node.Node
}

func NewRpcHandler(node *node.Node) *RpcHandler {
	return &RpcHandler{node: node}
}

// Echo
func (h *RpcHandler) Echo(content string, reply *string) error {
	logger.Rpc().Debug("RpcServer.Echo:" + content)
	*reply = content
	return nil
}

// QueryResource query resource info from worker node
func (h *RpcHandler) QueryResource(content string, reply *packet.RpcReply) error {
	if !h.node.IsWorker() {
		logger.Rpc().Warn("QueryResource failed: can not query resource from not worker node")
		*reply = packet.FailedReply(-1001, "can not query resource from not worker nodee")
		return nil
	}
	resource := &core.ResourceInfo{}
	resource.EndPoint = h.node.NodeInfo().EndPoint()
	resource.TaskCount = h.node.Runtime.TaskService.Count()
	resource.CpuRate = sysx.GetCpuUsedPercent()
	resource.MemoryRate = sysx.GetMemoryUsedPercent()

	logger.Rpc().DebugS("RpcServer.QueryResource success", *resource)
	*reply = packet.SuccessRpcReply(resource)
	return nil
}

// QueryExecutors return executors in runtime by taskId
// if taskId is nil, return all executors
func (h *RpcHandler) QueryExecutors(taskId string, reply *packet.RpcReply) error {
	logTitle := "RpcServer.QueryExecutors [" + taskId + "] "
	if !h.getNode().IsWorker() {
		logger.Rpc().Warn(logTitle + "unworker node can not query executor")
		*reply = packet.FailedReply(-1001, "unworker node can not query executor")
		return nil
	}
	configs := h.getNode().Runtime.QueryAllExecutorConfig()
	if taskId != "" {
		exec, isOk := configs[taskId]
		if !isOk {
			*reply = packet.FailedReply(-2001, "not exists this taskId")
			return nil
		} else {
			logger.Rpc().Debug(logTitle + "success")
			configs = make(map[string]core.TaskConfig)
			configs[taskId] = exec
		}
	} else {
		logger.Rpc().Debug(logTitle + "success, config count = " + strconv.Itoa(len(configs)))
	}
	execInfos := make(map[string]*core.ExecutorInfo)
	for k, v := range configs {
		execInfos[k] = &core.ExecutorInfo{
			TaskConfig: &v,
			Worker:     h.getNode().NodeInfo(),
		}
	}
	*reply = packet.SuccessRpcReply(execInfos)
	return nil
}

func (h *RpcHandler) getNode() *node.Node {
	return h.node
}

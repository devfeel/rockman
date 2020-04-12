package handler

import (
	"github.com/devfeel/rockman/core"
	"github.com/devfeel/rockman/logger"
	"github.com/devfeel/rockman/rpc/packet"
)

// RegisterExecutor register executor to runtime in worker node
func (h *RpcHandler) RegisterExecutor(config *core.TaskConfig, reply *packet.RpcReply) error {
	lt := "RpcServer.RegisterExecutor: "
	if !h.getNode().IsWorker() {
		logger.Rpc().Warn("unworker node can not register executor")
		*reply = packet.FailedReply(-1001, "unworker node can not register executor")
		return nil
	}
	result := h.getNode().RegisterExecutor(config)
	if result.Error != nil {
		logger.Rpc().Warn(lt + "error:" + result.Error.Error())
		*reply = packet.FailedReply(-9001, "error:"+result.Error.Error())
		return nil
	}
	if !result.IsSuccess() {
		logger.Rpc().Warn(lt + "failed:" + result.Message())
		*reply = packet.FailedReply(-9001, "failed:"+result.Message())
		return nil
	}
	logger.Rpc().DebugS(lt+"success", config)
	*reply = packet.SuccessRpcReply(h.getNode().Runtime.Executors)
	return nil
}

// StartExecutor start executor by taskId
func (h *RpcHandler) StartExecutor(taskId string, reply *packet.RpcReply) error {
	logTitle := "RpcServer.StartExecutor[" + taskId + "] "
	if !h.getNode().IsWorker() {
		logger.Rpc().Warn("unworker node can not start executor")
		*reply = packet.FailedReply(-1001, "unworker node can not start executor")
		return nil
	}
	err := h.getNode().Runtime.StartExecutor(taskId)
	if err != nil {
		logger.Rpc().Debug(logTitle + "error:" + err.Error())
		logger.Rpc().Error(err, logTitle+"error")
		*reply = packet.FailedReply(-2001, err.Error())
	}
	logger.Rpc().Debug(logTitle + "success")
	*reply = packet.SuccessRpcReply(nil)
	return nil
}

// StopExecutor stop executor by taskId
func (h *RpcHandler) StopExecutor(taskId string, reply *packet.RpcReply) error {
	logTitle := "RpcServer.StopExecutor[" + taskId + "] "
	if !h.getNode().IsWorker() {
		logger.Rpc().Warn(logTitle + "unworker node can not stop executor")
		*reply = packet.FailedReply(-1001, "unworker node can not stop executor")
		return nil
	}

	err := h.getNode().Runtime.StopExecutor(taskId)
	if err != nil {
		logger.Rpc().Debug(logTitle + "error:" + err.Error())
		logger.Rpc().Error(err, logTitle+"error")
		*reply = packet.FailedReply(-2001, logTitle+"error:"+err.Error())
	}
	logger.Rpc().Debug(logTitle + "success")
	*reply = packet.SuccessRpcReply(nil)
	return nil
}

// RemoveExecutor remove executor by taskId
// if task is running, auto stop it first
func (h *RpcHandler) RemoveExecutor(taskId string, reply *packet.RpcReply) error {
	logTitle := "RpcServer.RemoveExecutor[" + taskId + "] "
	if !h.getNode().IsWorker() {
		logger.Rpc().Warn(logTitle + "unworker node can not remove executor")
		*reply = packet.FailedReply(-1001, "unworker node can not remove executor")
		return nil
	}
	err := h.getNode().Runtime.RemoveExecutor(taskId)
	if err != nil {
		logger.Rpc().Debug(logTitle + "error:" + err.Error())
		logger.Rpc().Error(err, logTitle+"error")
		*reply = packet.FailedReply(-2001, logTitle+"error:"+err.Error())
	}
	logger.Rpc().Debug(logTitle + "success")
	*reply = packet.SuccessRpcReply(h.getNode().Runtime.Executors)
	return nil
}

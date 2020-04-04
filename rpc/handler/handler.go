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

// RegisterExecutor register executor to runtime in worker node
func (h *RpcHandler) RegisterExecutor(config *core.TaskConfig, result *core.JsonResult) error {
	logTitle := "RpcServer.RegisterExecutor: "
	if !h.getNode().Config().Node.IsWorker {
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
	if !h.getNode().IsWorker() {
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
	if !h.getNode().IsWorker() {
		logger.Default().Warn(logTitle + "unworker node can not start executor")
		*result = core.JsonResult{-1001, "unworker node can not start executor", nil}
		return nil
	}
	if !h.getNode().IsWorker() {
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
	if !h.getNode().IsWorker() {
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
	if !h.getNode().IsWorker() {
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

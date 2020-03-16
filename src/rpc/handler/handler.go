package handler

import (
	"github.com/devfeel/mapper"
	"github.com/devfeel/rockman/src/logger"
	"github.com/devfeel/rockman/src/node"
	"github.com/devfeel/rockman/src/packets"
	"github.com/devfeel/rockman/src/runtime/executor"
)

type RpcHandler struct {
	node *node.Node
}

func NewRpcHandler(node *node.Node) *RpcHandler {
	return &RpcHandler{node: node}
}

// Echo
func (h *RpcHandler) Echo(content string, result *packets.JsonResult) error {
	logger.Default().Debug("RpcHandler:Echo:" + content)
	result = &packets.JsonResult{0, "ok", content}
	return nil
}

// RegisterWorker register worker node to leader node
func (h *RpcHandler) RegisterWorker(worker *packets.WorkerInfo, result *packets.JsonResult) error {
	if !h.getNode().IsLeader {
		result = &packets.JsonResult{-1001, "can not register to not leader node", nil}
		return nil
	}

	err := h.getNode().Cluster.AddWorker(worker)
	if err != nil {
		result = &packets.JsonResult{-9001, "can not add node to cluster:" + err.Error(), nil}
		return nil
	}

	result = &packets.JsonResult{0, "ok", h.getNode().Cluster.Workers}
	return nil
}

func (h *RpcHandler) RegisterExecutor(config interface{}, result *packets.JsonResult) error {
	if !h.getNode().Config.IsWorker {
		logger.Default().Warn("unworker node can not register executor")
		result = &packets.JsonResult{-1001, "unworker node can not register executor", nil}
		return nil
	}

	taskConfig := &executor.TaskConfig{}
	err := mapper.MapperMap(config.(map[string]interface{}), taskConfig)
	if err != nil {
		logger.Default().Error(err, "mapper config to TaskConfig error:"+err.Error())
		result = &packets.JsonResult{-2001, "mapper config to TaskConfig error:" + err.Error(), nil}
		return nil
	}
	realTaskConfig, err := executor.ConvertRealTaskConfig(taskConfig)
	if err != nil {
		logger.Default().Error(err, "convert real task config err:"+err.Error())
		result = &packets.JsonResult{-2002, "convert real task config err:" + err.Error(), nil}
		return nil
	}

	exec, err := h.getNode().Runtime.CreateExecutor(taskConfig.TaskID, taskConfig.TargetType, realTaskConfig)
	if err != nil {
		logger.Default().Error(err, "CreateExecutor error:"+err.Error())
		result = &packets.JsonResult{-9001, "CreateExecutor error:" + err.Error(), nil}
		return nil
	} else {
		if exec.GetTaskConfig().IsRun {
			exec.GetTask().Start()
		}
	}

	logger.Default().DebugS("RegisterExecutor success", config)
	result = &packets.JsonResult{0, "ok", h.getNode().Runtime.Executors}
	return nil
}

func (h *RpcHandler) StartExecutor(taskId string, result *packets.JsonResult) error {
	logTitle := "StartExecutor[" + taskId + "] "
	err := h.getNode().Runtime.StartExecutor(taskId)
	if err != nil {
		logger.Default().Debug(logTitle + "error:" + err.Error())
		logger.Default().Error(err, logTitle+"error")
		result = &packets.JsonResult{-2001, logTitle + "error:" + err.Error(), nil}
	}
	logger.Default().Debug(logTitle + "success")
	result = &packets.JsonResult{0, "ok", nil}
	return nil
}

func (h *RpcHandler) StopExecutor(taskId string, result *packets.JsonResult) error {
	logTitle := "StopExecutor[" + taskId + "] "
	err := h.getNode().Runtime.StopExecutor(taskId)
	if err != nil {
		logger.Default().Debug(logTitle + "error:" + err.Error())
		logger.Default().Error(err, logTitle+"error")
		result = &packets.JsonResult{-2001, logTitle + "error:" + err.Error(), nil}
	}
	logger.Default().Debug(logTitle + "success")
	result = &packets.JsonResult{0, "ok", nil}
	return nil
}

func (h *RpcHandler) RemoveExecutor(taskId string, result *packets.JsonResult) error {
	logTitle := "RemoveExecutor[" + taskId + "] "
	err := h.getNode().Runtime.RemoveExecutor(taskId)
	if err != nil {
		logger.Default().Debug(logTitle + "error:" + err.Error())
		logger.Default().Error(err, logTitle+"error")
		result = &packets.JsonResult{-2001, logTitle + "error:" + err.Error(), nil}
	}
	logger.Default().Debug(logTitle + "success")
	result = &packets.JsonResult{0, "ok", h.getNode().Runtime.Executors}
	return nil
}

func (h *RpcHandler) getNode() *node.Node {
	return h.node
}

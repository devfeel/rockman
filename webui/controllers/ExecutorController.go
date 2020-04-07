package controllers

import (
	"github.com/devfeel/dotweb"
	"github.com/devfeel/rockman-webui/src/protected/viewModel/task"
	"github.com/devfeel/rockman/core"
	"github.com/devfeel/rockman/protected/model"
	"github.com/devfeel/rockman/protected/service/executor"
	_const "github.com/devfeel/rockman/webui/const"
)

type ExecutorController struct {
	executorService *executor.ExecutorService
}

func NewExecutorController() *ExecutorController {
	return &ExecutorController{
		executorService: executor.NewExecutorService(),
	}
}

func (c *ExecutorController) SaveExecutor(ctx dotweb.Context) error {
	model := &model.ExecutorInfo{}
	err := ctx.Bind(model)
	if err != nil {
		return ctx.WriteJson(FailedResponse(-1001, "parameter bind failed: "+err.Error()))
	}

	var result *core.Result
	if model.ID > 0 {
		result = c.executorService.UpdateExecutor(model)
	} else {
		result = c.executorService.AddExecutor(model)
	}

	if !result.IsSuccess() {
		return ctx.WriteJson(FailedResponse(-2001, "save failed: "+result.Message()))
	}

	return ctx.WriteJson(SuccessResponse(nil))
}

// QueryById
func (c *ExecutorController) QueryById(ctx dotweb.Context) error {
	model := task.ExecutorInfo{}
	err := ctx.Bind(&model)
	if err != nil {
		return ctx.WriteJson(FailedResponse(-1001, "parameter bind failed: "+err.Error()))
	}
	result, err := c.executorService.QueryExecutorById(model.ID)
	if err != nil {
		return ctx.WriteJson(FailedResponse(-2001, "query failed: "+err.Error()))
	}
	return ctx.WriteJson(SuccessResponse(result))
}

// ShowExecutors
func (c *ExecutorController) ShowExecutors(ctx dotweb.Context) error {
	nodeId := ctx.QueryString("node")
	pageIndex := ctx.QueryInt64("pageindex")
	pageSize := ctx.QueryInt64("pagesize")
	pageReq := new(model.PageRequest)
	pageReq.PageIndex = pageIndex
	pageReq.PageSize = pageSize

	if pageReq.PageSize <= 0 {
		pageReq.PageSize = _const.DefaultPageSize
	}
	taskService := executor.NewExecutorService()
	result, err := taskService.QueryExecutors(nodeId, pageReq)
	if err != nil {
		return ctx.WriteJson(FailedResponse(-2001, "Query error: "+err.Error()))
	}
	return ctx.WriteJson(SuccessResponse(result))
}

// ShowExecLogs
func (c *ExecutorController) ShowExecLogs(ctx dotweb.Context) error {
	taskId := ctx.QueryString("task")
	pageIndex := ctx.QueryInt64("pageindex")
	pageSize := ctx.QueryInt64("pagesize")
	pageReq := new(model.PageRequest)
	pageReq.PageIndex = pageIndex
	pageReq.PageSize = pageSize

	if pageReq.PageSize <= 0 {
		pageReq.PageSize = _const.DefaultPageSize
	}

	taskService := executor.NewExecutorService()
	result, err := taskService.QueryExecLogs(taskId, pageReq)
	if err != nil {
		return ctx.WriteJson(FailedResponse(-2001, "Query error: "+err.Error()))
	}
	return ctx.WriteJson(SuccessResponse(result))
}

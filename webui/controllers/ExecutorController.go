package controllers

import (
	"github.com/devfeel/dotweb"
	"github.com/devfeel/rockman/protected/model"
	service "github.com/devfeel/rockman/protected/service"
	_const "github.com/devfeel/rockman/webui/const"
	"github.com/devfeel/rockman/webui/contract"
)

type ExecutorController struct {
}

func (c *ExecutorController) ShowExecutors(ctx dotweb.Context) error {
	taskService := service.NewExecutorService()
	result, err := taskService.QueryExecutors()
	if err != nil {
		return ctx.WriteJson(contract.CreateResponse(-2001, "Query Error", err))
	}
	return ctx.WriteJson(contract.CreateResponse(0, "", result))
}

// ShowExecLogs
func (c *ExecutorController) ShowExecLogs(ctx dotweb.Context) error {
	taskId := ctx.QueryString("taskid")
	pageIndex := ctx.QueryInt64("pageindex")
	pageSize := ctx.QueryInt64("pagesize")
	pageReq := new(model.PageRequest)
	pageReq.PageIndex = pageIndex
	pageReq.PageSize = pageSize

	if pageReq.PageSize <= 0 {
		pageReq.PageSize = _const.DefaultPageSize
	}

	taskService := service.NewExecutorService()
	result, err := taskService.QueryExecLogs(taskId, pageReq)
	if err != nil {
		return ctx.WriteJson(contract.CreateResponse(-2001, "Query Error:"+err.Error(), err))
	}
	return ctx.WriteJson(contract.CreateResponse(0, "", result))
}

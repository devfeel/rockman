package controllers

import (
	"github.com/devfeel/dotweb"
	"github.com/devfeel/rockman/protected/model"
	service "github.com/devfeel/rockman/protected/service"
	_const "github.com/devfeel/rockman/webui/const"
	"github.com/devfeel/rockman/webui/contract"
)

type TaskController struct {
}

func (c *TaskController) ShowTasks(ctx dotweb.Context) error {
	taskService := service.NewTaskService()
	result, err := taskService.QueryTasks()
	if err != nil {
		return ctx.WriteJson(contract.CreateResponse(-2001, "Query Error", err))
	}
	return ctx.WriteJson(contract.CreateResponse(0, "", result))
}

// ShowExecLogs
func (c *TaskController) ShowExecLogs(ctx dotweb.Context) error {
	taskId := ctx.QueryString("taskid")
	pageIndex := ctx.QueryInt64("pageindex")
	pageSize := ctx.QueryInt64("pagesize")
	pageReq := new(model.PageRequest)
	pageReq.PageIndex = pageIndex
	pageReq.PageSize = pageSize

	if pageReq.PageSize <= 0 {
		pageReq.PageSize = _const.DefaultPageSize
	}

	taskService := service.NewTaskService()
	result, err := taskService.QueryExecLogs(taskId, pageReq)
	if err != nil {
		return ctx.WriteJson(contract.CreateResponse(-2001, "Query Error:"+err.Error(), err))
	}
	return ctx.WriteJson(contract.CreateResponse(0, "", result))
}

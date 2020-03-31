package controllers

import (
	"github.com/devfeel/dotweb"
	service "github.com/devfeel/rockman/protected/service"
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

func (c *TaskController) ShowLogs(ctx dotweb.Context) error {
	taskService := service.NewTaskService()
	result, err := taskService.QueryLogs()
	if err != nil {
		return ctx.WriteJson(contract.CreateResponse(-2001, "Query Error", err))
	}
	return ctx.WriteJson(contract.CreateResponse(0, "", result))
}

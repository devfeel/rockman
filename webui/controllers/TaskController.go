package controllers

import (
	"github.com/devfeel/dotweb"
	service2 "github.com/devfeel/rockman/protected/service"
	"github.com/devfeel/rockman/webui/contract"
)

type TaskController struct {
}

func (c *TaskController) ShowTasks(ctx dotweb.Context) error {
	service := service2.NewTaskService()
	tasks, err := service.QueryTasks()
	if err != nil {
		return ctx.WriteJson(contract.CreateResponse(-2001, "Query Error", err))
	}
	return ctx.WriteJson(contract.CreateResponse(0, "", tasks))
}

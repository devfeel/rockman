package controllers

import (
	"github.com/devfeel/dotweb"
	"github.com/devfeel/rockman/protected/service/task"
	"github.com/devfeel/rockman/webui/contract"
	"github.com/devfeel/rockman/webui/validate"
)

type TaskController struct {
}

func (c *TaskController) ShowTaskListByNodeID(ctx dotweb.Context) error {
	nodeID := ctx.QueryString("NodeID")
	if rep, err := validate.IsNilString(nodeID, -1001, "NodeID is null"); err != nil {
		return ctx.WriteJson(rep)
	}

	service := task.NewTaskService()
	tasks, err := service.QueryTasksByNodeID(nodeID)
	if err != nil {
		return ctx.WriteJson(contract.CreateResponse(-2001, "Query Error", err))
	}
	return ctx.WriteJson(contract.CreateResponse(0, "", tasks))
}

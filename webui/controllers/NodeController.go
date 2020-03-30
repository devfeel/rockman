package controllers

import (
	"github.com/devfeel/dotweb"
	"github.com/devfeel/rockman/node"
	_const "github.com/devfeel/rockman/webui/const"
	"github.com/devfeel/rockman/webui/contract"
)

type NodeController struct {
}

func (c *NodeController) ShowNodeList(ctx dotweb.Context) error {
	item, isExists := ctx.AppItems().Get(_const.ItemKey_Node)
	if !isExists {
		return ctx.WriteJson(contract.CreateResponse(-1001, "not exists node in app items", nil))
	}
	node, isOk := item.(*node.Node)
	if !isOk {
		return ctx.WriteJson(contract.CreateResponse(-1002, "not exists correct node in app items", nil))
	}
	return ctx.WriteJson(contract.CreateResponse(0, "", node.Cluster.Nodes))
}

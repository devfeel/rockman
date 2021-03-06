package controllers

import (
	"github.com/devfeel/dotweb"
	"github.com/devfeel/rockman/core"
	"github.com/devfeel/rockman/node"
	_const "github.com/devfeel/rockman/webui/const"
)

type NodeController struct {
}

func (c *NodeController) ShowNodes(ctx dotweb.Context) error {
	item, isExists := ctx.AppItems().Get(_const.ItemKeyNode)
	if !isExists {
		return ctx.WriteJson(NewResponse(-1001, "not exists node in app items", nil))
	}
	node, isOk := item.(*node.Node)
	if !isOk {
		return ctx.WriteJson(NewResponse(-1002, "not exists correct node in app items", nil))
	}
	var nodes []*core.NodeInfo
	for _, n := range node.Cluster.Nodes {
		nodes = append(nodes, n)
	}
	return ctx.WriteJson(NewResponse(0, "", nodes))
}

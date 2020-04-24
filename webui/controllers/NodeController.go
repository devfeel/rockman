package controllers

import (
	"fmt"
	"github.com/devfeel/dotweb"
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
	fmt.Println(node)
	return ctx.WriteJson(NewResponse(0, "", node))
}

package controllers

import (
	"github.com/devfeel/dotweb"
	"github.com/devfeel/rockman/config"
	"github.com/devfeel/rockman/node"
	"github.com/devfeel/rockman/rpc/client"
	_const "github.com/devfeel/rockman/webui/const"
	"github.com/devfeel/rockman/webui/contract"
)

func getNode(ctx dotweb.Context) *node.Node {
	nodeItem, exists := ctx.AppItems().Get(_const.ItemKeyNode)
	if !exists {
		return nil
	}
	return nodeItem.(*node.Node)
}

func getLeader(ctx dotweb.Context) string {
	leader, _ := getNode(ctx).Cluster.GetLeaderInfo()
	return leader
}

func GetRpcClient(endPoint string) *client.RpcClient {
	config := config.GetProfile()
	return client.NewRpcClient(endPoint, config.Rpc.EnableTls, config.Rpc.ClientCertFile, config.Rpc.ClientKeyFile)
}

func NewResponse(retCode int, retMsg string, message interface{}) *contract.Response {
	return contract.NewResponse(retCode, retMsg, message)
}

func SuccessResponse(message interface{}) *contract.Response {
	return contract.SuccessResponse(message)
}

func FailedResponse(retCode int, retMsg string) *contract.Response {
	return contract.FailedResponse(retCode, retMsg)
}

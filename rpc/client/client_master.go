package client

import (
	"github.com/devfeel/rockman/core"
	"github.com/devfeel/rockman/logger"
	"github.com/devfeel/rockman/rpc/packet"
)

func (c *RpcClient) CallQueryExecutorInfos(taskId string) (error, *packet.RpcReply) {
	client, err := c.getConnClient()
	if err != nil {
		logger.Default().Error(err, "getConnClient error")
		return err, nil
	}
	var reply packet.RpcReply
	err = client.Call("Rpc.QueryExecutorInfos", taskId, &reply)
	if err != nil {
		return err, nil
	}
	return nil, &reply
}

func (c *RpcClient) CallNotifyExecutorChange(taskId string) (error, *packet.RpcReply) {
	client, err := c.getConnClient()
	if err != nil {
		logger.Default().Error(err, "getConnClient error")
		return err, nil
	}
	var reply packet.RpcReply
	err = client.Call("Rpc.NotifyExecutorChange", taskId, &reply)
	if err != nil {
		return err, nil
	}
	return nil, &reply
}

func (c *RpcClient) CallRegisterNode(worker *core.NodeInfo) (error, *packet.RpcReply) {
	client, err := c.getConnClient()
	if err != nil {
		logger.Default().Error(err, "getConnClient error")
		return err, nil
	}
	var reply packet.RpcReply
	err = client.Call("Rpc.RegisterNode", worker, &reply)
	if err != nil {
		return err, nil
	}
	return nil, &reply
}

func (c *RpcClient) CallQueryNodes(pageInfo *core.PageInfo) (error, *packet.RpcReply) {
	client, err := c.getConnClient()
	if err != nil {
		logger.Default().Error(err, "getConnClient error")
		return err, nil
	}
	var reply packet.RpcReply
	err = client.Call("Rpc.QueryNodes", pageInfo, &reply)
	if err != nil {
		return err, nil
	}
	return nil, &reply
}

func (c *RpcClient) CallSubmitExecutor(execInfo *core.ExecutorInfo) (error, *packet.RpcReply) {
	client, err := c.getConnClient()
	if err != nil {
		logger.Default().Error(err, "getConnClient error")
		return err, nil
	}
	var reply packet.RpcReply
	err = client.Call("Rpc.SubmitExecutor", execInfo, &reply)
	if err != nil {
		return err, nil
	}
	return nil, &reply
}

func (c *RpcClient) CallSubmitStartExecutor(taskId string) (error, *packet.RpcReply) {
	client, err := c.getConnClient()
	if err != nil {
		logger.Default().Error(err, "getConnClient error")
		return err, nil
	}
	var reply packet.RpcReply
	err = client.Call("Rpc.SubmitStartExecutor", taskId, &reply)
	if err != nil {
		return err, nil
	}
	return nil, &reply
}

func (c *RpcClient) CallSubmitStopExecutor(taskId string) (error, *packet.RpcReply) {
	client, err := c.getConnClient()
	if err != nil {
		logger.Default().Error(err, "getConnClient error")
		return err, nil
	}
	var reply packet.RpcReply
	err = client.Call("Rpc.SubmitStopExecutor", taskId, &reply)
	if err != nil {
		return err, nil
	}
	return nil, &reply
}

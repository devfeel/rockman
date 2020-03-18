package client

import (
	"fmt"
	"github.com/devfeel/rockman/logger"
	"github.com/devfeel/rockman/packets"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type RpcClient struct {
	serviceUrl string
	client     *rpc.Client
}

func NewRpcClient(serverUrl string) *RpcClient {
	return &RpcClient{serviceUrl: serverUrl}
}

// getConnClient
func (c *RpcClient) getConnClient() (*rpc.Client, error) {
	if c.client != nil {
		return c.client, nil
	}
	client, err := jsonrpc.Dial("tcp", c.serviceUrl)
	if err != nil {
		logger.Default().Error(err, "RpcClient connServer dial error:")
	} else {
		c.client = client
	}
	return client, err
}

func (c *RpcClient) CallEcho(message string) (error, string) {
	client, err := c.getConnClient()
	if err != nil {
		logger.Default().Error(err, "getConnClient error")
		return err, ""
	}
	reply := new(string)
	err = client.Call("Rpc.Echo", message, reply)
	fmt.Println(*reply)
	if err != nil {
		return err, ""
	}
	return nil, *reply
}

func (c *RpcClient) CallRegisterWorker(worker *packets.WorkerInfo) (error, *packets.JsonResult) {
	client, err := c.getConnClient()
	if err != nil {
		logger.Default().Error(err, "getConnClient error")
		return err, nil
	}
	var reply packets.JsonResult
	err = client.Call("Rpc.RegisterWorker", worker, &reply)
	if err != nil {
		return err, nil
	}
	return nil, &reply
}

func (c *RpcClient) CallRegisterExecutor(conf interface{}) (error, *packets.JsonResult) {
	client, err := c.getConnClient()
	if err != nil {
		logger.Default().Error(err, "getConnClient error")
		return err, nil
	}
	var reply packets.JsonResult
	err = client.Call("Rpc.RegisterExecutor", conf, &reply)
	if err != nil {
		return err, nil
	}
	return nil, &reply
}

func (c *RpcClient) CallStartExecutor(taskId string) (error, *packets.JsonResult) {
	client, err := c.getConnClient()
	if err != nil {
		logger.Default().Error(err, "getConnClient error")
		return err, nil
	}
	var reply packets.JsonResult
	err = client.Call("Rpc.StartExecutor", taskId, &reply)
	if err != nil {
		return err, nil
	}
	return nil, &reply
}

func (c *RpcClient) CallStopExecutor(taskId string) (error, *packets.JsonResult) {
	client, err := c.getConnClient()
	if err != nil {
		logger.Default().Error(err, "getConnClient error")
		return err, nil
	}
	var reply packets.JsonResult
	err = client.Call("Rpc.StopExecutor", taskId, &reply)
	if err != nil {
		return err, nil
	}
	return nil, &reply
}

func (c *RpcClient) CallRemoveExecutor(taskId string) (error, *packets.JsonResult) {
	client, err := c.getConnClient()
	if err != nil {
		logger.Default().Error(err, "getConnClient error")
		return err, nil
	}
	var reply packets.JsonResult
	err = client.Call("Rpc.RemoveExecutor", taskId, &reply)
	if err != nil {
		return err, nil
	}
	return nil, &reply
}

func (c *RpcClient) CallQueryExecutorConfig(taskId string) (error, *packets.JsonResult) {
	client, err := c.getConnClient()
	if err != nil {
		logger.Default().Error(err, "getConnClient error")
		return err, nil
	}
	var reply packets.JsonResult
	err = client.Call("Rpc.QueryExecutorConfig", taskId, &reply)
	if err != nil {
		return err, nil
	}
	return nil, &reply
}

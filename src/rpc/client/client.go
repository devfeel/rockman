package client

import (
	"github.com/devfeel/rockman/src/logger"
	"github.com/devfeel/rockman/src/packets"
	"github.com/devfeel/rockman/src/rpc/packet"
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
	var reply packet.JsonResult
	err = client.Call("Rpc.Echo", message, &reply)
	if err != nil {
		return err, ""
	}
	return nil, reply.Message.(string)
}

func (c *RpcClient) CallRegisterWorker(worker *packets.WorkerInfo) (error, map[string]interface{}) {
	client, err := c.getConnClient()
	if err != nil {
		logger.Default().Error(err, "getConnClient error")
		return err, nil
	}
	var reply packet.JsonResult
	err = client.Call("Rpc.RegisterWorker", worker, &reply)
	if err != nil {
		return err, nil
	}
	if reply.Message == nil {
		return nil, make(map[string]interface{})
	} else {
		return nil, reply.Message.(map[string]interface{})
	}
}

func (c *RpcClient) CallRegisterExecutor(conf interface{}) (error, interface{}) {
	client, err := c.getConnClient()
	if err != nil {
		logger.Default().Error(err, "getConnClient error")
		return err, ""
	}
	var reply packet.JsonResult
	err = client.Call("Rpc.RegisterExecutor", conf, &reply)
	if err != nil {
		return err, ""
	}
	return nil, reply
}

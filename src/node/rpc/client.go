package rpc

import (
	"github.com/devfeel/rockman/src/logger"
	"github.com/michain/dotcoin/server/packet"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type RpcClient struct {
	serviceUrl string
	client     *rpc.Client
}

func NewRpcClient(host, port string) *RpcClient {
	return &RpcClient{serviceUrl: host + ":" + port}
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

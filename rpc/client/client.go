package client

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/devfeel/rockman/core"
	"github.com/devfeel/rockman/logger"
	"github.com/devfeel/rockman/rpc/packet"
	"io/ioutil"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type RpcClient struct {
	serviceUrl string
	enableTls  bool
	certFile   string
	keyFile    string
	client     *rpc.Client
}

func NewRpcClient(serverUrl string, enableTls bool, certFile, keyFile string) *RpcClient {
	return &RpcClient{serviceUrl: serverUrl, enableTls: enableTls, certFile: certFile, keyFile: keyFile}
}

// getConnClient
func (c *RpcClient) getConnClient() (*rpc.Client, error) {
	if c.client != nil {
		return c.client, nil
	}
	if c.enableTls {
		tlsConfig, err := c.createTlsConfig()
		if err != nil {
			logger.Default().Error(err, "RpcClient createTlsConfig error")
			return nil, err
		}
		conn, err := tls.Dial("tcp", c.serviceUrl, tlsConfig)
		if err != nil {
			logger.Default().Error(err, "RpcClient connServer tls.dial error:")
			return nil, err
		} else {
			c.client = jsonrpc.NewClient(conn)
		}
	} else {
		conn, err := net.Dial("tcp", c.serviceUrl)
		if err != nil {
			logger.Default().Error(err, "RpcClient connServer net.dial error:")
			return nil, err
		} else {
			c.client = jsonrpc.NewClient(conn)
		}
	}

	return c.client, nil
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

func (c *RpcClient) CallQueryResource() (error, *packet.RpcReply) {
	client, err := c.getConnClient()
	if err != nil {
		logger.Default().Error(err, "getConnClient error")
		return err, nil
	}
	var reply packet.RpcReply
	err = client.Call("Rpc.QueryResource", "", &reply)
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

func (c *RpcClient) CallSubmitExecutor(submit *core.SubmitInfo) (error, *packet.RpcReply) {
	client, err := c.getConnClient()
	if err != nil {
		logger.Default().Error(err, "getConnClient error")
		return err, nil
	}
	var reply packet.RpcReply
	err = client.Call("Rpc.SubmitExecutor", submit, &reply)
	if err != nil {
		return err, nil
	}
	return nil, &reply
}

func (c *RpcClient) CallRegisterExecutor(conf *core.TaskConfig) (error, *packet.RpcReply) {
	client, err := c.getConnClient()
	if err != nil {
		logger.Default().Error(err, "getConnClient error")
		return err, nil
	}
	var reply packet.RpcReply
	err = client.Call("Rpc.RegisterExecutor", conf, &reply)
	if err != nil {
		return err, nil
	}
	return nil, &reply
}

func (c *RpcClient) CallStartExecutor(taskId string) (error, *packet.RpcReply) {
	client, err := c.getConnClient()
	if err != nil {
		logger.Default().Error(err, "getConnClient error")
		return err, nil
	}
	var reply packet.RpcReply
	err = client.Call("Rpc.StartExecutor", taskId, &reply)
	if err != nil {
		return err, nil
	}
	return nil, &reply
}

func (c *RpcClient) CallStopExecutor(taskId string) (error, *packet.RpcReply) {
	client, err := c.getConnClient()
	if err != nil {
		logger.Default().Error(err, "getConnClient error")
		return err, nil
	}
	var reply packet.RpcReply
	err = client.Call("Rpc.StopExecutor", taskId, &reply)
	if err != nil {
		return err, nil
	}
	return nil, &reply
}

func (c *RpcClient) CallRemoveExecutor(taskId string) (error, *packet.RpcReply) {
	client, err := c.getConnClient()
	if err != nil {
		logger.Default().Error(err, "getConnClient error")
		return err, nil
	}
	var reply packet.RpcReply
	err = client.Call("Rpc.RemoveExecutor", taskId, &reply)
	if err != nil {
		return err, nil
	}
	return nil, &reply
}

func (c *RpcClient) CallQueryExecutorConfig(taskId string) (error, *packet.RpcReply) {
	client, err := c.getConnClient()
	if err != nil {
		logger.Default().Error(err, "getConnClient error")
		return err, nil
	}
	var reply packet.RpcReply
	err = client.Call("Rpc.QueryExecutorConfig", taskId, &reply)
	if err != nil {
		return err, nil
	}
	return nil, &reply
}

func (c *RpcClient) createTlsConfig() (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(c.certFile, c.keyFile)
	if err != nil {
		return nil, err
	}
	certBytes, err := ioutil.ReadFile(c.certFile)
	if err != nil {
		return nil, err
	}
	clientCertPool := x509.NewCertPool()
	ok := clientCertPool.AppendCertsFromPEM(certBytes)
	if !ok {
		fmt.Println("AppendCertsFromPEM err")
		return nil, err
	}
	tlsConfig := &tls.Config{
		RootCAs:            clientCertPool,
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
	}
	return tlsConfig, nil
}

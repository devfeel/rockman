package rpc

import (
	"github.com/devfeel/rockman/src/config"
	"github.com/devfeel/rockman/src/logger"
	"github.com/devfeel/rockman/src/node"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type RpcServer struct {
	RpcHost     string
	RpcPort     string
	RpcProtocol string
	Node        *node.Node
}

func NewRpcServer(profile *config.Profile, node *node.Node) *RpcServer {
	s := new(RpcServer)
	s.Node = node
	s.RpcHost = profile.Rpc.RpcHost
	s.RpcPort = profile.Rpc.RpcPort
	s.RpcProtocol = profile.Rpc.RpcProtocol
	logger.Default().Debug("RpcServer Init Success!")
	return s
}

func (s *RpcServer) Listen() error {
	lis, err := net.Listen("tcp", s.RpcHost+":"+s.RpcPort)
	if err != nil {
		return err
	}
	defer lis.Close()

	srv := rpc.NewServer()
	if err := srv.RegisterName("Rpc", NewRpcHandler(s)); err != nil {
		logger.Default().Error(err, "lis.RegisterName error")
		return err
	}

	logger.Default().DebugF("RPCServer begin listen %s", lis.Addr())

	for {
		conn, err := lis.Accept()
		if err != nil {
			logger.Default().Error(err, "lis.Accept() error")
			continue
		}
		go srv.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}

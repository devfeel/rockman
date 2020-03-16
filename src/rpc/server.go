package rpc

import (
	"github.com/devfeel/rockman/src/config"
	"github.com/devfeel/rockman/src/logger"
	"github.com/devfeel/rockman/src/node"
	"github.com/devfeel/rockman/src/rpc/handler"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

const (
	DefaultHost    = "127.0.0.1"
	DefaultRpcPort = "2398" //2398 = 1983+0415 my birthday!
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
	logger.Default().Debug("RpcServer init success.")
	return s
}

func (s *RpcServer) Listen() error {
	lis, err := net.Listen("tcp", s.RpcHost+":"+s.RpcPort)
	if err != nil {
		return err
	}
	defer lis.Close()

	srv := rpc.NewServer()
	if err := srv.RegisterName("Rpc", handler.NewRpcHandler(s.Node)); err != nil {
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

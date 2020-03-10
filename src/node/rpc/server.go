package rpc

import (
	"github.com/devfeel/rockman/src/logger"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type RpcServer struct {
	RpcHost     string
	RpcPort     string
	RpcProtocol string
}

func NewRpcServer(host, port, protocol string) *RpcServer {
	s := new(RpcServer)
	s.RpcHost = host
	s.RpcPort = port
	s.RpcProtocol = protocol
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
	if err := srv.RegisterName("Rpc", &RpcHandler{}); err != nil {
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

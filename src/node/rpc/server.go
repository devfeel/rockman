package rpc

import (
	"github.com/devfeel/rockman/src/logger"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type RpcServer struct {
}

func NewRpcServer() *RpcServer {
	s := new(RpcServer)
	logger.Default().Debug("RpcServer Init Success!")
	return s
}

func (s *RpcServer) Listen(host, port string) error {
	lis, err := net.Listen("tcp", host+":"+port)
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

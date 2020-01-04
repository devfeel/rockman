package rpc

import (
	"github.com/devfeel/rockman/src/logger"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

var Logger = logger.GetLogger(logger.LoggerName_Default)

type RpcServer struct {
}

func NewRpcServer() *RpcServer {
	s := new(RpcServer)
	logger.Default().Debug("RpcServer Init Success!")
	return s
}

func (s *RpcServer) Listen(host, port string) {
	lis, err := net.Listen("tcp", host+":"+port)
	if err != nil {
		return
	}
	defer lis.Close()

	srv := rpc.NewServer()
	if err := srv.RegisterName("Rpc", &RpcHandler{}); err != nil {
		return
	}

	Logger.DebugF("RPCServer begin listen %s", lis.Addr())

	for {
		conn, err := lis.Accept()
		if err != nil {
			Logger.Error(err, "lis.Accept() error")
			continue
		}
		go srv.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}

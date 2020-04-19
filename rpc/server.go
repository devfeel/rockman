package rpc

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"github.com/devfeel/rockman/config"
	"github.com/devfeel/rockman/logger"
	"github.com/devfeel/rockman/node"
	"github.com/devfeel/rockman/rpc/handler"
	"io/ioutil"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

const (
	DefaultHost    = "127.0.0.1"
	DefaultRpcPort = "2398" //2398 = 1983+0415 my birthday!
)

type RpcServer struct {
	config      *config.Profile
	RpcHost     string
	RpcPort     string
	RpcProtocol string
	Node        *node.Node
}

func NewRpcServer(profile *config.Profile, node *node.Node) *RpcServer {
	s := new(RpcServer)
	s.config = profile
	s.Node = node
	s.RpcHost = profile.Rpc.RpcHost
	s.RpcPort = profile.Rpc.RpcPort
	s.RpcProtocol = profile.Rpc.RpcProtocol
	logger.Default().Debug("RpcServer init success.")
	return s
}

func (s *RpcServer) Listen() error {
	var listener net.Listener
	var err error

	if s.config.Rpc.EnableTls {
		tlsConfig, err := s.createTlsConfig()
		if err != nil {
			logger.Default().Error(err, "RPCServer createTlsConfig error")
			return err
		}
		listener, err = tls.Listen("tcp", s.RpcHost+":"+s.RpcPort, tlsConfig)
	} else {
		listener, err = net.Listen("tcp", s.RpcHost+":"+s.RpcPort)
	}

	if err != nil {
		return err
	}
	defer listener.Close()

	srv := rpc.NewServer()
	if err := srv.RegisterName("Rpc", handler.NewRpcHandler(s.Node)); err != nil {
		logger.Default().Error(err, "RPCServer lis.RegisterName error")
		return err
	}

	logger.Default().DebugF("RPCServer begin listen %s", listener.Addr())

	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.Default().Error(err, "lis.Accept() error")
			continue
		}
		go srv.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}

func (s *RpcServer) createTlsConfig() (*tls.Config, error) {
	serverCertFile := s.config.Rpc.ServerCertFile
	serverKeyFile := s.config.Rpc.ServerKeyFile
	clientCertFile := s.config.Rpc.ClientCertFile
	cert, err := tls.LoadX509KeyPair(serverCertFile, serverKeyFile)
	if err != nil {
		return nil, err
	}
	certBytes, err := ioutil.ReadFile(clientCertFile)
	if err != nil {
		return nil, err
	}
	clientCertPool := x509.NewCertPool()
	ok := clientCertPool.AppendCertsFromPEM(certBytes)
	if !ok {
		err = errors.New("AppendCertsFromPEM failed")
		return nil, err
	}
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    clientCertPool,
	}
	return tlsConfig, nil
}

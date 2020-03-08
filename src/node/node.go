package node

import (
	"errors"
	"github.com/devfeel/rockman/src/cluster"
	"github.com/devfeel/rockman/src/config"
	"github.com/devfeel/rockman/src/logger"
	"github.com/devfeel/rockman/src/node/rpc"
	"github.com/devfeel/rockman/src/runtime"
	"github.com/devfeel/rockman/src/runtime/executor"
	"github.com/devfeel/rockman/src/webui"
	"strconv"
)

const (
	defaultHost     = "127.0.0.1"
	defaultRpcPort  = 2398 //2398 = 1983+0415 my birthday!
	defaultHttpPort = 8080
)

type (
	Node struct {
		NodeId    string
		NodeName  string
		Status    int
		Config    *NodeConfig
		Runtime   *runtime.Runtime
		Cluster   *cluster.Cluster
		WebServer *webui.WebServer
		RpcServer *rpc.RpcServer
	}

	NodeConfig struct {
		RpcHost        string
		RpcPort        int
		RpcProtocol    string
		HttpHost       string
		HttpPort       int
		IsMaster       bool
		IsWorker       bool
		LogFilePath    string
		RegistryServer string
	}
)

func NewNode(profile *config.Profile) (*Node, error) {
	node := &Node{NodeId: profile.Node.NodeId}

	//init logger
	logger.Default().Debug("Node {" + node.NodeId + "} Start...")

	//init config
	err := node.initConfig(profile)
	if err != nil {
		return nil, errors.New("Node Init Config error: " + err.Error())
	}

	node.Cluster, err = cluster.NewCluster(profile.Registry.ServerUrl)
	if err != nil {
		return nil, errors.New("Node New Cluster error: " + err.Error())
	}
	node.RpcServer = rpc.NewRpcServer()

	if node.Config.IsMaster {
		node.WebServer = webui.NewWebServer(profile.Logger.LogPath)
	}

	if node.Config.IsWorker {
		// create runtime
		node.Runtime = runtime.NewRuntime()

		// load tasks
		// TODO load tasks from mysql
		registerDemoExecutors(node.Runtime)
	}

	return node, err
}

func (n *Node) Start() error {
	if n.Config.IsWorker {
		go n.Runtime.Start()
	}
	if n.Config.IsMaster {
		go n.WebServer.ListenAndServe(n.Config.HttpHost + ":" + strconv.Itoa(n.Config.HttpPort))
	}

	// start rpcserver listen
	go n.RpcServer.Listen(n.Config.RpcHost, strconv.Itoa(n.Config.RpcPort))

	//n.Cluster.Registry.Register(n.Config.RegistryServer)
	return nil
}

func (n *Node) initConfig(conf *config.Profile) error {
	n.Config = new(NodeConfig)
	n.Config.HttpHost = defaultHost
	n.Config.HttpPort = defaultHttpPort
	n.Config.RpcHost = conf.Node.RpcHost
	n.Config.RpcPort = defaultRpcPort
	n.Config.RpcProtocol = conf.Node.RpcProtocol

	if conf.Node.RpcPort > 0 {
		n.Config.RpcPort = conf.Node.RpcPort
	}

	n.Config.HttpHost = conf.Node.HttpHost
	if conf.Node.HttpPort > 0 {
		n.Config.HttpPort = conf.Node.HttpPort
	}

	n.Config.IsMaster = conf.Node.IsMaster
	n.Config.IsWorker = conf.Node.IsWorker

	logger.Default().Debug("Config Init Success!")
	return nil
}

func registerDemoExecutors(r *runtime.Runtime) {
	logger.Default().Debug("Register Demo Executors Begin")
	goExec := executor.NewDebugGoExecutor(("go"))
	err := r.RegisterExecutor(goExec)
	if err != nil {
		logger.Default().Error(err, "service.CreateCronTask {go.exec} error!")
	}

	httpExec := executor.NewDebugHttpExecutor("http")
	err = r.RegisterExecutor(httpExec)
	if err != nil {
		logger.Default().Error(err, "service.CreateCronTask {http.exec} error!")
	}

	shellExec := executor.NewDebugShellExecutor("shell")
	err = r.RegisterExecutor(shellExec)
	if err != nil {
		logger.Default().Error(err, "service.CreateCronTask {shell.exec} error!")
	}
	logger.Default().Debug("Register Demo Executors Success!")
}

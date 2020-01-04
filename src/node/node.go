package node

import (
	"errors"
	"github.com/devfeel/rockman/src/config"
	"github.com/devfeel/rockman/src/logger"
	"github.com/devfeel/rockman/src/runtime"
	"github.com/devfeel/rockman/src/runtime/executor"
)

const (
	RunMode_Single  = "single"
	RunMode_Cluster = "cluster"

	defaultRpcPort  = 2020
	defaultHttpPort = 8080
)

var NodeLogger = logger.GetLogger(logger.LoggerName_Node)

type (
	Node struct {
		NodeId   string
		NodeName string
		Logger   logger.Logger
		Config   *NodeConfig
		Status   int //NodeStatus
		Runtime  *runtime.Runtime
	}

	NodeConfig struct {
		RpcPort        int
		RpcProtocol    string
		HttpPort       int
		RunMode        string
		IsMaster       bool
		IsWorker       bool
		LogFilePath    string
		RegistryServer string
	}
)

func NewNode(profile *config.Profile) (*Node, error) {
	node := &Node{NodeId: profile.Node.NodeId}

	//init logger
	node.Logger = logger.GetLogger(logger.LoggerName_Node)
	node.Logger.Debug("Node {" + node.NodeId + "} Start...")

	//init config
	err := node.initConfig(profile)
	if err != nil {
		return nil, errors.New("Node Init Config error: " + err.Error())
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
		n.Runtime.Start()
	}
	return nil
}

func (n *Node) initConfig(conf *config.Profile) error {
	n.Config = new(NodeConfig)
	n.Config.HttpPort = defaultHttpPort
	n.Config.RpcPort = defaultRpcPort
	n.Config.RpcProtocol = conf.Node.RpcProtocol
	if conf.Node.RpcPort > 0 {
		n.Config.RpcPort = conf.Node.RpcPort
	}
	if conf.Node.HttpPort > 0 {
		n.Config.HttpPort = conf.Node.HttpPort
	}
	n.Config.RunMode = conf.Node.RunMode
	n.Config.IsMaster = conf.Node.IsMaster
	n.Config.IsWorker = conf.Node.IsWorker

	return nil
}

func registerDemoExecutors(r *runtime.Runtime) {
	goExec := executor.NewDebugGoExecutor(("go"))
	err := r.RegisterExecutor(goExec)
	if err != nil {
		NodeLogger.Error(err, "service.CreateCronTask {go.exec} error!")
	}

	httpExec := executor.NewDebugHttpExecutor("http")
	err = r.RegisterExecutor(httpExec)
	if err != nil {
		NodeLogger.Error(err, "service.CreateCronTask {http.exec} error!")
	}

	shellExec := executor.NewDebugShellExecutor("shell")
	err = r.RegisterExecutor(shellExec)
	if err != nil {
		NodeLogger.Error(err, "service.CreateCronTask {shell.exec} error!")
	}
}

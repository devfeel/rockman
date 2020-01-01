package node

import (
	"errors"
	"github.com/devfeel/rockman/src/config"
	"github.com/devfeel/rockman/src/logger"
	"github.com/devfeel/rockman/src/runtime"
)

const (
	RunMode_Single  = "single"
	RunMode_Cluster = "cluster"

	defaultRpcPort  = 2020
	defaultHttpPort = 8080
)

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
	node.Logger = logger.NewFileLogger(profile.Logger.LogPath)
	node.Logger.Debug("Node {" + node.NodeId + "} Start...")

	//init config
	err := node.initConfig(profile)
	if err != nil {
		return nil, errors.New("Node Init Config error: " + err.Error())
	}

	node.Runtime = runtime.NewRuntime()

	return node, err

}

func (n *Node) Start() error {
	n.Runtime.Start()
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

package node

import (
	"errors"
	"github.com/devfeel/rockman/cluster"
	"github.com/devfeel/rockman/config"
	"github.com/devfeel/rockman/logger"
	"github.com/devfeel/rockman/packets"
	"github.com/devfeel/rockman/runtime"
	"github.com/devfeel/rockman/runtime/executor"
	"time"
)

type (
	Node struct {
		NodeId    string
		NodeName  string
		ClusterId string
		IsLeader  bool
		Status    int
		Config    *NodeConfig
		Cluster   *cluster.Cluster
		Runtime   *runtime.Runtime
	}

	NodeConfig struct {
		IsMaster       bool
		IsWorker       bool
		RpcServer      string //current node's rpc server info
		LogFilePath    string
		RegistryServer string
		Profile        *config.Profile
	}
)

func NewNode(profile *config.Profile) (*Node, error) {
	node := &Node{NodeId: profile.Node.NodeId, NodeName: profile.Node.NodeName}

	logger.Default().Debug("Node {" + node.NodeId + "} start...")

	//init config
	err := node.initConfig(profile)
	if err != nil {
		return nil, errors.New("Node Config init error: " + err.Error())
	}

	//init Registry
	cluster, err := cluster.NewCluster(profile.Cluster.ClusterId, profile.Cluster.RegistryServer, getLeaderKey(profile.Cluster.ClusterId))
	if err != nil {
		return nil, err
	}
	node.Cluster = cluster

	if node.Config.IsMaster {
		//register master role
		go node.ElectionLeader()
	}

	if node.Config.IsWorker {
		// create runtime
		node.Runtime = runtime.NewRuntime()

		//register worker
		go func() {
			worker := &packets.WorkerInfo{NodeID: node.NodeId, Host: profile.Rpc.RpcHost, Port: profile.Rpc.RpcPort}
		RegisterWorker:
			for {
				err := node.Cluster.RegisterWorker(worker)
				if err != nil {
					logger.Node().DebugS("Node RegisterWorker error, will retry after 10 seconds")
					logger.Node().Error(err, "Node RegisterWorker error.")
					time.Sleep(time.Second * 10)
					continue RegisterWorker
				} else {
					logger.Node().DebugS("Node RegisterWorker success:", worker)
					break
				}
			}

		}()
	}

	logger.Node().Debug("Node init success.")
	return node, err
}

func (n *Node) Start() error {
	if n.Config.IsWorker {
		// load self tasks
		// TODO load self tasks

		go n.Runtime.Start()
	}
	//n.Cluster.Registry.Register(n.Config.RegistryServer)
	return nil
}

func (n *Node) ElectionLeader() error {
	isLeader, err := n.Cluster.ElectionLeader(n.Config.RpcServer, "")
	if err == nil {
		n.IsLeader = isLeader
	} else {
		logger.Default().Error(err, "Node {"+n.NodeId+"} election leader role with key {"+n.Cluster.LeaderKey+"} error:"+err.Error())
	}

	if n.IsLeader {
		//TODO do something when change to leader
		logger.Default().Debug("Node {" + n.NodeId + "} election leader role success with key {" + n.Cluster.LeaderKey + "}")
	}
	return nil
}

// initConfig init node config from config profile
func (n *Node) initConfig(conf *config.Profile) error {
	n.Config = new(NodeConfig)
	n.Config.RpcServer = conf.Rpc.RpcHost + ":" + conf.Rpc.RpcPort
	n.Config.RegistryServer = conf.Cluster.RegistryServer
	n.Config.IsMaster = conf.Node.IsMaster
	n.Config.IsWorker = conf.Node.IsWorker
	n.Config.Profile = conf

	logger.Default().Debug("Node Config init success.")
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

func loadHttpTaskConfigs() []*executor.HttpTaskConfig {
	//TODO load http task config from mysql
	return []*executor.HttpTaskConfig{}
}

func loadShellTaskConfigs() []*executor.ShellTaskConfig {
	//TODO load shell task config from mysql
	return []*executor.ShellTaskConfig{}
}

func loadGoTaskConfigs() []*executor.GoTaskConfig {
	//TODO load go task config from mysql
	return []*executor.GoTaskConfig{}
}

func getLeaderKey(clusterId string) string {
	return "devfeel/rockman:" + clusterId + ":leader:locker"
}

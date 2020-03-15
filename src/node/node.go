package node

import (
	"errors"
	"fmt"
	"github.com/devfeel/rockman/src/config"
	"github.com/devfeel/rockman/src/logger"
	"github.com/devfeel/rockman/src/rpc"
	"github.com/devfeel/rockman/src/runtime"
	"github.com/devfeel/rockman/src/runtime/executor"
	"sync"
)

type (
	Node struct {
		NodeId       string
		NodeName     string
		ClusterId    string
		IsLeader     bool
		Status       int
		Workers      map[string]*WorkerInfo
		workerLocker *sync.RWMutex
		Config       *NodeConfig
		Registry     *Registry
		Runtime      *runtime.Runtime
	}

	NodeConfig struct {
		IsMaster       bool
		IsWorker       bool
		RpcServer      string //current node's rpc server info
		LogFilePath    string
		RegistryServer string
		Profile        *config.Profile
	}

	WorkerInfo struct {
		NodeID string
		Host   string
		Port   string
	}
)

func NewNode(profile *config.Profile) (*Node, error) {
	node := &Node{NodeId: profile.Node.NodeId, NodeName: profile.Node.NodeName, ClusterId: profile.Node.ClusterId}

	logger.Default().Debug("Node {" + node.NodeId + "} start...")

	//init config
	err := node.initConfig(profile)
	if err != nil {
		return nil, errors.New("Node Config init error: " + err.Error())
	}

	//init Registry
	register, err := initRegistry(profile.Registry.ServerUrl, getLeaderKey(profile.Node.ClusterId))
	if err != nil {
		return nil, err
	}
	node.Registry = register

	//init workers
	node.Workers = make(map[string]*WorkerInfo)
	node.workerLocker = new(sync.RWMutex)

	if node.Config.IsMaster {
		//register master role
		go node.ElectionLeader()
	}

	if node.Config.IsWorker {
		// create runtime
		node.Runtime = runtime.NewRuntime()

		// get leader info
		leaderServer, err := node.Registry.GetLeaderInfo()
		if err != nil {
			logger.Node().DebugS("Node GetLeaderInfo error:", err.Error())
			logger.Node().Error(err, "Node GetLeaderInfo error.")
		} else {
			logger.Node().DebugS("Node GetLeaderInfo success:", leaderServer)
			//register worker
			rpcClient := rpc.NewRpcClient(leaderServer)
			worker := WorkerInfo{NodeID: node.NodeId, Host: profile.Rpc.RpcHost, Port: profile.Rpc.RpcPort}
			err, _ := rpcClient.CallRegisterWorker(worker)
			if err != nil {
				logger.Node().DebugS("Node RegisterWorker error:", err.Error())
				logger.Node().Error(err, "Node RegisterWorker error.")
			} else {
				logger.Node().DebugS("Node RegisterWorker success:", worker)
			}
		}
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
	isLeader, err := n.Registry.electionLeader(n.Config.RpcServer, "")
	if err == nil {
		n.IsLeader = isLeader
	} else {
		logger.Default().Error(err, "Node {"+n.NodeId+"} election leader role with key {"+n.Registry.LeaderKey+"} error:"+err.Error())
	}

	if n.IsLeader {
		//TODO do something when change to leader
		logger.Default().Debug("Node {" + n.NodeId + "} election leader role success with key {" + n.Registry.LeaderKey + "}")
	}
	return nil
}

// AddWorker add worker into node workers
func (n *Node) AddWorker(worker *WorkerInfo) error {
	key := worker.Host + "," + worker.Port
	n.workerLocker.Lock()
	defer n.workerLocker.Unlock()
	rawWorker, isExists := n.Workers[key]
	if isExists {
		logger.Default().Debug("replace master node:" + fmt.Sprint(rawWorker, worker))
	} else {
		logger.Default().Debug("add master node:" + fmt.Sprint(worker))
	}
	n.Workers[key] = worker
	return nil
}

// initConfig init node config from config profile
func (n *Node) initConfig(conf *config.Profile) error {
	n.Config = new(NodeConfig)
	n.Config.RpcServer = conf.Rpc.RpcHost + ":" + conf.Rpc.RpcPort
	n.Config.RegistryServer = conf.Registry.ServerUrl
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

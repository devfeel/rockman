package node

import (
	"errors"
	"fmt"
	"github.com/devfeel/rockman/src/config"
	"github.com/devfeel/rockman/src/logger"
	"github.com/devfeel/rockman/src/runtime"
	"github.com/devfeel/rockman/src/runtime/executor"
	"github.com/devfeel/rockman/src/util/consul"
	"github.com/hashicorp/consul/api"
	"sync"
)

const (
	registryLockerKey = "master:locker"
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
		ServerUrl      string
		LogFilePath    string
		RegistryServer string
		Profile        *config.Profile
	}

	Registry struct {
		ServerUrl string
		RegServer *consul.ConsulClient
	}

	WorkerInfo struct {
		NodeID string
		Host   string
		Port   string
	}
)

func NewNode(profile *config.Profile) (*Node, error) {
	node := &Node{NodeId: profile.Node.NodeId}

	logger.Default().Debug("Node {" + node.NodeId + "} start...")

	//init config
	err := node.initConfig(profile)
	if err != nil {
		return nil, errors.New("Node Config init error: " + err.Error())
	}

	//init Registry
	err = node.initRegistry()
	if err != nil {
		return nil, err
	}

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
	isLeader, err := n.electionLeader(n.Config.ServerUrl, "")
	if err == nil {
		n.IsLeader = isLeader
	} else {
		logger.Default().Error(err, "Node {"+n.NodeId+"} ElectionLeader error:"+err.Error())
	}

	if n.IsLeader {
		//TODO do something when change to leader
		logger.Default().Debug("Node {" + n.NodeId + "} election leader role success")
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

// electionLeader election leader role to registry server
func (n *Node) electionLeader(serverUrl string, checkUrl string) (bool, error) {
	opts := &api.LockOptions{
		Key:         getRegistryLockerKey(n.ClusterId),
		Value:       []byte(serverUrl),
		SessionTTL:  "10s",
		SessionName: serverUrl,
	}
	locker, err := n.Registry.RegServer.CreateLockerOpts(opts)
	if err != nil {
		return false, err
	}

	_, err = locker.Locker.Lock(nil)
	if err != nil {
		return false, err
	}
	return true, nil
}

// initRegistry init Registry and reg server
func (n *Node) initRegistry() error {
	n.Registry = new(Registry)
	n.Registry.ServerUrl = n.Config.RegistryServer
	regServer, err := consul.NewConsulClient(n.Registry.ServerUrl)
	if err != nil {
		logger.Node().Debug(fmt.Sprint("Registry init error", err.Error()))
		logger.Node().Error(err, "Registry init error")
		return err
	}
	n.Registry.RegServer = regServer
	logger.Node().Debug("Registry init success.")
	return nil
}

// initConfig init node config from config profile
func (n *Node) initConfig(conf *config.Profile) error {
	n.Config = new(NodeConfig)
	n.Config.ServerUrl = conf.Rpc.RpcHost + ":" + conf.Rpc.RpcPort
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

func getRegistryLockerKey(clusterId string) string {
	return "devfeel/rockman:" + clusterId + ":" + registryLockerKey
}

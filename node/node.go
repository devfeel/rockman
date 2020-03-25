package node

import (
	"errors"
	"fmt"
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
		NodeId           string
		NodeName         string
		IsLeader         bool
		Status           int
		Config           *NodeConfig
		submitList       map[string]*packets.SubmitInfo
		submitQueue      chan *packets.SubmitInfo
		submitRetryQueue chan *packets.SubmitInfo
		Cluster          *cluster.Cluster
		Runtime          *runtime.Runtime
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

var (
	ErrorCanNotSubmitToNotLeaderNode = errors.New("can not submit to not leader node")
)

func NewNode(profile *config.Profile) (*Node, error) {
	logger.Node().Debug("Node {" + profile.Node.NodeId + "} start...")

	node := &Node{
		NodeId:           profile.Node.NodeId,
		NodeName:         profile.Node.NodeName,
		submitList:       make(map[string]*packets.SubmitInfo),
		submitQueue:      make(chan *packets.SubmitInfo),
		submitRetryQueue: make(chan *packets.SubmitInfo),
	}

	nodeInfo := &packets.NodeInfo{NodeID: node.NodeId, Host: profile.Rpc.RpcHost, Port: profile.Rpc.RpcPort}
	nodeKey := nodeInfo.GetNodeKey(profile.Cluster.ClusterId)

	//init config
	err := node.initConfig(profile)
	if err != nil {
		logger.Node().Debug("Node init config error: " + err.Error())
		return nil, err
	}

	//init Registry
	cluster, err := cluster.NewCluster(profile.Cluster.ClusterId, profile.Cluster.RegistryServer, getLeaderKey(profile.Cluster.ClusterId))
	if err != nil {
		return nil, err
	}
	node.Cluster = cluster

	if node.Config.IsMaster {
		// election leader role
		go node.ElectionLeader()
	}

	if node.Config.IsWorker {
		// create runtime
		node.Runtime = runtime.NewRuntime()
	}

	// create session with node info
	go func() {
		err := node.Cluster.CreateSession(nodeKey, nodeInfo)
		if err != nil {
			logger.Node().Debug("Node{" + node.NodeId + "} create session to registry error: " + err.Error())
		} else {
			logger.Node().Debug("Node{" + node.NodeId + "} create session to registry success with key {" + nodeKey + "}")
		}
	}()

	// register node to cluster
	go func() {
		node.registerNode(nodeInfo)
	}()

	logger.Node().Debug("Node init success.")
	return node, err
}

func (n *Node) Start() error {
	if n.Config.IsWorker {
		// load self tasks
		// TODO load self tasks
		go n.Runtime.Start()
	}

	if n.IsLeader {
		go n.distributeSubmit()
	}

	return nil
}

func (n *Node) ElectionLeader() error {
	isLeader, err := n.Cluster.ElectionLeader(n.Config.RpcServer, "")
	if err == nil {
		n.IsLeader = isLeader
	} else {
		logger.Node().Error(err, "Node {"+n.NodeId+"} election leader role with key {"+n.Cluster.LeaderKey+"} error:"+err.Error())
	}

	if n.IsLeader {
		//TODO do something when change to leader
		logger.Node().Debug("Node {" + n.NodeId + "} election leader role success with key {" + n.Cluster.LeaderKey + "}")
	}
	return nil
}

func (n *Node) SubmitExecutor(submit *packets.SubmitInfo) error {
	if !n.IsLeader {
		return ErrorCanNotSubmitToNotLeaderNode
	}
	n.submitQueue <- submit
	logger.Node().Debug("SubmitExecutor[" + fmt.Sprint(submit.ExecutorConfig) + "] into queue success")
	//TODO log submit to db log
	return nil
}

// registerNode register node to cluster
func (n *Node) registerNode(nodeInfo *packets.NodeInfo) {
RegisterWorker:
	for {
		err := n.Cluster.RegisterNode(nodeInfo)
		if err != nil {
			logger.Node().DebugS("RegisterNode error, will retry after 10 seconds")
			logger.Node().Error(err, "RegisterNode error.")
			time.Sleep(time.Second * 10)
			continue RegisterWorker
		} else {
			logger.Node().DebugS("RegisterNode success:", nodeInfo)
			break
		}
	}
}

// distributeSubmit distribute submit from queue, send to worker node
func (n *Node) distributeSubmit() {
	for {
		defer func() {
			if err := recover(); err != nil {
				errInfo := errors.New(fmt.Sprintln(err))
				logger.Node().Error(errInfo, "distributeSubmit error")
			}
		}()
		for {
			submit := <-n.submitQueue
			worker := submit.Worker
			var err error

			// get low balance worker
			if worker == nil {
				worker, err = n.Cluster.GetLowBalanceWorker()
				if err != nil {
					logger.Node().Error(err, "GetLowBalanceWorker error")
					//TODO log submit result to db log
					return
				}
			}

			//submit executor to the specified worker node
			rpcClient := n.Cluster.GetRpcClient(submit.Worker.Host, submit.Worker.Port)
			err, reply := rpcClient.CallRegisterExecutor(submit.ExecutorConfig)
			if err != nil {
				n.submitRetryQueue <- submit
				logger.Node().DebugS("distributeSubmit into retry queue, error:", err.Error())
				//TODO log submit result to db log
			} else {
				if reply.RetCode != reply.CorrectCode() {
					n.submitRetryQueue <- submit
					logger.Node().DebugS("distributeSubmit into retry queue, failed:", reply.RetCode)
					//TODO log submit result to db log
				} else {
					n.Cluster.Scheduler.AddJobInfo(worker.EndPoint(), 1)
					//TODO log submit result to db log
				}
			}
		}
	}
}

// initConfig init node config from config profile
func (n *Node) initConfig(conf *config.Profile) error {
	n.Config = new(NodeConfig)
	n.Config.RpcServer = conf.Rpc.RpcHost + ":" + conf.Rpc.RpcPort
	n.Config.RegistryServer = conf.Cluster.RegistryServer
	n.Config.IsMaster = conf.Node.IsMaster
	n.Config.IsWorker = conf.Node.IsWorker
	n.Config.Profile = conf

	logger.Node().Debug("Node Config init success.")
	return nil
}

func registerDemoExecutors(r *runtime.Runtime) {
	logger.Node().Debug("Register Demo Executors Begin")
	goExec := executor.NewDebugGoExecutor(("go"))
	err := r.RegisterExecutor(goExec)
	if err != nil {
		logger.Node().Error(err, "service.CreateCronTask {go.exec} error!")
	}

	httpExec := executor.NewDebugHttpExecutor("http")
	err = r.RegisterExecutor(httpExec)
	if err != nil {
		logger.Node().Error(err, "service.CreateCronTask {http.exec} error!")
	}

	shellExec := executor.NewDebugShellExecutor("shell")
	err = r.RegisterExecutor(shellExec)
	if err != nil {
		logger.Node().Error(err, "service.CreateCronTask {shell.exec} error!")
	}
	logger.Node().Debug("Register Demo Executors Success!")
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

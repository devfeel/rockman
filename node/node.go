package node

import (
	"errors"
	"github.com/devfeel/rockman/cluster"
	"github.com/devfeel/rockman/config"
	"github.com/devfeel/rockman/core"
	"github.com/devfeel/rockman/logger"
	"github.com/devfeel/rockman/protected/service"
	"github.com/devfeel/rockman/registry"
	"github.com/devfeel/rockman/rpc/client"
	"github.com/devfeel/rockman/runtime"
	"strconv"
	"time"
)

const defaultLockerTTL = "10s"

type (
	Node struct {
		NodeId                  string
		NodeName                string
		isLeader                bool
		Status                  int
		config                  *config.Profile
		nodeInfo                *core.NodeInfo
		Cluster                 *cluster.Cluster
		Registry                *registry.Registry
		Runtime                 *runtime.Runtime
		shutdownChan            chan string
		isSTW                   bool //stop the world flag
		logLogic                *service.LogService
		isRunCycleLoadExecutors bool
		executorFlagLastIndex   uint64
	}
)

var (
	ErrorCanNotSubmitToNotLeaderNode = errors.New("can not submit to not leader node")
	ErrorStopTheWorld                = errors.New("node is stop the world")
)

func NewNode(profile *config.Profile, shutdown chan string) (*Node, error) {
	logger.Node().Debug("Node {" + profile.Node.NodeId + "} begin init...")

	node := &Node{
		NodeId:       profile.Node.NodeId,
		NodeName:     profile.Node.NodeName,
		config:       profile,
		shutdownChan: shutdown,
		logLogic:     service.NewLogService(),
	}

	registry, err := registry.NewRegistry(profile.Cluster.RegistryServer)
	if err != nil {
		return nil, err
	}
	registry.OnServerOnline = node.onRegistryOnline
	registry.OnServerOffline = node.onRegistryOffline
	node.Registry = registry

	//init cluster
	cluster := cluster.NewCluster(profile, registry)
	cluster.OnLeaderChange = node.onLeaderChange
	cluster.OnLeaderChangeFailed = node.onLeaderChangeFailed

	node.Cluster = cluster
	if node.config.Node.IsWorker {
		node.Runtime = runtime.NewRuntime(node.NodeInfo(), profile)
	}

	logger.Node().Debug("Node init success.")
	return node, err
}

func (n *Node) Start() error {
	logger.Default().Debug("Node start...")
	// create session with node info
	err := n.createSession(n.NodeInfo().GetNodeKey(n.Cluster.ClusterId))
	if err != nil {
		return err
	}

	if n.config.Node.IsMaster {
		n.electionLeader()
	}

	// register node to cluster
	err = n.registerNode()
	if err != nil {
		return err
	}

	err = n.Registry.Start()
	if err != nil {
		return err
	}

	err = n.Cluster.Start()
	if err != nil {
		return err
	}

	if n.config.Node.IsWorker {
		err = n.Runtime.Start()
		if err != nil {
			return err
		}
	}

	return nil
}

func (n *Node) Config() *config.Profile {
	return n.config
}

func (n *Node) IsMaster() bool {
	return n.config.Node.IsMaster
}

func (n *Node) IsWorker() bool {
	return n.config.Node.IsWorker
}

func (n *Node) IsLeader() bool {
	return n.isLeader
}

func (n *Node) ClusterId() string {
	return n.Cluster.ClusterId
}

func (n *Node) Shutdown() {
	logTitle := "Node Shutdown "
	logger.Node().Debug(logTitle + "doing.")
	n.stopTheWorld()
	n.shutdownChan <- "ok"
}

func (n *Node) NodeInfo() *core.NodeInfo {
	if n.nodeInfo != nil {
		return n.nodeInfo
	}
	n.nodeInfo = &core.NodeInfo{
		NodeID:    n.NodeId,
		Cluster:   n.config.Cluster.ClusterId,
		OuterHost: n.config.Rpc.OuterHost,
		OuterPort: n.config.Rpc.OuterPort,
		Host:      n.config.Rpc.RpcHost,
		Port:      n.config.Rpc.RpcPort,
		IsMaster:  n.config.Node.IsMaster,
		IsWorker:  n.config.Node.IsWorker,
		IsOnline:  true,
	}
	if n.Runtime != nil {
		n.nodeInfo.Executors = n.Runtime.GetTaskIDs()
	}
	return n.nodeInfo
}

func (n *Node) stopTheWorld() {
	lt := "Node stopTheWorld "
	logger.Node().Debug(lt + "begin.")
	logger.Node().Debug(lt + "set SWT flag true.")
	n.isSTW = true
	n.Cluster.Stop()
	n.Runtime.Stop()
	logger.Node().Debug(lt + "success.")
}

func (n *Node) startTheWorld() {
	lt := "Node startTheWorld "
	logger.Node().Debug(lt + "begin.")

	logger.Node().Debug(lt + "set SWT flag false.")
	n.isSTW = false

	if n.config.Node.IsMaster {
		cluster := cluster.NewCluster(n.config, n.Registry)
		cluster.OnLeaderChange = n.onLeaderChange
		cluster.OnLeaderChangeFailed = n.onLeaderChangeFailed
		n.Cluster = cluster
	}

	if n.config.Node.IsWorker {
		n.Runtime = runtime.NewRuntime(n.NodeInfo(), n.config)
	}

	err := n.Start()
	if err != nil {
		logger.Node().Debug(lt + "failed, error: " + err.Error())
		n.Shutdown()
	} else {
		logger.Node().Debug(lt + "success.")
	}
}

// registerNode register node to cluster
func (n *Node) registerNode() error {
	logTitle := "Node registerNode "
	var leaderServer string
	var err error
	var retryCount int
	nodeInfo := n.NodeInfo()
	logger.Node().Debug(logTitle + "begin.")
RegisterNode:
	for {
		if n.isSTW {
			return ErrorStopTheWorld
		}
		if retryCount > n.config.Global.RetryLimit {
			err = errors.New("retry more than 5 times and stop it")
			logger.Node().DebugS(logTitle + "error: " + err.Error())
			return err
		}
		retryCount += 1
		// get leader info
		leaderServer, err = n.Cluster.GetLeaderInfo()
		if err != nil {
			logger.Node().Debug(logTitle + "GetLeaderInfo error:" + err.Error() + ", will retry 10 seconds after.")
			time.Sleep(time.Second * 10)
			continue RegisterNode
		} else {
			logger.Node().Debug(logTitle + "GetLeaderInfo success [" + leaderServer + "]")
			rpcClient := client.NewRpcClient(leaderServer, n.config.Rpc.EnableTls, n.config.Rpc.ClientCertFile, n.config.Rpc.ClientKeyFile)
			err, reply := rpcClient.CallRegisterNode(nodeInfo)
			if err != nil {
				logger.Node().Debug(logTitle + "CallRegisterNode error:" + err.Error() + ", will retry 10 seconds after.")
				time.Sleep(time.Second * 10)
				continue RegisterNode
			}
			if !reply.IsSuccess() {
				logger.Node().Debug(logTitle + "CallRegisterNode failed:" + strconv.Itoa(reply.RetCode) + ", will retry 10 seconds after.")
				time.Sleep(time.Second * 10)
				continue RegisterNode
			} else {
				retryCount = 0
				// if node is leader and register to self, mean cluster is init, remove old init flag
				if leaderServer == n.NodeInfo().EndPoint() {
					if err != nil {
						logger.Node().Warn(logTitle + "delete executor-init flag error:" + err.Error())
					}
				}
				logger.Node().DebugS(logTitle + "success.")
			}
			break
		}
	}
	return nil
}

// createSession create session to registry server
func (n *Node) createSession(nodeKey string) error {
	lt := "Node create session "
	logger.Node().Debug(lt + "begin.")

	locker, err := n.Registry.CreateLocker(nodeKey, n.NodeInfo().Json(), defaultLockerTTL)
	if err != nil {
		logger.Node().Debug(lt + "error: " + err.Error())
	}
	_, err = locker.Lock()
	if err != nil {
		logger.Node().Debug(lt + "error: " + err.Error())
		return err
	}
	logger.Node().Debug(lt + "success with key {" + nodeKey + "}")
	return nil
}

// onLeaderChange do something when leader is changed
func (n *Node) onLeaderChange() {
	err := n.registerNode()
	if err != nil {
		logger.Node().DebugS("Node.onLeaderChange registerNode error:", err.Error())
	} else {
		logger.Node().Debug("Node.onLeaderChange registerNode success")
	}
	if n.IsLeader() {
		if n.Cluster.LeaderServer != n.NodeInfo().EndPoint() {
			n.becomeLeaderRole()
		}
	}

	if n.IsMaster() && !n.IsLeader() {
		if n.Cluster.LeaderServer == n.NodeInfo().EndPoint() {
			n.removeLeaderRole()
		}
	}
}

// onLeaderChangeFailed
func (n *Node) onLeaderChangeFailed() {
	logger.Node().DebugS("Node.onLeaderChangeFailed, now will shutdown node.")
	n.Shutdown()
}

// onWorkerNodeOffline
func (n *Node) onWorkerNodeOffline(nodeInfo *core.NodeInfo) {
	logTitle := "Node.onWorkerNodeOffline[" + nodeInfo.NodeID + "] "
	if !n.isLeader {
		logger.Node().Warn(logTitle + "is be called, but it's not leader")
		return
	}
	var needReSubmits []*core.ExecutorInfo
	for _, v := range n.Cluster.ExecutorInfos {
		if v.Worker.NodeID == nodeInfo.NodeID {
			needReSubmits = append(needReSubmits, v)
		}
	}
	go func() {
		for _, exec := range needReSubmits {
			n.SubmitExecutor(exec)
		}
	}()
}

func (n *Node) onRegistryOnline() {
	logger.Node().DebugS("Node.onRegistryOnline registry online, now start the world.")
	if n.isSTW {
		n.startTheWorld()
	}
}

func (n *Node) onRegistryOffline() {
	logger.Node().DebugS("Node.onRegistryOffline registry offline, now stop the world.")
	n.stopTheWorld()
}

func (n *Node) refreshNodeInfo() *core.NodeInfo {
	n.nodeInfo = &core.NodeInfo{
		NodeID:    n.NodeId,
		Cluster:   n.config.Cluster.ClusterId,
		OuterHost: n.config.Rpc.OuterHost,
		OuterPort: n.config.Rpc.OuterPort,
		Host:      n.config.Rpc.RpcHost,
		Port:      n.config.Rpc.RpcPort,
		IsMaster:  n.config.Node.IsMaster,
		IsWorker:  n.config.Node.IsWorker,
		IsOnline:  true,
	}
	if n.IsWorker() {
		n.nodeInfo.Executors = n.Runtime.GetTaskIDs()
	}
	return n.nodeInfo
}

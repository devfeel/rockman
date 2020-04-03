package node

import (
	"errors"
	"github.com/devfeel/rockman/cluster"
	"github.com/devfeel/rockman/config"
	"github.com/devfeel/rockman/logger"
	"github.com/devfeel/rockman/packets"
	"github.com/devfeel/rockman/rpc/client"
	"github.com/devfeel/rockman/runtime"
	"strconv"
	"sync"
	"time"
)

type (
	Node struct {
		NodeId             string
		NodeName           string
		IsLeader           bool
		Status             int
		Config             *NodeConfig
		profile            *config.Profile
		onlineSubmits      map[string]*packets.SubmitInfo
		onlineSubmitLocker *sync.RWMutex
		Cluster            *cluster.Cluster
		Runtime            *runtime.Runtime
		shutdownChan       chan string
	}

	NodeConfig struct {
		IsMaster       bool
		IsWorker       bool
		LogFilePath    string
		RegistryServer string
		Profile        *config.Profile
	}
)

var (
	ErrorCanNotSubmitToNotLeaderNode = errors.New("can not submit to not leader node")
)

func NewNode(profile *config.Profile, shutdown chan string) (*Node, error) {
	logger.Node().Debug("Node {" + profile.Node.NodeId + "} begin init...")

	node := &Node{
		NodeId:             profile.Node.NodeId,
		NodeName:           profile.Node.NodeName,
		onlineSubmits:      make(map[string]*packets.SubmitInfo),
		onlineSubmitLocker: new(sync.RWMutex),
		profile:            profile,
		shutdownChan:       shutdown,
	}

	//init config
	err := node.initConfig(profile)
	if err != nil {
		logger.Node().Debug("Node init config error: " + err.Error())
		return nil, err
	}

	//init cluster
	cluster, err := cluster.NewCluster(
		profile.Cluster.ClusterId,
		profile.Cluster.RegistryServer)
	if err != nil {
		return nil, err
	}
	cluster.OnLeaderChange = node.onLeaderChange

	node.Cluster = cluster

	if node.Config.IsWorker {
		// create runtime
		node.Runtime = runtime.NewRuntime(node.getNodeInfo())
	}

	logger.Node().Debug("Node init success.")
	return node, err
}

func (n *Node) Start() error {
	logger.Node().Debug("Node begin start...")
	// create session with node info
	err := n.createSession(n.getNodeInfo().GetNodeKey(n.Cluster.ClusterId))
	if err != nil {
		return err
	}

	if n.profile.Node.IsMaster {
		n.electionLeader()
	}

	if n.Config.IsWorker {
		go n.Runtime.Start()
	}

	if n.profile.Node.IsMaster {
		n.Cluster.LoadOnlineNodes()
	}

	// register node to cluster
	err = n.registerNode()
	return err
}

func (n *Node) SubmitExecutor(submit *packets.SubmitInfo) (error, *packets.JsonResult) {
	logTitle := "Node SubmitExecutor [" + submit.TaskConfig.TaskID + "] "
	if !n.IsLeader {
		logger.Node().Debug("Node SubmitExecutor [" + submit.TaskConfig.TaskID + "] failed, Current node is not leader.")
		return ErrorCanNotSubmitToNotLeaderNode, nil
	}

	var err error
	// get low balance worker
	if submit.Worker == nil {
		submit.Worker, err = n.Cluster.GetLowBalanceWorker()
		if err != nil {
			logger.Node().Error(err, logTitle+"GetLowBalanceWorker error")
			//TODO log submit result to db log
			return err, nil
		}
	}

	//submit executor to the specified worker node
	rpcClient := n.Cluster.GetRpcClient(submit.Worker.EndPoint())
	err, reply := rpcClient.CallRegisterExecutor(submit.TaskConfig)
	//TODO log submit result to db log
	if err != nil {
		logger.Node().DebugS(logTitle+"to ["+submit.Worker.EndPoint()+"] error:", err.Error())
		return err, reply
	} else {
		if reply.RetCode != reply.CorrectCode() {
			logger.Node().DebugS(logTitle+"to ["+submit.Worker.EndPoint()+"] failed, result:", reply.RetCode)
		} else {
			n.Cluster.Scheduler.AddOnlineSubmit(submit)
			logger.Node().Debug(logTitle + "to [" + submit.Worker.EndPoint() + "] success.")
		}
		return err, reply
	}
}

func (n *Node) Shutdown() {
	logTitle := "Node Shutdown "
	//TODO add some check
	logger.Node().Debug(logTitle + "start.")
	n.shutdownChan <- "ok"
}

// electionLeader
func (n *Node) electionLeader() {
	logTitle := "Node election leader "
	logger.Node().Debug(logTitle + "begin...")

	go func() {
		var retryCount int
		limit := n.profile.Global.RetryLimit
		for {
			if retryCount > limit {
				err := errors.New(logTitle + "retry count bigger than " + strconv.Itoa(limit) + ", now stop it.")
				logger.Node().DebugS(logTitle + "error:" + err.Error())
				return
			}
			retryCount += 1

			err := n.Cluster.ElectionLeader(n.getNodeInfo().EndPoint(), "")
			if err == nil {
				logger.Node().Debug(logTitle + "success with key {" + n.Cluster.LeaderKey + "}")
				n.becomeLeaderRole()
			} else {
				logger.Node().DebugS(logTitle + "error: " + err.Error())
				logger.Node().Error(err, logTitle+"error")
			}
		}
	}()
}

// registerNode register node to cluster
func (n *Node) registerNode() error {
	logTitle := "Node registerNode "
	var leaderServer string
	var err error
	var retryCount int
	nodeInfo := n.getNodeInfo()
	logger.Cluster().Debug(logTitle + "begin...")
RegisterNode:
	for {
		if retryCount > n.profile.Global.RetryLimit {
			err = errors.New(logTitle + "retry count bigger than 5, now stop it")
			logger.Node().DebugS(logTitle + "error: " + err.Error())
			return err
		}
		retryCount += 1
		// get leader info
		leaderServer, err = n.Cluster.GetLeaderInfo()
		if err != nil {
			logger.Node().Debug(logTitle + "GetLeaderInfo error, will retry 10 seconds after.")
			time.Sleep(time.Second * 10)
			continue RegisterNode
		} else {
			logger.Node().Debug(logTitle + "GetLeaderInfo success [" + leaderServer + "]")
			rpcClient := client.NewRpcClient(leaderServer)
			err, result := rpcClient.CallRegisterNode(nodeInfo)
			if err != nil {
				logger.Node().Debug(logTitle + "CallRegisterNode error will retry 10 seconds after.")
				time.Sleep(time.Second * 10)
				continue RegisterNode
			}
			if result.RetCode != result.CorrectCode() {
				logger.Node().Debug(logTitle + "CallRegisterNode error will retry 10 seconds after.")
				time.Sleep(time.Second * 10)
				continue RegisterNode
			} else {
				// watch leader change
				n.watchLeader()
			}
			break
		}
	}
	logger.Node().DebugS(logTitle + "success.")
	return nil
}

// addOnlineSubmit
func (n *Node) addOnlineSubmit(submit *packets.SubmitInfo) {
	n.onlineSubmitLocker.Lock()
	defer n.onlineSubmitLocker.Unlock()
	n.onlineSubmits[submit.TaskConfig.TaskID] = submit
}

// onLeaderChange do something when leader is changed
func (n *Node) onLeaderChange() {
	err := n.registerNode()
	if err != nil {
		logger.Node().DebugS("Node.onLeaderChange registerNode error:", err.Error())
	} else {
		logger.Node().Debug("Node.onLeaderChange registerNode success")
	}
	if n.IsLeader {
		if n.Cluster.LeaderServer != n.getNodeInfo().EndPoint() {
			n.removeLeaderRole()
		}
	}
}

// createSession create session to registry server
func (n *Node) createSession(nodeKey string) error {
	logger.Node().Debug("Node create session begin...")
	err := n.Cluster.CreateSession(nodeKey, n.getNodeInfo())
	if err != nil {
		logger.Node().Debug("Node create session error: " + err.Error())
	} else {
		logger.Node().Debug("Node create session success with key {" + nodeKey + "}")
	}
	return err
}

func (n *Node) becomeLeaderRole() {
	logTitle := "Node "
	//TODO do something when become to leader
	logger.Node().Debug(logTitle + "become to leader role")
	n.IsLeader = true

}

func (n *Node) removeLeaderRole() {
	logTitle := "Node "
	//TODO do something when become to not leader
	logger.Node().Debug(logTitle + "remove leader role")
	n.IsLeader = false
}

func (n *Node) watchLeader() {
	logTitle := "Node.watchLeader "
	logger.Cluster().Debug(logTitle + "running.")
	go func() {
		var retryCount int
		for {
			err := n.Cluster.WatchLeader()
			if err != nil {
				if retryCount > config.CurrentProfile.Cluster.WatchLeaderRetryLimit {
					logger.Cluster().DebugS(logTitle + "error num bigger than max limit, now will shutdown node.")
					n.Shutdown()
				} else {
					retryCount += 1
					logger.Cluster().DebugS(logTitle+"error, will retry after 10 seconds:", err.Error())
				}
				time.Sleep(time.Second * 10)
			}
		}
	}()
}

// initConfig init node config from config profile
func (n *Node) initConfig(conf *config.Profile) error {
	n.Config = new(NodeConfig)
	n.Config.RegistryServer = conf.Cluster.RegistryServer
	n.Config.IsMaster = conf.Node.IsMaster
	n.Config.IsWorker = conf.Node.IsWorker
	n.Config.Profile = conf

	logger.Node().Debug("Node Config init success.")
	return nil
}

func (n *Node) getNodeInfo() *packets.NodeInfo {
	nodeInfo := &packets.NodeInfo{
		NodeID:    n.NodeId,
		Cluster:   n.profile.Cluster.ClusterId,
		OuterHost: n.profile.Rpc.OuterHost,
		OuterPort: n.profile.Rpc.OuterPort,
		Host:      n.profile.Rpc.RpcHost,
		Port:      n.profile.Rpc.RpcPort,
		IsMaster:  n.profile.Node.IsMaster,
		IsWorker:  n.profile.Node.IsWorker,
		IsOnline:  true,
	}
	return nodeInfo
}

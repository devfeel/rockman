package node

import (
	"errors"
	"github.com/devfeel/rockman/cluster"
	"github.com/devfeel/rockman/config"
	"github.com/devfeel/rockman/core"
	"github.com/devfeel/rockman/logger"
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
		isLeader           bool
		Status             int
		config             *config.Profile
		nodeInfo           *core.NodeInfo
		onlineSubmits      map[string]*core.SubmitInfo
		onlineSubmitLocker *sync.RWMutex
		Cluster            *cluster.Cluster
		Runtime            *runtime.Runtime
		shutdownChan       chan string
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
		onlineSubmits:      make(map[string]*core.SubmitInfo),
		onlineSubmitLocker: new(sync.RWMutex),
		config:             profile,
		shutdownChan:       shutdown,
	}

	//init cluster
	cluster, err := cluster.NewCluster(profile)
	if err != nil {
		return nil, err
	}
	cluster.OnLeaderChange = node.onLeaderChange
	cluster.OnLeaderChangeFailed = node.onLeaderChangeFailed

	node.Cluster = cluster
	if node.config.Node.IsWorker {
		// create runtime
		node.Runtime = runtime.NewRuntime(node.NodeInfo())
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

	if n.config.Node.IsMaster {
		err = n.Cluster.Start()
		if err != nil {
			return err
		}
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

func (n *Node) SubmitExecutor(submit *core.SubmitInfo) *core.Result {
	logTitle := "Node SubmitExecutor [" + submit.TaskConfig.TaskID + "] "
	if !n.IsLeader() {
		logger.Node().Debug(logTitle + "failed, current node is not leader.")
		return core.CreateResult(-1001, "current node is not leader", nil)
	}

	if submit.Worker != nil {
		if submit.Worker.Cluster != n.Cluster.ClusterId {
			logger.Node().Debug(logTitle + "failed, not match cluster [" + submit.Worker.Cluster + ", " + n.Cluster.ClusterId + "]")
			return core.CreateResult(-1002, "not match cluster ["+submit.Worker.Cluster+", "+n.Cluster.ClusterId+"]", nil)
		}

		endPoint := submit.Worker.EndPoint()
		node, exists := n.Cluster.FindNode(endPoint)
		if !exists {
			logger.Node().Debug(logTitle + "failed, can not find node[" + endPoint + "] in cluster")
			return core.CreateResult(-1003, "can not find node["+endPoint+"] in cluster", nil)
		}

		if node.NodeID != submit.Worker.NodeID {
			logger.Node().Debug(logTitle + "failed, not match node id [" + submit.Worker.NodeID + ", " + node.NodeID + "]")
			return core.CreateResult(-1004, "not match node id ["+submit.Worker.NodeID+", "+node.NodeID+"]", nil)
		}
	}

	var err error
	// get low balance worker
	if submit.Worker == nil {
		submit.Worker, err = n.Cluster.GetLowBalanceWorker()
		if err != nil {
			logger.Node().Error(err, logTitle+"GetLowBalanceWorker error")
			//TODO log submit result to db log
			return core.CreateErrorResult(err)
		}
	}

	//submit executor to the specified worker node
	rpcClient := n.Cluster.GetRpcClient(submit.Worker.EndPoint())
	err, reply := rpcClient.CallRegisterExecutor(submit.TaskConfig)
	//TODO log submit result to db log
	if err != nil {
		logger.Node().DebugS(logTitle+"to ["+submit.Worker.EndPoint()+"] error:", err.Error())
		return core.CreateErrorResult(err)
	} else {
		if reply.IsSuccess() {
			logger.Node().DebugS(logTitle+"to ["+submit.Worker.EndPoint()+"] failed, result:", reply.RetCode)
		} else {
			n.Cluster.Scheduler.AddOnlineSubmit(submit)
			logger.Node().Debug(logTitle + "to [" + submit.Worker.EndPoint() + "] success.")
		}
		return core.CreateResult(reply.RetCode, reply.RetMsg, nil)
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
	logger.Node().Debug(logTitle + "begin.")

	go func() {
		var retryCount int
		limit := n.config.Global.RetryLimit
		for {
			err := n.Cluster.ElectionLeader(n.NodeInfo().EndPoint())
			if err != nil {
				retryCount += 1
				if retryCount > limit {
					err := errors.New("retry count bigger than " + strconv.Itoa(limit) + ", now will stop node.")
					logger.Node().DebugS(logTitle + "error:" + err.Error())
					n.Shutdown()
					return
				}
				logger.Node().DebugS(logTitle + "error: " + err.Error())
				logger.Node().Error(err, logTitle+"error")
			} else {
				logger.Node().Debug(logTitle + "success with key {" + n.Cluster.LeaderKey + "}")
				n.becomeLeaderRole()
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
	nodeInfo := n.NodeInfo()
	logger.Cluster().Debug(logTitle + "begin.")
RegisterNode:
	for {
		if retryCount > n.config.Global.RetryLimit {
			err = errors.New(logTitle + "retry more than 5 times, now will be stop")
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
			}
			break
		}
	}
	logger.Node().DebugS(logTitle + "success.")
	return nil
}

// addOnlineSubmit
func (n *Node) addOnlineSubmit(submit *core.SubmitInfo) {
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
	if n.IsLeader() {
		if n.Cluster.LeaderServer != n.NodeInfo().EndPoint() {
			n.removeLeaderRole()
		}
	}

	if n.IsMaster() && !n.IsLeader() {
		if n.Cluster.LeaderServer == n.NodeInfo().EndPoint() {
			n.becomeLeaderRole()
		}
	}
}

// onLeaderChangeFailed
func (n *Node) onLeaderChangeFailed() {
	logger.Node().DebugS("Node.onLeaderChangeFailed, now will shutdown node.")
	n.Shutdown()
}

// createSession create session to registry server
func (n *Node) createSession(nodeKey string) error {
	logger.Node().Debug("Node create session begin.")
	err := n.Cluster.CreateSession(nodeKey, n.NodeInfo())
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
	n.isLeader = true

}

func (n *Node) removeLeaderRole() {
	logTitle := "Node "
	//TODO do something when become to not leader
	logger.Node().Debug(logTitle + "remove leader role")
	n.isLeader = false
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
	return n.nodeInfo
}

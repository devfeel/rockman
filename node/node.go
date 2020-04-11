package node

import (
	"errors"
	"fmt"
	"github.com/devfeel/rockman/cluster"
	"github.com/devfeel/rockman/config"
	"github.com/devfeel/rockman/core"
	"github.com/devfeel/rockman/logger"
	"github.com/devfeel/rockman/protected/model"
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
		NodeId       string
		NodeName     string
		isLeader     bool
		Status       int
		config       *config.Profile
		nodeInfo     *core.NodeInfo
		Cluster      *cluster.Cluster
		Registry     *registry.Registry
		Runtime      *runtime.Runtime
		shutdownChan chan string
		isSTW        bool //stop the world flag
		logLogic     *service.LogService
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

	n.Registry.Start()

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

func (n *Node) ClusterId() string {
	return n.Cluster.ClusterId
}

func (n *Node) SubmitExecutor(execInfo *core.ExecutorInfo) *core.Result {
	logTitle := "Node SubmitExecutor [" + execInfo.TaskConfig.TaskID + "] "
	if !n.IsLeader() {
		logger.Node().Debug(logTitle + "failed, current node is not leader.")
		return core.FailedResult(-1001, "current node is not leader")
	}

	if execInfo.Worker != nil {
		if execInfo.Worker.Cluster != n.Cluster.ClusterId {
			logger.Node().Debug(logTitle + "failed, not match cluster [" + execInfo.Worker.Cluster + ", " + n.Cluster.ClusterId + "]")
			return core.FailedResult(-1002, "not match cluster ["+execInfo.Worker.Cluster+", "+n.Cluster.ClusterId+"]")
		}

		endPoint := execInfo.Worker.EndPoint()
		node, exists := n.Cluster.FindNode(endPoint)
		if !exists {
			logger.Node().Debug(logTitle + "failed, can not find node[" + endPoint + "] in cluster")
			return core.FailedResult(-1003, "can not find node["+endPoint+"] in cluster")
		}

		if node.NodeID != execInfo.Worker.NodeID {
			logger.Node().Debug(logTitle + "failed, not match node id [" + execInfo.Worker.NodeID + ", " + node.NodeID + "]")
			return core.FailedResult(-1004, "not match node id ["+execInfo.Worker.NodeID+", "+node.NodeID+"]")
		}
	}

	var err error
	// get low balance worker
	if execInfo.Worker == nil {
		execInfo.Worker, err = n.Cluster.GetLowBalanceWorker()
		if err != nil {
			logger.Node().Error(err, logTitle+"GetLowBalanceWorker error")
			//TODO log submit result to db log
			return core.ErrorResult(err)
		}
	}

	//submit executor to the specified worker node
	rpcClient := n.Cluster.GetRpcClient(execInfo.Worker.EndPoint())
	err, reply := rpcClient.CallRegisterExecutor(execInfo.TaskConfig)
	submitLog := &model.TaskSubmitLog{
		TaskID:       execInfo.TaskConfig.TaskID,
		NodeID:       execInfo.Worker.NodeID,
		NodeEndPoint: execInfo.Worker.EndPoint(),
		IsSuccess:    err != nil && reply.IsSuccess(),
	}

	if err != nil {
		logger.Node().DebugS(logTitle+"to ["+execInfo.Worker.EndPoint()+"] error:", err.Error())
		submitLog.FailureType = "error"
		submitLog.FailureCause = err.Error()
		n.logLogic.WriteSubmitLog(submitLog)
		return core.ErrorResult(err)
	} else {
		if !reply.IsSuccess() {
			submitLog.FailureType = "failure"
			submitLog.FailureCause = reply.FailureMessage()
			logger.Node().DebugS(logTitle+"to ["+execInfo.Worker.EndPoint()+"] failed, result:", reply.RetCode)
		} else {
			n.Cluster.AddExecutor(execInfo)
			submitLog.IsSuccess = true
			logger.Node().Debug(logTitle + "to [" + execInfo.Worker.EndPoint() + "] success.")
		}
		n.logLogic.WriteSubmitLog(submitLog)
		return core.NewResult(reply.RetCode, reply.RetMsg, nil)
	}
}

func (n *Node) SubmitStopExecutor(taskId string) *core.Result {
	logTitle := "Node SubmitStopExecutor [" + taskId + "] "
	if !n.IsLeader() {
		logger.Node().Debug(logTitle + "failed, current node is not leader.")
		return core.FailedResult(-1001, "current node is not leader")
	}

	runExecInfo, exists := n.Cluster.FindExecutor(taskId)
	if !exists {
		logger.Node().Debug(logTitle + "failed, can not find executor is running cluster.")
		return core.FailedResult(-1001, "can not find executor is running cluster")
	}

	if !runExecInfo.TaskConfig.IsRun {
		return core.SuccessResult()
	}

	//submit executor to the specified worker node
	rpcClient := n.Cluster.GetRpcClient(runExecInfo.Worker.EndPoint())
	err, reply := rpcClient.CallStopExecutor(taskId)
	//TODO log submit result to db log
	if err != nil {
		logger.Node().DebugS(logTitle+"to ["+runExecInfo.Worker.EndPoint()+"] error:", err.Error())
		return core.ErrorResult(err)
	} else {
		if !reply.IsSuccess() {
			logger.Node().DebugS(logTitle+"to ["+runExecInfo.Worker.EndPoint()+"] failed, result:", reply.RetCode)
		} else {
			runExecInfo.TaskConfig.IsRun = false
			logger.Node().Debug(logTitle + "to [" + runExecInfo.Worker.EndPoint() + "] success.")
		}
		return core.NewResult(reply.RetCode, reply.RetMsg, nil)
	}
}

func (n *Node) SubmitStartExecutor(taskId string) *core.Result {
	logTitle := "Node SubmitStartExecutor [" + taskId + "] "
	if !n.IsLeader() {
		logger.Node().Debug(logTitle + "failed, current node is not leader.")
		return core.FailedResult(-1001, "current node is not leader")
	}

	runExecInfo, exists := n.Cluster.FindExecutor(taskId)
	if !exists {
		logger.Node().Debug(logTitle + "failed, can not find executor is running cluster.")
		return core.FailedResult(-1001, "can not find executor is running cluster")
	}

	if runExecInfo.TaskConfig.IsRun {
		return core.SuccessResult()
	}

	//submit executor to the specified worker node
	rpcClient := n.Cluster.GetRpcClient(runExecInfo.Worker.EndPoint())
	err, reply := rpcClient.CallStartExecutor(runExecInfo.TaskConfig.TaskID)
	//TODO log submit result to db log
	if err != nil {
		logger.Node().DebugS(logTitle+"to ["+runExecInfo.Worker.EndPoint()+"] error:", err.Error())
		return core.ErrorResult(err)
	} else {
		if !reply.IsSuccess() {
			logger.Node().DebugS(logTitle+"to ["+runExecInfo.Worker.EndPoint()+"] failed, result:", reply.RetCode)
		} else {
			runExecInfo.TaskConfig.IsRun = true
			logger.Node().Debug(logTitle + "to [" + runExecInfo.Worker.EndPoint() + "] success.")
		}
		return core.NewResult(reply.RetCode, reply.RetMsg, nil)
	}
}

// RegisterExecutor
func (n *Node) RegisterExecutor(taskInfo *core.TaskConfig) *core.Result {
	logTitle := "Node RegisterExecutor [" + taskInfo.TaskID + "] "
	if !n.IsWorker() {
		logger.Node().Debug(logTitle + "failed, current node is not worker.")
		return core.FailedResult(-1001, "current node is not worker")
	}
	_, err := n.Runtime.CreateExecutor(taskInfo)
	if err != nil {
		logger.Node().Warn(logTitle + "CreateExecutor error:" + err.Error())
		return core.FailedResult(-2001, err.Error())
	} else {
		// reg to registry server
		execInfo := new(core.ExecutorInfo)
		execInfo.TaskConfig = taskInfo
		execInfo.Worker = n.NodeInfo()
		_, err := n.Registry.Set(execInfo.GetExecutorKey(n.ClusterId()), execInfo.Json(), nil)
		if err != nil {
			logger.Node().Warn(logTitle + "sync to registry error:" + err.Error())
		}
		return core.SuccessResult()
	}
}

func (n *Node) Shutdown() {
	logTitle := "Node Shutdown "
	//TODO add some check
	logger.Node().Debug(logTitle + "start.")
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

// electionLeader
func (n *Node) electionLeader() {
	logTitle := "Node election leader "
	logger.Node().Debug(logTitle + "begin.")

	doQuery := func() error {
		err := n.Cluster.ElectionLeader(n.NodeInfo().EndPoint())
		if err != nil {
			logger.Node().DebugS(logTitle + "error: " + err.Error() + ", will retry 10 seconds after")
			logger.Node().Error(err, logTitle+"error")
			time.Sleep(time.Second * 10)
			return err
		} else {
			logger.Node().Debug(logTitle + "success with key {" + n.Cluster.LeaderKey + "}")
			n.becomeLeaderRole()
			return nil
		}
	}

	go func() {
		var retryCount int
		limit := n.config.Global.RetryLimit
		for {
			if n.isSTW {
				return
			}
			err := doQuery()
			if err != nil {
				retryCount += 1
				if retryCount > limit {
					err := errors.New("retry count bigger than " + strconv.Itoa(limit) + ", now will stop node.")
					logger.Node().DebugS(logTitle + "error:" + err.Error())
					n.Shutdown()
					return
				}
			} else {
				retryCount = 0
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
				logger.Node().DebugS(logTitle + "success.")
				n.initExecutorsFromDB()
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

func (n *Node) becomeLeaderRole() {
	logTitle := "Node.becomeLeaderRole "
	logger.Node().Debug(logTitle + "become to leader role")
	n.isLeader = true

	//TODO sync all executors
	n.Cluster.OnNodeOffline = n.onWorkerNodeOffline

}

func (n *Node) removeLeaderRole() {
	logTitle := "Node "
	//TODO do something when become to not leader
	logger.Node().Debug(logTitle + "remove leader role")
	n.Cluster.OnNodeOffline = nil
	n.isLeader = false
}

// initExecutorsFromDB init executors from db
// must check init flag on registry
func (n *Node) initExecutorsFromDB() {
	logTitle := "Node initExecutorsFromDB "
	if !n.IsLeader() {
		return
	}
	var successCount, failureCount int
	doQuery := func() {
		defer func() {
			if err := recover(); err != nil {
				errInfo := errors.New(fmt.Sprintln(err))
				logger.Cluster().Error(errInfo, logTitle+"throw unhandled error:"+errInfo.Error())
			}
		}()
		execInfos, err := service.NewExecutorService().QueryAllExecutors()
		if err != nil {
			logger.Node().Debug(logTitle + "NewExecutorService error:" + err.Error())
			return
		}
		if execInfos == nil {
			return
		}
		for _, exec := range execInfos {
			submit := new(core.ExecutorInfo)
			submit.TaskConfig = exec.TaskConfig()
			if submit.TaskConfig == nil || submit.TaskConfig.TargetConfig == nil {
				logger.Node().Debug(logTitle + "init submit error: TaskConfig is nil or target config is nil")
				failureCount += 1
				continue
			}
			result := n.SubmitExecutor(submit)
			if result.Error != nil {
				logger.Node().DebugS(logTitle+"SubmitExecutor error:", result.Error.Error())
				failureCount += 1
				//TODO log to db
				continue
			}

			if !result.IsSuccess() {
				logger.Node().DebugS(logTitle + "SubmitExecutor failed, " + result.Message())
				failureCount += 1
				//TODO log to db
			} else {
				logger.Node().DebugS(logTitle + "SubmitExecutor success")
				successCount += 1
				//TODO log to db
			}
		}
	}

	logger.Node().Debug(logTitle + "init begin.")
	flag, err := n.getInitFlag()
	if err != nil {
		logger.Node().Warn(logTitle + "get init flag error:" + err.Error())
	} else {
		if !flag {
			doQuery()
			err := n.setInitFlag()
			if err != nil {
				logger.Node().Warn(logTitle + "set init flag error:" + err.Error())
			}
			logger.Node().Debug(logTitle + "init finish. Success[" + strconv.Itoa(successCount) + "] Failure[" + strconv.Itoa(failureCount) + "]")
		}
	}
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
	for _, v := range n.Cluster.Executors {
		if v.Worker.NodeID == nodeInfo.NodeID {
			needReSubmits = append(needReSubmits, v)
		}
	}
	go func() {
		for _, exec := range needReSubmits {
			result := n.SubmitExecutor(exec)
			if result.Error != nil {
				logger.Node().DebugS(logTitle+"HA SubmitExecutor error:", result.Error.Error())
				//TODO log to db
			} else {
				if !result.IsSuccess() {
					logger.Node().DebugS(logTitle + "HA SubmitExecutor failed, " + result.Message())
					//TODO log to db
				} else {
					logger.Node().DebugS(logTitle + "HA SubmitExecutor success")
					//TODO log to db
				}
			}
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

func (n *Node) getInitFlag() (bool, error) {
	kv, _, err := n.Registry.Get(getInitFlagKey(n.ClusterId()), nil)
	if err == nil {
		return false, err
	}
	if kv == nil {
		return false, nil
	}
	return true, nil
}

func (n *Node) setInitFlag() error {
	_, err := n.Registry.Set(getInitFlagKey(n.ClusterId()), "true", nil)
	return err
}

func getInitFlagKey(clusterId string) string {
	return core.ClusterKeyPrefix + clusterId + "/flags/init"
}

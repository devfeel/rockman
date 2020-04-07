package cluster

import (
	"errors"
	"fmt"
	"github.com/devfeel/mapper"
	"github.com/devfeel/rockman/config"
	"github.com/devfeel/rockman/core"
	"github.com/devfeel/rockman/logger"
	"github.com/devfeel/rockman/registry"
	"github.com/devfeel/rockman/rpc/client"
	"github.com/devfeel/rockman/scheduler"
	"github.com/hashicorp/consul/api"
	"strconv"
	"sync"
	"time"
)

const (
	MinQueryResourceInterval = 60
)

type (
	Cluster struct {
		ClusterId             string
		Registry              *registry.Registry
		LeaderKey             string
		LeaderServer          string
		leaderLastIndex       uint64
		lastGetLeaderTime     time.Time
		Nodes                 map[string]*core.NodeInfo
		nodesLastIndex        uint64
		nodesLocker           *sync.RWMutex
		lastLoadNodesTime     time.Time
		Executors             map[string]*core.ExecutorInfo
		executorsLastIndex    uint64
		executorsLocker       *sync.RWMutex
		lastLoadExecutorsTime time.Time
		rpcClients            map[string]*client.RpcClient
		rpcClientLocker       *sync.RWMutex
		Scheduler             *scheduler.Scheduler
		config                *config.Profile
		isSTW                 bool //stop the world flag
		OnNodesChange         WatchChangeHandle
		OnNodeOffline         NodeOfflineHandle
		OnLeaderChange        WatchChangeHandle
		OnLeaderChangeFailed  WatchFailedHandle
	}

	WatchChangeHandle func()
	WatchFailedHandle func()
	NodeOfflineHandle func(info *core.NodeInfo)
)

// NewCluster new cluster and reg server
func NewCluster(profile *config.Profile, registry *registry.Registry) *Cluster {
	cluster := new(Cluster)
	cluster.config = profile
	cluster.ClusterId = profile.Cluster.ClusterId
	cluster.LeaderKey = getLeaderKey(profile.Cluster.ClusterId)
	cluster.Nodes = make(map[string]*core.NodeInfo)
	cluster.nodesLocker = new(sync.RWMutex)
	cluster.Executors = make(map[string]*core.ExecutorInfo)
	cluster.executorsLocker = new(sync.RWMutex)
	cluster.rpcClients = make(map[string]*client.RpcClient)
	cluster.rpcClientLocker = new(sync.RWMutex)

	cluster.Registry = registry

	cluster.Scheduler = scheduler.NewScheduler()
	logger.Node().Debug("Cluster init success.")
	return cluster
}

func (c *Cluster) Start() error {
	logger.Default().Debug("Cluster start...")
	err := c.loadOnlineNodes()
	if err != nil {
		return err
	}
	c.watchLeader()
	c.watchOnlineNodes()
	c.cycleQueryWorkerResource()
	return nil
}

func (c *Cluster) Stop() error {
	logger.Default().Debug("Cluster Stop.")
	c.isSTW = true
	return nil
}

// electionLeader election leader role to registry server
func (c *Cluster) ElectionLeader(leaderServer string) error {
	opts := &api.LockOptions{
		Key:   c.LeaderKey,
		Value: []byte(leaderServer),
		SessionOpts: &api.SessionEntry{
			Name:     leaderServer,
			TTL:      "10s",
			Behavior: "delete",
		},
	}
	locker, err := c.Registry.CreateLockerOpts(opts)
	if err != nil {
		return err
	}

	_, err = locker.Locker.Lock(nil)
	if err != nil {
		return err
	}
	return nil
}

// GetLeaderInfo get leader info from leader key
// must check is locked by leader session
// leader changed by watchLeaderChange
func (c *Cluster) GetLeaderInfo() (string, error) {
	if c.LeaderServer != "" {
		return c.LeaderServer, nil
	}
	kvPair, meta, err := c.Registry.Get(c.LeaderKey, nil)
	if err != nil {
		return "", err
	} else {
		if kvPair.Session == "" {
			return "", errors.New("lock session is nil")
		} else {
			c.LeaderServer = string(kvPair.Value)
			c.leaderLastIndex = meta.LastIndex
			c.lastGetLeaderTime = time.Now()
			return c.LeaderServer, nil
		}
	}
}

// AddNodeInfo add node into node list
// it will query remote resource
func (c *Cluster) AddNodeInfo(nodeInfo *core.NodeInfo) *core.Result {
	if nodeInfo.Cluster != c.ClusterId {
		return core.FailedResult(-1001, "not match cluster")
	}
	key := nodeInfo.EndPoint()
	resource, result := c.QueryNodeResource(key)
	if result.Error != nil {
		resource = nodeInfo.GetEmptyResource()
		logger.Cluster().Warn("AddNodeInfo.QueryResource[" + key + "] error: " + result.Error.Error())
		logger.Cluster().Error(result.Error, "AddNodeInfo.QueryResource["+key+"] error")
	} else {
		if !result.IsSuccess() {
			resource = nodeInfo.GetEmptyResource()
			logger.Cluster().Warn("AddNodeInfo.QueryResource[" + key + "] failed: " + result.Message())
		}
	}
	c.nodesLocker.Lock()
	defer c.nodesLocker.Unlock()
	c.Scheduler.SetResource(resource)
	c.Nodes[key] = nodeInfo
	return core.SuccessResult()
}

// AddExecutor add executor info into executor list
func (c *Cluster) AddExecutor(execInfo *core.ExecutorInfo) *core.Result {
	if execInfo.Worker.Cluster != c.ClusterId {
		return core.FailedResult(-1001, "not match cluster")
	}
	c.executorsLocker.Lock()
	defer c.executorsLocker.Unlock()
	c.Executors[execInfo.TaskConfig.TaskID] = execInfo
	return core.SuccessResult()
}

// FindNode find node info by endpoint
func (c *Cluster) FindNode(endPoint string) (*core.NodeInfo, bool) {
	c.nodesLocker.RLock()
	defer c.nodesLocker.RUnlock()
	node, exists := c.Nodes[endPoint]
	return node, exists
}

// GetRpcClient get rpc client with endpoint
func (c *Cluster) GetRpcClient(endPoint string) *client.RpcClient {
	//TODO check endpoint is in cluster
	defer c.rpcClientLocker.Unlock()
	c.rpcClientLocker.Lock()
	var rpcClient *client.RpcClient
	var isExists bool
	if rpcClient, isExists = c.rpcClients[endPoint]; !isExists {
		rpcClient = client.NewRpcClient(endPoint, c.config.Rpc.EnableTls, c.config.Rpc.ClientCertFile, c.config.Rpc.ClientKeyFile)
		c.rpcClients[endPoint] = rpcClient
	}
	return rpcClient
}

// GetLeaderRpcClient get leader rpc client
func (c *Cluster) GetLeaderRpcClient() *client.RpcClient {
	return c.GetRpcClient(c.LeaderServer)
}

// GetLowBalanceWorker get lower balance worker, if not match, it will try 3 times
func (c *Cluster) GetLowBalanceWorker() (*core.NodeInfo, error) {
	resources, err := c.Scheduler.Schedule(scheduler.Balance_LowerLoad)
	if err != nil {
		return nil, err
	}

	c.nodesLocker.RLock()
	defer c.nodesLocker.RUnlock()

	resource := resources[0]
	rawWorker, isExists := c.Nodes[resource.EndPoint]
	if isExists {
		return rawWorker, nil
	}
	logger.Cluster().Debug("try get lower load worker[" + resource.EndPoint + "] failed 1 times, try get next")
	if len(resources) > 1 {
		resource := resources[1]
		rawWorker, isExists := c.Nodes[resource.EndPoint]
		if isExists {
			return rawWorker, nil
		}
	}
	logger.Cluster().Debug("try get lower load worker[" + resource.EndPoint + "] failed 2 times, try get next.")
	if len(resources) > 2 {
		resource := resources[2]
		rawWorker, isExists := c.Nodes[resource.EndPoint]
		if isExists {
			return rawWorker, nil
		}
	}
	logger.Cluster().Debug("try get lower load worker[" + resource.EndPoint + "] failed 3 times.")
	return nil, errors.New("no match resource with worker")
}

// QueryNodeResource query node resource by endpoint
func (c *Cluster) QueryNodeResource(endPoint string) (*core.ResourceInfo, *core.Result) {
	client := c.GetRpcClient(endPoint)
	err, reply := client.CallQueryResource()
	if err != nil {
		return nil, core.ErrorResult(err)
	} else {
		if !reply.IsSuccess() {
			return nil, core.FailedResult(-1001, "query failed["+strconv.Itoa(reply.RetCode)+", "+reply.RetMsg+"]")
		} else {
			resource := new(core.ResourceInfo)
			err := mapper.MapperMap(reply.Message.(map[string]interface{}), resource)
			if err != nil {
				return nil, core.NewResult(core.ErrorCode, err.Error(), err)
			}
			return resource, core.SuccessResult()
		}
	}
}

// ClusterInfo return ClusterInfo
func (c *Cluster) ClusterInfo() *core.ClusterInfo {
	return &core.ClusterInfo{
		ClusterId:             c.ClusterId,
		RegistryServerUrl:     c.Registry.ServerUrl,
		LeaderKey:             c.LeaderKey,
		LeaderServer:          c.LeaderServer,
		NodeNum:               len(c.Nodes),
		WatchLeaderRetryLimit: c.config.Cluster.QueryResourceInterval,
		QueryResourceInterval: c.config.Cluster.QueryResourceInterval,
	}
}

// loadOnlineNodes load all online nodes from Registry
func (c *Cluster) loadOnlineNodes() error {
	logTitle := "Cluster.loadOnlineNodes "
	nodeKVs, meta, err := c.Registry.ListKV(core.GetNodeKeyPrefix(c.ClusterId), nil)
	if err != nil {
		logger.Cluster().Debug(logTitle + "error: " + err.Error())
		return errors.New(logTitle + "error: " + err.Error())
	} else {
		logger.Cluster().Debug(logTitle + "success.")
	}
	c.nodesLastIndex = meta.LastIndex
	c.refreshOnlineNodes(nodeKVs)
	return nil
}

// refreshOnlineNodes
func (c *Cluster) refreshOnlineNodes(nodeKVs api.KVPairs) {
	nodes := make(map[string]*core.NodeInfo)
	for _, s := range nodeKVs {
		if s.Session == "" {
			continue
		}
		node := new(core.NodeInfo)
		if err := node.LoadFromJson(string(s.Value)); err != nil {
			continue
		}
		if node.Cluster != c.ClusterId {
			continue
		}
		nodes[node.EndPoint()] = node
		if _, exists := c.Nodes[node.EndPoint()]; !exists {
			c.AddNodeInfo(node)
		}
	}
	c.nodesLocker.Lock()
	defer c.nodesLocker.Unlock()
	for _, node := range c.Nodes {
		if _, exists := nodes[node.EndPoint()]; !exists {
			node.IsOnline = false
			if c.OnNodeOffline != nil {
				c.OnNodeOffline(node)
			}
		} else {
			node.IsOnline = true
		}
	}
	c.lastLoadNodesTime = time.Now()
}

// watchOnlineNodes watch online nodes change
func (c *Cluster) watchOnlineNodes() {
	logTitle := "Cluster.watchOnlineNodes "
	logger.Cluster().Debug(logTitle + "running...")
	doQuery := func() error {
		defer func() {
			if err := recover(); err != nil {
				errInfo := errors.New(fmt.Sprintln(err))
				logger.Cluster().Error(errInfo, logTitle+"throw unhandled error:"+errInfo.Error())
			}
		}()

		opt := &api.QueryOptions{
			WaitIndex: c.nodesLastIndex,
			WaitTime:  time.Minute * 10,
		}
		nodeKVs, meta, err := c.Registry.ListKV(core.GetNodeKeyPrefix(c.ClusterId), opt)
		if err != nil {
			return err
		}
		if meta.LastIndex != c.nodesLastIndex {
			logger.Cluster().Debug(logTitle + "some nodes changed.")
			c.nodesLastIndex = meta.LastIndex
			c.refreshOnlineNodes(nodeKVs)
			if c.OnNodesChange != nil {
				c.OnNodesChange()
			}
		}
		return nil
	}

	go func() {
		for {
			if c.isSTW {
				return
			}
			err := doQuery()
			if err != nil {
				logger.Cluster().DebugS(logTitle+"error, will retry after 10 seconds:", err.Error())
				time.Sleep(time.Second * 10)
			}
		}
	}()
}

// watchLeader watch leader change
func (c *Cluster) watchLeader() {
	logTitle := "Cluster.watchLeader "
	logger.Cluster().Debug(logTitle + "running...")

	doQuery := func() error {
		defer func() {
			if err := recover(); err != nil {
				errInfo := errors.New(fmt.Sprintln(err))
				logger.Cluster().Error(errInfo, logTitle+"throw unhandled error:"+errInfo.Error())
			}
		}()

		opt := &api.QueryOptions{
			WaitIndex: c.leaderLastIndex,
			WaitTime:  time.Minute * 10,
		}
		kvPair, meta, err := c.Registry.Get(c.LeaderKey, opt)
		if err != nil {
			return err
		}
		if kvPair.Session == "" {
			return errors.New("leader lock session is nil")
		}
		if meta.LastIndex != c.leaderLastIndex {
			logger.Cluster().Debug("Cluster.watchLeaderChange: leader changed.")
			c.leaderLastIndex = meta.LastIndex
			c.LeaderServer = string(kvPair.Value)
			c.lastGetLeaderTime = time.Now()
			if c.OnLeaderChange != nil {
				c.OnLeaderChange()
			}
		}
		return nil
	}
	go func() {
		var retryCount int
		for {
			if c.isSTW {
				return
			}
			retryWaitSeconds := (retryCount + 1) * 10
			err := doQuery()
			if err != nil {
				retryCount += 1
				if retryCount > config.CurrentProfile.Cluster.WatchLeaderRetryLimit {
					logger.Cluster().DebugS(logTitle + "error count bigger than max limit")
					if c.OnLeaderChangeFailed != nil {
						c.OnLeaderChangeFailed()
					}
				} else {
					logger.Cluster().DebugS(logTitle+"error, will retry after "+strconv.Itoa(retryWaitSeconds)+" seconds:", err.Error())
				}
				time.Sleep(time.Second * time.Duration(retryWaitSeconds))
			}
		}
	}()
}

// CycleLoadWorkerResource
func (c *Cluster) cycleQueryWorkerResource() {
	logTitle := "Cluster.cycleQueryWorkerResource "
	logger.Cluster().Debug(logTitle + "running...")
	doQuery := func() {
		defer func() {
			if err := recover(); err != nil {
				errInfo := errors.New(fmt.Sprintln(err))
				logger.Cluster().Error(errInfo, logTitle+"throw unhandled error:"+errInfo.Error())
			}
		}()
		queryNodes := 0
		failedNodes := 0
		for _, n := range c.Nodes {
			if !n.IsOnline && !n.IsWorker {
				continue
			}
			queryNodes += 1
			resource, result := c.QueryNodeResource(n.EndPoint())
			if result.Error != nil {
				failedNodes += 1
				logger.Cluster().Warn(logTitle + "QueryResource[" + n.EndPoint() + "] error: " + result.Message())
				logger.Cluster().Error(result.Error, logTitle+"QueryResource["+n.EndPoint()+"] error")
				continue
			} else {
				if !result.IsSuccess() {
					logger.Cluster().Warn(logTitle + "QueryResource[" + n.EndPoint() + "] failed: " + result.Message())
					continue
				}
				c.Scheduler.SetResource(resource)
			}
		}
		logger.Cluster().Debug(logTitle + "success, query nodes[" + strconv.Itoa(queryNodes) + "], failed[" + strconv.Itoa(failedNodes) + "]")
	}

	go func() {
		for {
			interval := c.config.Cluster.QueryResourceInterval
			if interval < MinQueryResourceInterval {
				interval = MinQueryResourceInterval
			}
			time.Sleep(time.Second * time.Duration(interval))

			if c.isSTW {
				return
			}
			doQuery()
		}
	}()
}

func getLeaderKey(clusterId string) string {
	return "devfeel/rockman/" + clusterId + "/leader"
}

package cluster

import (
	"errors"
	"fmt"
	"github.com/devfeel/rockman/cluster/consul"
	"github.com/devfeel/rockman/config"
	"github.com/devfeel/rockman/core"
	"github.com/devfeel/rockman/logger"
	"github.com/devfeel/rockman/rpc/client"
	"github.com/devfeel/rockman/scheduler"
	"github.com/hashicorp/consul/api"
	"sync"
	"time"
)

type (
	Cluster struct {
		ClusterId         string
		RegistryServerUrl string
		RegistryClient    *consul.ConsulClient
		LeaderKey         string
		LeaderServer      string
		leaderLastIndex   uint64
		lastGetLeaderTime time.Time
		OnLeaderChange    WatchChangeHandle
		Nodes             map[string]*core.NodeInfo
		nodeKVs           api.KVPairs
		nodesLastIndex    uint64
		nodesLocker       *sync.RWMutex
		OnNodesChange     WatchChangeHandle
		lastLoadNodesTime time.Time
		rpcClients        map[string]*client.RpcClient
		rpcClientLocker   *sync.RWMutex
		Scheduler         *scheduler.Scheduler
		profile           *config.Profile
	}

	WatchChangeHandle func()
)

// NewCluster new cluster and reg server
func NewCluster(profile *config.Profile) (*Cluster, error) {
	cluster := new(Cluster)
	cluster.profile = profile
	cluster.ClusterId = profile.Cluster.ClusterId
	cluster.RegistryServerUrl = profile.Cluster.RegistryServer
	cluster.LeaderKey = getLeaderKey(profile.Cluster.ClusterId)
	regClient, err := consul.NewConsulClient(profile.Cluster.RegistryServer)
	if err != nil {
		logger.Node().Debug(fmt.Sprint("Cluster init error", err.Error()))
		logger.Node().Error(err, "Cluster init error")
		return nil, err
	}
	cluster.RegistryClient = regClient
	cluster.Nodes = make(map[string]*core.NodeInfo)
	cluster.nodesLocker = new(sync.RWMutex)
	cluster.rpcClients = make(map[string]*client.RpcClient)
	cluster.rpcClientLocker = new(sync.RWMutex)

	cluster.Scheduler = scheduler.NewScheduler()
	logger.Node().Debug("Cluster init success.")
	return cluster, nil
}

// electionLeader election leader role to registry server
func (c *Cluster) ElectionLeader(leaderServer string, checkUrl string) error {
	opts := &api.LockOptions{
		Key:   c.LeaderKey,
		Value: []byte(leaderServer),
		SessionOpts: &api.SessionEntry{
			Name:     leaderServer,
			TTL:      "10s",
			Behavior: "delete",
		},
	}
	locker, err := c.RegistryClient.CreateLockerOpts(opts)
	if err != nil {
		return err
	}

	_, err = locker.Locker.Lock(nil)
	if err != nil {
		return err
	}
	return nil
}

// CreateSession create session to registry with node info
func (c *Cluster) CreateSession(nodeKey string, nodeInfo *core.NodeInfo) error {
	opts := &api.LockOptions{
		Key:   nodeKey,
		Value: []byte(nodeInfo.Json()),
		SessionOpts: &api.SessionEntry{
			Name:     nodeKey,
			TTL:      "10s",
			Behavior: "delete",
		},
	}
	locker, err := c.RegistryClient.CreateLockerOpts(opts)
	if err != nil {
		return err
	}

	_, err = locker.Locker.Lock(nil)
	if err != nil {
		return err
	}
	return nil
}

// LoadOnlineNodes load all online nodes from Registry
func (c *Cluster) LoadOnlineNodes() error {
	logTitle := "Cluster.LoadOnlineNodes "
	nodeKVs, meta, err := c.RegistryClient.ListKV(c.getNodeKeyPrefix(), nil)
	if err != nil {
		logger.Cluster().Debug(logTitle + "error: " + err.Error())
		return errors.New(logTitle + "error: " + err.Error())
	}
	c.nodeKVs = nodeKVs
	c.nodesLastIndex = meta.LastIndex
	c.refreshOnlineNodes(nodeKVs)
	c.watchOnlineNodes()
	return nil
}

// GetLeaderInfo get leader info from leader key
// must check is locked by leader session
// leader changed by watchLeaderChange
func (c *Cluster) GetLeaderInfo() (string, error) {
	if c.LeaderServer != "" {
		return c.LeaderServer, nil
	}
	kvPair, meta, err := c.RegistryClient.Get(c.LeaderKey, nil)
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

// addNodeToList add node into node list
func (c *Cluster) AddNodeInfo(nodeInfo *core.NodeInfo) {
	key := nodeInfo.EndPoint()
	c.nodesLocker.Lock()
	defer c.nodesLocker.Unlock()
	//TODO get remote worker's resource
	c.Scheduler.SetResource(key, 0, 0, 0)
	c.Nodes[key] = nodeInfo
}

func (c *Cluster) GetRpcClient(endPoint string) *client.RpcClient {
	defer c.rpcClientLocker.Unlock()
	c.rpcClientLocker.Lock()
	var rpcClient *client.RpcClient
	var isExists bool
	if rpcClient, isExists = c.rpcClients[endPoint]; !isExists {
		rpcClient = client.NewRpcClient(endPoint, c.profile.Rpc.ClientCertFile, c.profile.Rpc.ClientKeyFile)
		c.rpcClients[endPoint] = rpcClient
	}
	return rpcClient
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
		} else {
			node.IsOnline = true
		}
	}
	c.lastLoadNodesTime = time.Now()
}

// watchLeader
func (c *Cluster) WatchLeader() error {
	logTitle := "Cluster.WatchLeader "
	defer func() {
		if err := recover(); err != nil {
			errInfo := errors.New(fmt.Sprintln(err))
			logger.Cluster().Error(errInfo, logTitle+"error")
		}
	}()

	opt := &api.QueryOptions{
		WaitIndex: c.leaderLastIndex,
		WaitTime:  time.Minute * 10,
	}
	kvPair, meta, err := c.RegistryClient.Get(c.LeaderKey, opt)
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

// watchOnlineNodes
func (c *Cluster) watchOnlineNodes() error {
	logTitle := "Cluster.watchOnlineNodes "
	logger.Cluster().Debug(logTitle + "running.")
	doQuery := func() error {
		defer func() {
			if err := recover(); err != nil {
				errInfo := errors.New(fmt.Sprintln(err))
				logger.Cluster().Error(errInfo, logTitle+"error, now will retry it")
			}
		}()

		opt := &api.QueryOptions{
			WaitIndex: c.nodesLastIndex,
			WaitTime:  time.Minute * 10,
		}
		nodeKVs, meta, err := c.RegistryClient.ListKV(c.getNodeKeyPrefix(), opt)
		if err != nil {
			return err
		}
		if meta.LastIndex != c.nodesLastIndex {
			logger.Cluster().Debug("Cluster.watchNodesChange: some nodes changed.")
			c.nodeKVs = nodeKVs
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
			err := doQuery()
			if err != nil {
				logger.Cluster().DebugS(logTitle+"error, will retry after 10 seconds:", err.Error())
				time.Sleep(time.Second * 10)
			}
		}
	}()

	return nil
}

func (c *Cluster) getNodeKeyPrefix() string {
	return core.NodeKeyPrefix + c.ClusterId + "/"
}

func getLeaderKey(clusterId string) string {
	return "devfeel/rockman/" + clusterId + "/leader/locker"
}

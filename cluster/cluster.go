package cluster

import (
	"errors"
	"fmt"
	"github.com/devfeel/rockman/logger"
	"github.com/devfeel/rockman/packets"
	"github.com/devfeel/rockman/rpc/client"
	"github.com/devfeel/rockman/scheduler"
	"github.com/devfeel/rockman/util/consul"
	"github.com/hashicorp/consul/api"
	"sync"
	"time"
)

type (
	Cluster struct {
		ClusterId             string
		RegistryServerUrl     string
		RegistryClient        *consul.ConsulClient
		LeaderKey             string
		LeaderServer          string
		leaderLastIndex       uint64
		lastGetLeaderInfoTime time.Time
		OnLeaderChange        LeaderChangeHandle
		Nodes                 map[string]*packets.NodeInfo
		nodeLocker            *sync.RWMutex
		rpcClients            map[string]*client.RpcClient
		rpcClientLocker       *sync.RWMutex
		Scheduler             *scheduler.Scheduler
	}

	LeaderChangeHandle func(leader string)
)

// NewCluster new cluster and reg server
func NewCluster(clusterId string, registryServer string, leaderKey string) (*Cluster, error) {
	cluster := new(Cluster)
	cluster.ClusterId = clusterId
	cluster.RegistryServerUrl = registryServer
	cluster.LeaderKey = leaderKey
	regClient, err := consul.NewConsulClient(registryServer)
	if err != nil {
		logger.Node().Debug(fmt.Sprint("Cluster init error", err.Error()))
		logger.Node().Error(err, "Cluster init error")
		return nil, err
	}
	cluster.RegistryClient = regClient
	cluster.Nodes = make(map[string]*packets.NodeInfo)
	cluster.nodeLocker = new(sync.RWMutex)
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
func (c *Cluster) CreateSession(nodeKey string, nodeInfo *packets.NodeInfo) error {
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

// RefreshNodes refresh node state from Registry
func (c *Cluster) RefreshNodes() error {
	nodeKVs, _, err := c.RegistryClient.ListKV(packets.NodeKeyPrefix)
	if err != nil {
		return errors.New("RefreshNodes error: " + err.Error())
	}
	nodes := make(map[string]*packets.NodeInfo)
	for _, s := range nodeKVs {
		if s.Session == "" {
			continue
		}
		node := new(packets.NodeInfo)
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
	c.nodeLocker.Lock()
	defer c.nodeLocker.Unlock()
	for _, node := range c.Nodes {
		if _, exists := nodes[node.EndPoint()]; !exists {
			node.IsOnline = false
		} else {
			node.IsOnline = true
		}
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
	kvPair, meta, err := c.RegistryClient.Get(c.LeaderKey, nil)
	if err != nil {
		return "", err
	} else {
		if kvPair.Session == "" {
			return "", errors.New("lock session is nil")
		} else {
			c.LeaderServer = string(kvPair.Value)
			c.leaderLastIndex = meta.LastIndex
			c.lastGetLeaderInfoTime = time.Now()
			go c.watchLeaderChange()
			logger.Cluster().Debug("Cluster.GetLeaderInfo success [" + c.LeaderServer + "]")
			return c.LeaderServer, nil
		}
	}
}

// addNodeToList add node into node list
func (c *Cluster) AddNodeInfo(nodeInfo *packets.NodeInfo) {
	key := nodeInfo.EndPoint()
	c.nodeLocker.Lock()
	defer c.nodeLocker.Unlock()
	//TODO get remote worker's resource
	c.Scheduler.SetResource(key, 0, 0, 0)
	c.Nodes[key] = nodeInfo
}

func (c *Cluster) GetRpcClient(host, port string) *client.RpcClient {
	serverUrl := host + ":" + port
	defer c.rpcClientLocker.Unlock()
	c.rpcClientLocker.Lock()
	var rpcClient *client.RpcClient
	var isExists bool
	if rpcClient, isExists = c.rpcClients[serverUrl]; !isExists {
		rpcClient = client.NewRpcClient(serverUrl)
		c.rpcClients[serverUrl] = rpcClient
	}
	return rpcClient
}

// GetLowBalanceWorker get lower balance worker, if not match, it will try 3 times
func (c *Cluster) GetLowBalanceWorker() (*packets.NodeInfo, error) {
	resources, err := c.Scheduler.Schedule(scheduler.Balance_LowerLoad)
	if err != nil {
		return nil, err
	}

	c.nodeLocker.RLock()
	defer c.nodeLocker.RUnlock()

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

// watchLeaderChange
func (c *Cluster) watchLeaderChange() error {
	logger.Cluster().Debug("Cluster.watchLeaderChange start.")
	doQuery := func() {
		defer func() {
			if err := recover(); err != nil {
				errInfo := errors.New(fmt.Sprintln(err))
				logger.Cluster().Error(errInfo, "Cluster.watchLeaderChange error")
			}
		}()

		opt := &api.QueryOptions{
			WaitIndex: c.leaderLastIndex,
			WaitTime:  time.Minute * 10,
		}
		kvPair, meta, err := c.RegistryClient.Get(c.LeaderKey, opt)
		if err != nil {
			logger.Cluster().DebugS("Cluster.watchLeaderChange error:", err.Error())
			return
		}
		if kvPair.Session == "" {
			logger.Cluster().DebugS("Cluster.watchLeaderChange error: lock session is nil")
			return
		}
		if meta.LastIndex != c.leaderLastIndex {
			c.leaderLastIndex = meta.LastIndex
			c.LeaderServer = string(kvPair.Value)
			c.lastGetLeaderInfoTime = time.Now()
			if c.OnLeaderChange != nil {
				c.OnLeaderChange(c.LeaderServer)
			}
		}
	}

	for {
		doQuery()
	}

	return nil
}

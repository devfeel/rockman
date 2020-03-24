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
		lastGetLeaderInfoTime time.Time
		Workers               map[string]*packets.NodeInfo
		workerLocker          *sync.RWMutex
		isRegisterWorker      bool
		rpcClients            map[string]*client.RpcClient
		rpcClientLocker       *sync.RWMutex
		Scheduler             *scheduler.Scheduler
	}
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
	cluster.Workers = make(map[string]*packets.NodeInfo)
	cluster.workerLocker = new(sync.RWMutex)
	cluster.rpcClients = make(map[string]*client.RpcClient)
	cluster.rpcClientLocker = new(sync.RWMutex)

	cluster.Scheduler = scheduler.NewScheduler()
	logger.Node().Debug("Cluster init success.")
	return cluster, nil
}

// electionLeader election leader role to registry server
func (c *Cluster) ElectionLeader(leaderServer string, checkUrl string) (bool, error) {
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
		return false, err
	}

	_, err = locker.Locker.Lock(nil)
	if err != nil {
		return false, err
	}
	return true, nil
}

// RegisterNode register node info to registry server
func (c *Cluster) RegisterNode(nodeKey string, node *packets.NodeInfo) error {
	opts := &api.LockOptions{
		Key:   nodeKey,
		Value: []byte(node.EndPoint()),
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

// GetLeaderInfo get leader info from leader key
// must check is locked by leader session
// cache in memory with 1 minute
func (c *Cluster) GetLeaderInfo() (string, error) {
	if c.LeaderServer != "" && time.Now().Sub(c.lastGetLeaderInfoTime) < time.Minute {
		return c.LeaderServer, nil
	}
	kvPair, err := c.RegistryClient.Get(c.LeaderKey)
	if err != nil {
		return "", err
	} else {
		if kvPair.Session == "" {
			return "", errors.New("no leader info exists")
		} else {
			c.LeaderServer = string(kvPair.Value)
			c.lastGetLeaderInfoTime = time.Now()
			return c.LeaderServer, nil
		}
	}
}

// RegisterWorker register worker node to leader server
func (c *Cluster) RegisterWorker(worker *packets.NodeInfo) error {
	if c.isRegisterWorker {
		return nil
	}
	var leaderServer string
	var err error
GetLeader:
	for {
		// get leader info
		leaderServer, err = c.GetLeaderInfo()
		if err != nil {
			logger.Cluster().Debug("Cluster.RegisterWorker GetLeaderInfo error, will retry 10 seconds after.")
			time.Sleep(time.Second * 10)
			continue GetLeader
		} else {
			logger.Cluster().Debug("Cluster.RegisterWorker GetLeaderInfo success.")
			break
		}
	}
	rpcClient := client.NewRpcClient(leaderServer)
	err, _ = rpcClient.CallRegisterWorker(worker)
	if err == nil {
		c.isRegisterWorker = true
	}
	return err
}

// AddWorker add worker into workers
func (c *Cluster) AddWorker(worker *packets.NodeInfo) error {
	key := worker.EndPoint()
	c.workerLocker.Lock()
	defer c.workerLocker.Unlock()
	rawWorker, isExists := c.Workers[key]
	if isExists {
		logger.Cluster().Debug("Cluster replace worker node:" + fmt.Sprint(rawWorker, worker))
	} else {
		logger.Cluster().Debug("Cluster add worker node:" + fmt.Sprint(worker))
	}
	//TODO get remote worker's resource
	c.Scheduler.SetResource(key, 0, 0, 0)
	c.Workers[key] = worker
	return nil
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

	c.workerLocker.RLock()
	defer c.workerLocker.RUnlock()

	resource := resources[0]
	rawWorker, isExists := c.Workers[resource.EndPoint]
	if isExists {
		return rawWorker, nil
	}
	logger.Cluster().Debug("try get lower load worker[" + resource.EndPoint + "] failed 1 times, try get next")
	if len(resources) > 1 {
		resource := resources[1]
		rawWorker, isExists := c.Workers[resource.EndPoint]
		if isExists {
			return rawWorker, nil
		}
	}
	logger.Cluster().Debug("try get lower load worker[" + resource.EndPoint + "] failed 2 times, try get next.")
	if len(resources) > 2 {
		resource := resources[2]
		rawWorker, isExists := c.Workers[resource.EndPoint]
		if isExists {
			return rawWorker, nil
		}
	}
	logger.Cluster().Debug("try get lower load worker[" + resource.EndPoint + "] failed 3 times.")
	return nil, errors.New("no match resource with worker")
}

package cluster

import (
	"errors"
	"fmt"
	"github.com/devfeel/rockman/logger"
	"github.com/devfeel/rockman/packets"
	"github.com/devfeel/rockman/rpc/client"
	"github.com/devfeel/rockman/state"
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
		Workers               map[string]*packets.WorkerInfo
		workerLocker          *sync.RWMutex
		isRegisterWorker      bool
		rpcClients            map[string]*client.RpcClient
		rpcClientLocker       *sync.RWMutex
		state                 *state.State
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
	cluster.Workers = make(map[string]*packets.WorkerInfo)
	cluster.workerLocker = new(sync.RWMutex)
	cluster.rpcClients = make(map[string]*client.RpcClient)
	cluster.rpcClientLocker = new(sync.RWMutex)

	cluster.state = state.NewState()
	logger.Node().Debug("Cluster init success.")
	return cluster, nil
}

// electionLeader election leader role to registry server
func (c *Cluster) ElectionLeader(leaderServer string, checkUrl string) (bool, error) {
	opts := &api.LockOptions{
		Key:         c.LeaderKey,
		Value:       []byte(leaderServer),
		SessionTTL:  "10s",
		SessionName: leaderServer,
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
func (c *Cluster) RegisterWorker(worker *packets.WorkerInfo) error {
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
func (c *Cluster) AddWorker(worker *packets.WorkerInfo) error {
	key := worker.Host + "," + worker.Port
	c.workerLocker.Lock()
	defer c.workerLocker.Unlock()
	rawWorker, isExists := c.Workers[key]
	if isExists {
		logger.Cluster().Debug("Cluster replace worker node:" + fmt.Sprint(rawWorker, worker))
	} else {
		logger.Cluster().Debug("Cluster add worker node:" + fmt.Sprint(worker))
	}
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

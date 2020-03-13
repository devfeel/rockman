package cluster

import (
	"fmt"
	"github.com/devfeel/rockman/src/core/packets"
	"github.com/devfeel/rockman/src/logger"
	"github.com/devfeel/rockman/src/util/consul"
	"github.com/hashicorp/consul/api"
	"sync"
)

const (
	registryLockerKey = "master:locker"
)

type Cluster struct {
	Id         string
	Registry   *Registry
	IsMaster   bool
	Masters    map[string]*packets.NodeInfo
	Workers    map[string]*packets.NodeInfo
	nodeLocker *sync.RWMutex
}

type Registry struct {
	ServerUrl string
	RegServer *consul.ConsulClient
}

func NewCluster(clusterId string, registryUrl string) (*Cluster, error) {
	c := new(Cluster)
	c.Registry = new(Registry)
	c.Id = clusterId
	c.Registry.ServerUrl = registryUrl
	regServer, err := consul.NewConsulClient(c.Registry.ServerUrl)
	if err != nil {
		logger.Default().Debug(fmt.Sprint("Cluster Init error", err.Error()))
		return nil, err
	}

	c.Masters = make(map[string]*packets.NodeInfo)
	c.Workers = make(map[string]*packets.NodeInfo)
	c.nodeLocker = new(sync.RWMutex)

	c.Registry.RegServer = regServer
	logger.Default().Debug("Cluster Init Success!")
	return c, nil
}

// RegisterMaster register master role to registry server
func (c *Cluster) RegisterMaster(serverUrl string, checkUrl string) (bool, error) {
	opts := &api.LockOptions{
		Key:         c.getRegistryLockerKey(),
		Value:       []byte(serverUrl),
		SessionTTL:  "10s",
		SessionName: serverUrl,
	}
	locker, err := c.Registry.RegServer.CreateLockerOpts(opts)
	if err != nil {
		return false, err
	}

	_, err = locker.Locker.Lock(nil)
	if err != nil {
		return false, err
	}
	return true, nil
}

// AddNode add node into Cluster
func (c *Cluster) AddNode(node *packets.NodeInfo) error {
	key := node.Host + ":" + node.Port
	c.nodeLocker.Lock()
	defer c.nodeLocker.Unlock()
	if node.IsMaster {
		rawNode, isExists := c.Masters[key]
		if isExists {
			logger.Default().Debug("replace master node:" + fmt.Sprint(rawNode, node))
		} else {
			logger.Default().Debug("add master node:" + fmt.Sprint(node))
		}
		c.Masters[key] = node
	}

	if node.IsWorker {
		rawNode, isExists := c.Workers[key]
		if isExists {
			logger.Default().Debug("replace worker node:" + fmt.Sprint(rawNode, node))
		} else {
			logger.Default().Debug("add worker node:" + fmt.Sprint(node))
		}
		c.Workers[key] = node
	}
	return nil
}

func (c *Cluster) getRegistryLockerKey() string {
	return "devfeel/rockman:" + c.Id + ":" + registryLockerKey
}

func (reg *Registry) Register(address string, port string, checkUrl string) error {
	return nil
}

package cluster

import (
	"fmt"
	"github.com/devfeel/rockman/src/core"
	"github.com/devfeel/rockman/src/logger"
	"github.com/devfeel/rockman/src/util/consul"
	"github.com/hashicorp/consul/api"
)

const (
	registryLockerKey = "master:locker"
)

type Cluster struct {
	Id       string
	Registry *Registry
	IsMaster bool
	Workers  []*core.ServerInfo
}

type Registry struct {
	ServerUrl string
	RegServer *consul.ConsulClient
}

func NewCluster(clusterId string, serverUrl string) (*Cluster, error) {
	c := new(Cluster)
	c.Registry = new(Registry)
	c.Id = clusterId
	c.Registry.ServerUrl = serverUrl
	regServer, err := consul.NewConsulClient(c.Registry.ServerUrl)
	if err != nil {
		logger.Default().Debug(fmt.Sprint("Cluster Init error", err.Error()))
		return nil, err
	}

	c.Registry.RegServer = regServer

	fmt.Println(c.Workers)

	logger.Default().Debug("Cluster Init Success!")
	return c, nil
}

// RegisterMaster register master role to registry server
func (c *Cluster) RegisterMaster(address string, port string, checkUrl string) (bool, error) {
	opts := &api.LockOptions{
		Key:         c.getRegistryLockerKey(),
		Value:       []byte(address + "," + port),
		SessionTTL:  "10s",
		SessionName: address + "," + port,
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

func (reg *Registry) Register(address string, port string, checkUrl string) error {
	return nil
}

func (c *Cluster) getRegistryLockerKey() string {
	return "devfeel/rockman:" + c.Id + ":" + registryLockerKey
}

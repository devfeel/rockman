package consul

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	consulapi "github.com/hashicorp/consul/api"
)

type (
	ConsulClient struct {
		client *consulapi.Client
		addr   string
	}

	ServiceConfig struct {
		Name     string
		Tags     []string
		Address  string
		Port     int
		ChechUrl string
	}
)

var ErrorNotExistsKey = errors.New("not exists this key")

func NewConsulClient(addr string) (*ConsulClient, error) {
	client := &ConsulClient{}
	client.addr = addr

	config := consulapi.DefaultConfig()
	config.Address = addr
	innerClient, err := consulapi.NewClient(config)
	if err != nil {
		return nil, err
	}

	client.client = innerClient
	return client, nil
}

func (c *ConsulClient) RegisterService(service *ServiceConfig) error {
	registration := new(consulapi.AgentServiceRegistration)
	registration.ID = hashService(service)
	registration.Name = service.Name
	registration.Port = service.Port
	registration.Tags = service.Tags
	registration.Address = service.Address

	registration.Check = &consulapi.AgentServiceCheck{
		HTTP:                           service.ChechUrl,
		Timeout:                        "10s",
		Interval:                       "10s",
		DeregisterCriticalServiceAfter: "30s", //check失败后30秒删除本服务
	}
	err := c.GetClient().Agent().ServiceRegister(registration)
	return err

}

func (c *ConsulClient) SearchService(addr string, serviceName string, tag string) ([]*ServiceConfig, error) {
	services, _, err := c.GetClient().Catalog().Service(serviceName, tag, nil)
	if err != nil {
		return nil, err
	}

	var appServices []*ServiceConfig
	for _, v := range services {
		appServices = append(appServices, &ServiceConfig{
			Name:    v.ServiceName,
			Tags:    v.ServiceTags,
			Address: v.Address,
			Port:    v.ServicePort,
		})
	}
	return appServices, nil
}

func (c *ConsulClient) Get(key string, opt *consulapi.QueryOptions) (*consulapi.KVPair, *consulapi.QueryMeta, error) {
	client := c.GetClient()

	kvPair, meta, err := client.KV().Get(key, opt)
	if err != nil {
		return nil, nil, err
	}
	if kvPair == nil {
		return nil, nil, ErrorNotExistsKey
	}
	return kvPair, meta, nil
}

func (c *ConsulClient) CreateLockerOpts(opts *consulapi.LockOptions) (*Locker, error) {
	locker, err := c.GetClient().LockOpts(opts)
	if err != nil {
		return nil, err
	}
	return &Locker{Locker: locker}, nil
}

func (c *ConsulClient) CreateLocker(key string) (*Locker, error) {
	locker, err := c.GetClient().LockKey(key)
	if err != nil {
		return nil, err
	}
	return &Locker{Locker: locker}, nil
}

func (c *ConsulClient) ListSession() ([]*consulapi.SessionEntry, *consulapi.QueryMeta, error) {
	return c.GetClient().Session().List(nil)
}

func (c *ConsulClient) ListKV(prefix string, opt *consulapi.QueryOptions) (consulapi.KVPairs, *consulapi.QueryMeta, error) {
	return c.GetClient().KV().List(prefix, opt)
}

func (c *ConsulClient) GetClient() *consulapi.Client {
	return c.client
}

// hashService hash service to string
func hashService(service *ServiceConfig) string {
	data, err := json.Marshal(service)
	if err != nil {
		return ""
	}
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

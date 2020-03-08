package consul

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	consulapi "github.com/hashicorp/consul/api"
)

type (
	ConsulClient struct {
		client *consulapi.Client
		addr   string
	}

	ConsulLocker struct {
		Locker *consulapi.Lock
	}

	ServiceConfig struct {
		Name     string
		Tags     []string
		Address  string
		Port     int
		ChechUrl string
	}
)

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

func (c *ConsulClient) CreateLockerOpts(opts *consulapi.LockOptions) (*ConsulLocker, error) {
	locker, err := c.GetClient().LockOpts(opts)
	if err != nil {
		return nil, err
	}
	return &ConsulLocker{Locker: locker}, nil
}

func (c *ConsulClient) CreateLocker(addr string, key string) (*ConsulLocker, error) {
	locker, err := c.GetClient().LockKey(key)
	if err != nil {
		return nil, err
	}
	return &ConsulLocker{Locker: locker}, nil
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

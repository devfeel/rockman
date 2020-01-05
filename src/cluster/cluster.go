package cluster

import "github.com/devfeel/rockman/src/logger"

type Cluster struct {
	Registry *Registry
}

type Registry struct {
	ServerUrl string
}

func NewCluster() *Cluster {
	c := new(Cluster)
	c.Registry = new(Registry)
	logger.Default().Debug("Cluster Init Success!")
	return c
}

func (reg *Registry) Register() error {
	return nil
}

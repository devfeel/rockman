package cluster

type Cluster struct {
	Registry *Registry
}

type Registry struct {
	ServerUrl string
}

func NewCluster() *Cluster {
	c := new(Cluster)
	c.Registry = new(Registry)
	return c
}

func (reg *Registry) Register() error {
	return nil
}

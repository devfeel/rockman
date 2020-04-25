package node

import (
	"github.com/devfeel/rockman/core"
	"time"
)

func (n *Node) setExecutorChangeFlag() error {
	_, err := n.Registry.Set(getExecutorChangeFlagKey(n.ClusterId()), time.Now().Format(core.DefaultTimeLayout), nil)
	return err
}

func getExecutorChangeFlagKey(clusterId string) string {
	return core.ClusterKeyPrefix + clusterId + "/flags/executor-change"
}

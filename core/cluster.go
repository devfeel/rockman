package core

type ClusterInfo struct {
	ClusterId             string
	RegistryServerUrl     string
	LeaderKey             string
	LeaderServer          string
	NodeNum               int
	WatchLeaderRetryLimit int
	QueryResourceInterval int
}

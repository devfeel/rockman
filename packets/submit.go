package packets

type SubmitInfo struct {
	ExecutorConfig interface{}
	Worker         *WorkerInfo
	DistributeType int ``
}

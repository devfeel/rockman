package packets

import "github.com/devfeel/rockman/runtime/executor"

type SubmitInfo struct {
	Executor       executor.Executor
	Worker         *WorkerInfo
	DistributeType int ``
}

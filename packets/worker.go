package packets

type WorkerInfo struct {
	NodeID string
	Host   string
	Port   string
}

func (w *WorkerInfo) EndPoint() string {
	return w.Host + ":" + w.Port
}

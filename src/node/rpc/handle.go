package rpc

type RpcHandler struct {
	server *RpcServer
}

// Echo
func (h *RpcHandler) Echo(content string, result *JsonResult) error {
	*result = JsonResult{0, "ok", content}
	return nil
}

// RegisterNode
func (h *RpcHandler) RegisterNode(host string, port string, nodeId string, isWorker bool, result *JsonResult) error {
	if h.server.RpcHost == host && h.server.RpcPort == port {
		*result = JsonResult{-1001, "can not register node to self", nil}
		return nil
	}
	*result = JsonResult{0, "ok", nil}
	return nil
}

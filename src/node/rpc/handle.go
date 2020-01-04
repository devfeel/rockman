package rpc

type RpcHandler struct {
}

// Echo
func (h *RpcHandler) Echo(content string, result *JsonResult) error {
	*result = JsonResult{0, "ok", content}
	return nil
}

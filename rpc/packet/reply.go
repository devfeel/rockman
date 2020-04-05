package packet

import "github.com/devfeel/rockman/core"

type RpcReply struct {
	RetCode int
	RetMsg  string
	Message interface{}
}

func (r *RpcReply) IsSuccess() bool {
	return r.RetCode == core.SuccessCode
}

func CreateRpcReply(retCode int, retMsg string, message interface{}) RpcReply {
	return RpcReply{RetCode: retCode, RetMsg: retMsg, Message: message}
}

func CreateFailedReply(retCode int, retMsg string) RpcReply {
	return RpcReply{RetCode: retCode, RetMsg: retMsg}
}

func CreateSuccessRpcReply(message interface{}) RpcReply {
	return RpcReply{RetCode: core.SuccessCode, RetMsg: "", Message: message}
}

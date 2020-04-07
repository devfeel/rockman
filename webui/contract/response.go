package contract

import _const "github.com/devfeel/rockman/webui/const"

type Response struct {
	RetCode int
	RetMsg  string
	Message interface{}
}

func NewResponse(retCode int, retMsg string, message interface{}) *Response {
	return &Response{
		RetCode: retCode,
		RetMsg:  retMsg,
		Message: message,
	}
}

func SuccessResponse(message interface{}) *Response {
	return &Response{
		RetCode: _const.SuccessCode,
		RetMsg:  "",
		Message: message,
	}
}

func FailedResponse(retCode int, retMsg string) *Response {
	return &Response{
		RetCode: retCode,
		RetMsg:  retMsg,
	}
}

package controllers

const (
	SuccessCode = 0
	ErrorCode   = -9999
)

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
		RetCode: SuccessCode,
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

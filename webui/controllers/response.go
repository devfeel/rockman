package controllers

import (
	"github.com/devfeel/rockman/webui/contract"
)

func NewResponse(retCode int, retMsg string, message interface{}) *contract.Response {
	return contract.NewResponse(retCode, retMsg, message)
}

func SuccessResponse(message interface{}) *contract.Response {
	return contract.SuccessResponse(message)
}

func FailedResponse(retCode int, retMsg string) *contract.Response {
	return contract.FailedResponse(retCode, retMsg)
}

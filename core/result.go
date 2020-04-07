package core

import "strconv"

const (
	SuccessCode = 0
	ErrorCode   = -9999
)

type Result struct {
	RetCode int
	RetMsg  string
	Error   error
}

func NewResult(retCode int, retMsg string, err error) *Result {
	return &Result{RetCode: retCode, RetMsg: retMsg, Error: err}
}

func ErrorResult(err error) *Result {
	return &Result{RetCode: ErrorCode, RetMsg: err.Error(), Error: err}
}

func SuccessResult() *Result {
	return &Result{RetCode: SuccessCode, RetMsg: "", Error: nil}
}

func FailedResult(retCode int, retMsg string) *Result {
	return &Result{RetCode: retCode, RetMsg: retMsg, Error: nil}
}

func (r *Result) IsSuccess() bool {
	return r.Error == nil && r.RetCode == SuccessCode
}

func (r *Result) Message() string {
	if r.Error != nil {
		return r.Error.Error()
	}
	return strconv.Itoa(r.RetCode) + "," + r.RetMsg
}

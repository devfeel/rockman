package controllers

import (
	"github.com/pkg/errors"
)

func IsNilString(val string, errCode int, errMsg string) (*Response, error) {
	if val != "" {
		return SuccessResponse(nil), nil
	} else {
		return NewResponse(errCode, errMsg, nil), errors.New("val is nil")
	}
}

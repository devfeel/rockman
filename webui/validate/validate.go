package validate

import (
	"errors"
	"github.com/devfeel/rockman/webui/contract"
)

func IsNilString(val string, errCode int, errMsg string) (*contract.Response, error) {
	if val != "" {
		return contract.SuccessResponse(nil), nil
	} else {
		return contract.NewResponse(errCode, errMsg, nil), errors.New("val is nil")
	}
}

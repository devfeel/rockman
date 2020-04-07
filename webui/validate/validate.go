package validate

import (
	"github.com/devfeel/rockman/webui/contract"
	"github.com/pkg/errors"
)

func IsNilString(val string, errCode int, errMsg string) (*contract.Response, error) {
	if val != "" {
		return contract.SuccessResponse(nil), nil
	} else {
		return contract.NewResponse(errCode, errMsg, nil), errors.New("val is nil")
	}
}

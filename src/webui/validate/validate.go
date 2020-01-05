package validate

import (
	"github.com/devfeel/rockman/src/webui/contract"
	"github.com/pkg/errors"
)

func IsNilString(val string, errCode int, errMsg string) (*contract.ResponseInfo, error) {
	if val != "" {
		return contract.NewResonseInfo(), nil
	} else {
		return contract.CreateResponse(errCode, errMsg, nil), errors.New("val is nil")
	}
}

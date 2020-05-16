package sysx

import (
	"fmt"
	"testing"
)

func TestGetCpuTimeState(t *testing.T) {
	fmt.Println(GetCpuTimeState())
}

func TestGetCpuUsedPercent(t *testing.T) {
	fmt.Println(GetCpuUsedPercent())
}

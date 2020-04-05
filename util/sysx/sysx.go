package sysx

import (
	"errors"
	"github.com/shirou/gopsutil/cpu"
)

func GetCpuTimeState() (*cpu.TimesStat, error) {
	states, err := cpu.Times(false)
	if err != nil {
		return nil, err
	}
	for _, state := range states {
		if state.CPU == "cpu-total" {
			return &state, nil
		}
	}
	return nil, errors.New("not exists cpu info")
}

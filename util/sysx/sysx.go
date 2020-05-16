package sysx

import (
	"errors"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

func GetCpuUsedPercent() int {
	state, err := GetCpuTimeState()
	if err != nil {
		return 100
	}
	return (int)(((state.User + state.System) / (state.User + state.System + state.Idle)) * 100)
}

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

func GetMemoryUsedPercent() int {
	v, err := mem.VirtualMemory()
	if err != nil {
		return 100
	}
	return (int)(v.UsedPercent)
}

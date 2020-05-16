package handler

import (
	"fmt"
	"github.com/shirou/gopsutil/mem"
	"testing"
)

func Test_GetMemory(t *testing.T) {
	v, _ := mem.VirtualMemory()
	fmt.Printf("UsedPercent:%f\n", v.UsedPercent)
	fmt.Println((int)(v.UsedPercent))
}

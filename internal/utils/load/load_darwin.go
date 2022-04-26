//go:build darwin
// +build darwin

package load

import (
	"golang.org/x/sys/unix"
	"unsafe"
)

func Avg() (float64, error) {
	type loadavg struct {
		load  [3]uint32
		scale int
	}

	b, err := unix.SysctlRaw("vm.loadavg")
	if err != nil {
		return 0, err
	}
	load := *(*loadavg)(unsafe.Pointer(&b[0]))
	scale := float64(load.scale)

	return float64(load.load[0]) / scale, nil
}

//go:build linux
// +build linux

package cpu

import (
	"fmt"
	"github.com/vlyagusha/system_stats_daemon/internal/app"
	"runtime"
)

func Get() (*app.CpuStats, error) {
	return nil, fmt.Errorf("cpu statistics not implemented for: %s", runtime.GOOS)
}

//go:build linux
// +build linux

package disk

import (
	"fmt"
	"github.com/vlyagusha/system_stats_daemon/internal/app"
	"runtime"
)

func Get() (*app.DiskStats, error) {
	return nil, fmt.Errorf("disk statistics not implemented for: %s", runtime.GOOS)
}

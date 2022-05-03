//go:build linux
// +build linux

package cpu

import (
	"github.com/vlyagusha/system_stats_daemon/internal/app"
	"github.com/vlyagusha/system_stats_daemon/internal/utils"
)

func Get() (*app.CPUStats, error) {
	return nil, utils.ErrNotImplemented
}

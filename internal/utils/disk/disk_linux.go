//go:build linux
// +build linux

package disk

import (
	"github.com/vlyagusha/system_stats_daemon/internal/app"
	"github.com/vlyagusha/system_stats_daemon/internal/utils"
)

func Get() (*app.DiskStats, error) {
	return nil, utils.ErrNotImplemented
}

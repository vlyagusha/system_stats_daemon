//go:build linux
// +build linux

package disk

import (
	"github.com/stretchr/testify/require"
	"github.com/vlyagusha/system_stats_daemon/internal/utils"
	"testing"
)

func TestGetStat(t *testing.T) {
	t.Run("test not implemented get stats", func(t *testing.T) {
		diskStat, err := Get()

		require.Nil(t, diskStat)
		require.ErrorIs(t, err, utils.ErrNotImplemented)
	})
}

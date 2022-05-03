package pipeline

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vlyagusha/system_stats_daemon/internal/config"
)

func TestGenerator(t *testing.T) {
	t.Run("test empty config", func(t *testing.T) {
		emptyStatConfig := config.StatsConfig{
			LoadAvg: false,
			CPU:     false,
			Disk:    false,
			NetTop:  false,
			NetStat: false,
		}

		result := GetStages(emptyStatConfig)
		require.Empty(t, result)
	})

	t.Run("test not implemented config", func(t *testing.T) {
		notImplementedStatConfig := config.StatsConfig{
			LoadAvg: false,
			CPU:     false,
			Disk:    false,
			NetTop:  true,
			NetStat: true,
		}

		result := GetStages(notImplementedStatConfig)
		require.Empty(t, result)
	})

	t.Run("test implemented config", func(t *testing.T) {
		implementedStatConfig := config.StatsConfig{
			LoadAvg: true,
			CPU:     true,
			Disk:    true,
			NetTop:  false,
			NetStat: false,
		}

		result := GetStages(implementedStatConfig)
		require.Len(t, result, 3)
	})
}

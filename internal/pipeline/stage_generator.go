package pipeline

import (
	"github.com/vlyagusha/system_stats_daemon/internal/app"
	"github.com/vlyagusha/system_stats_daemon/internal/utils/cpu"
	"github.com/vlyagusha/system_stats_daemon/internal/utils/disk"
	"github.com/vlyagusha/system_stats_daemon/internal/utils/load"
	"log"
)

func GetStages() []Stage {
	return []Stage{
		stageGenerator(
			"Load Average Stage",
			func(stat app.SystemStats) app.SystemStats {
				loadAvg, err := load.Get()
				if err != nil {
					log.Fatalf("Failed to get load average: %s", err)
				}
				stat.Load = *loadAvg

				return stat
			}),
		stageGenerator(
			"CPU stage",
			func(stat app.SystemStats) app.SystemStats {
				cpuStat, err := cpu.Get()
				if err != nil {
					log.Fatalf("Failed to get CPU usage: %s", err)
				}
				stat.CPU = *cpuStat

				return stat
			}),
		stageGenerator(
			"Disk stage",
			func(stat app.SystemStats) app.SystemStats {
				diskStats, err := disk.Get()
				if err != nil {
					log.Fatalf("Failed to get disk stats: %s", err)
				}
				stat.Disk = *diskStats

				return stat
			}),
	}
}

func stageGenerator(_ string, f func(s app.SystemStats) app.SystemStats) Stage {
	return func(in In) Out {
		out := make(Bi)
		go func() {
			defer close(out)
			for v := range in {
				out <- f(v)
			}
		}()
		return out
	}
}

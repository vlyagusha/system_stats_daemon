package pipeline

import (
	"github.com/vlyagusha/system_stats_daemon/internal/app"
	"github.com/vlyagusha/system_stats_daemon/internal/utils/load"
	"log"
)

func GetStages() []Stage {
	return []Stage{
		stageGenerator(
			"Load Average Stage",
			func(stat app.SystemStats) app.SystemStats {
				l, err := load.Avg()
				if err != nil {
					log.Fatalf("Failed to get load average: %s", err)
				}

				stat.Main.Load = l
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

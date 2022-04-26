package pipeline

import "github.com/vlyagusha/system_stats_daemon/internal/app"

type (
	In  = <-chan app.SystemStats
	Out = In
	Bi  = chan app.SystemStats
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := in
	for _, stage := range stages {
		out = doneStage(done, out)
		out = stage(out)
	}

	return out
}

func doneStage(done In, in In) Out {
	out := make(Bi)

	go func() {
		defer close(out)
		for {
			select {
			case <-done:
				return
			case v, ok := <-in:
				if !ok {
					return
				}
				out <- v
			}
		}
	}()

	return out
}

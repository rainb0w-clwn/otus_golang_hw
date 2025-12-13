package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func drain(ch In) {
	for range ch {} //nolint:all
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	for _, stage := range stages {
		stageOut := stage(in)
		out := make(Bi)
		go func(_in In, _out Bi) {
			defer close(_out)
			for {
				select {
				case <-done:
					go drain(_in)
					return
				case v, ok := <-_in:
					if !ok {
						return
					}
					select {
					case <-done:
						go drain(_in)
						return
					case _out <- v:
					}
				}
			}
		}(stageOut, out)
		in = out
	}
	return in
}

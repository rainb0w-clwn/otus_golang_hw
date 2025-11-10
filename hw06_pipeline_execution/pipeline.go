package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	for i := 0; i < len(stages); i++ {
		// Свой выходной канал с возможностью закрытия
		out := make(Bi)
		go func(_in In, _out Bi) {
			defer func() {
				close(_out)
				// Нам необходимо "докрутить" канал входных данных, так как только это откроет wg
				for range _in {} //nolint:all
			}()
			for {
				select {
				case <-done:
					return
				case v, ok := <-_in:
					if !ok {
						return
					}
					_out <- v
				}
			}
		}(in, out)
		in = stages[i](out)
	}
	return in
}

package hw06pipelineexecution

import "fmt"

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	terminated := make(Bi)
	go func() {
		defer fmt.Println("Work exited.")
		defer close(terminated)
		i := 0
		for i < len(stages) {
			select {
			case terminated <- stages[i](in):
				fmt.Println(i)
				in = terminated
				i++
			case <-done:
				return
			default:
			}
		}
	}()
	return terminated
}

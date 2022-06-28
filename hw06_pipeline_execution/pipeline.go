package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	// output channel
	outputChan := make(Bi)
	// channel to load data into the first stage
	inputFirstStage := make(Bi)

	// goroutine to load data from in into inputFirstStage with the ability
	// to interrupt on an additional signal (done)
	go func() {
		defer close(inputFirstStage)
		for element := range in {
			select {
			case <-done:
				return
			case inputFirstStage <- element:
			}
		}
	}()

	var inputNextStage In = inputFirstStage
	var outputLastStage Out

	// loop to run stages
	for _, stage := range stages {
		outputLastStage = stage(inputNextStage)
		inputNextStage = outputLastStage
	}

	// goroutine to load the result into outputChan with the ability
	// to interrupt on an additional signal (done)
	// (with priority case <-done)
	go func() {
		defer close(outputChan)
		for {
			select {
			case <-done:
				return
			default:
				select {
				case s, ok := <-outputLastStage:
					if !ok {
						return
					}
					outputChan <- s
				case <-done:
					return
				}
			}
		}
	}()

	return outputChan
}

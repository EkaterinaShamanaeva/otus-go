package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	outChan := make(Bi)

	inFirstStage := make(Bi)
	var inNextStage In = inFirstStage
	var outLastStage Out

	go func() {
		for s := range in {
			inFirstStage <- s
		}
		close(inFirstStage)
	}()

	for _, stage := range stages {
		outLastStage = stage(inNextStage)
		inNextStage = outLastStage
	}

	go func() {
		//for s := range outLastStage {
		//	outChan <- s
		//}
		//close(outChan)
		//array := make([]interface{}, 0)
		for s := range outLastStage {
			select {
			case <-done:
				close(outChan)
				//for {
				//	_, ok := <-outChan
				//	if !ok {
				//		break
				//	}
				//}
				return
			default:
				select {
				case outChan <- s:
				case <-done:
					return
				}
			}
			//for x := range outChan {
			//	array = append(array, x)
			//}
			//return

		}
		close(outChan)

	}()

	return outChan
}

package hw06pipelineexecution

import (
	"fmt"
)

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	outChan := make(Bi)
	//defer close(outChan)
	//firstStage := make(Bi)
	var value interface{}
	fmt.Println("flag: main before loop")

	go func() {
		fmt.Println("flag: go 1 - first stage")
		//for s := range in {
		outChan <- <-in //s
		fmt.Println("flag: go 1 - read input OK")
		//}
	}()
	fmt.Println("flag: before loop")
	x := stages[0](outChan)
	fmt.Println("flag: 0 stage")
	for i := 1; i < len(stages); i++ {
		fmt.Println("flag: main - loop ", i)
		if i == len(stages)-1 {
			value = <-stages[i](x)
			fmt.Printf("flag: main - stage res: %v,  %T \n", value, value)
			outChan <- value
		} else {
			x = stages[i](x)
		}
	}

	go func() {
		fmt.Println("flag go2: start go")
		outChan <- value
		fmt.Println("flag go2: gor with outchan")
		//select {
		//case <-done:
		//	return
		//case outChan <- value:
		//}

	}()

	fmt.Println("flag: main end")
	//fmt.Println(<-outChan) //

	return outChan
}

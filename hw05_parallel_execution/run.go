package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	// counter of errors
	errCounter := 0

	wg := sync.WaitGroup{}

	tasksChan := make(chan error, len(tasks))
	defer close(tasksChan)

	// if number of goroutines > number of tasks
	if len(tasks) < n {
		n = len(tasks)
	}

	for i := 0; i < len(tasks); i += n {
		for j := 0; j < n; j++ {
			wg.Add(1)
			go func(idx int) {
				tasksChan <- tasks[idx]()
				wg.Done()
			}(i + j)
		}
		for k := 0; k < n; k++ {
			err := <-tasksChan
			if err != nil {
				errCounter++
			}
		}
		if m > 0 && errCounter >= m {
			return ErrErrorsLimitExceeded
		}
	}

	wg.Wait()

	return nil
}

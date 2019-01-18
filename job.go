/* Package for CPU bound jobs. */
package job

import (
	"context"
	"sync"
	"sync/atomic"
)

// Worker may run in parallel with other workers.
type Worker interface {
	Run(ctx context.Context) error
}

// Execute with maximum concurrent workers set via max.
// Return last error.
func Execute(ctx context.Context, workers []Worker, max int) error {
	var wg sync.WaitGroup
	var errLast atomic.Value
	sem := make(chan int, max)

	for _, w := range workers {
		wg.Add(1)
		sem <- 1

		// One of the workers error, do not start more Workers.
		if errLast.Load() != nil {
			break
		}

		go func(wrk Worker) {
			defer wg.Done()

			if err := wrk.Run(ctx); err != nil {
				errLast.Store(err)
			}
			<-sem
		}(w)
	}

	wg.Wait()
	if err := errLast.Load(); err != nil {
		return err.(error)
	}
	return nil
}

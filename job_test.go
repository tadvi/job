package job

import (
	"context"
	"errors"
	"testing"
	"time"
)

// execute as go test -v -race

type work struct {
	n         int
	withError bool
}

func (w *work) Run(ctx context.Context) error {
	for i := 0; i < w.n; i++ {
		// Do some work.
		time.Sleep(time.Millisecond)
		if w.withError && i > 7 {
			return errors.New("stop")
		}
	}
	return nil
}

func TestJobBasic(t *testing.T) {
	workers := make([]Worker, 10)
	for k := range workers {
		workers[k] = &work{k, false}
	}
	err := Execute(context.TODO(), workers, 4)
	if err != nil {
		t.Errorf("want no error, got error")
	}
}

func TestJobWithError(t *testing.T) {
	workers := make([]Worker, 10)
	for k := range workers {
		workers[k] = &work{k, true}
	}
	err := Execute(context.TODO(), workers, 4)
	if err == nil {
		t.Errorf("want error, got no error")
	}
}

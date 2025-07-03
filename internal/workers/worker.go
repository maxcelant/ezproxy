package workers

import (
	"context"
	"fmt"
	"time"
)

type worker struct {
	ctx context.Context
}

func NewWorker(ctx context.Context) *worker {
	return &worker{ctx}
}

func (w *worker) start() error {
	fmt.Println("starting worker....")
	for {
		select {
		case <-w.ctx.Done():
			fmt.Println("shutting down")
			return nil
		default:
			fmt.Println("running...")
			time.Sleep(5 * time.Second)
		}
	}
}

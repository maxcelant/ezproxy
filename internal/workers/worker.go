package workers

import (
	"context"
	"fmt"
	"time"
)

type Worker struct {
	ctx context.Context
}

func NewWorker(ctx context.Context) *Worker {
	return &Worker{ctx}
}

func (w *Worker) start() error {
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

package workers

import (
	"context"
	"fmt"
	"runtime"
	"sync"

	"github.com/maxcelant/ezproxy/internal/dispatch"
)

// Handles the lifecycle of the worker threads
type WorkerPool struct {
	workers         []*worker
	startCh         chan *worker
	notifyStartedCh chan struct{}
	errCh           chan error
	ctx             context.Context
	cancel          context.CancelFunc
	wg              sync.WaitGroup
}

func NewWorkerPool() *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())
	return &WorkerPool{
		workers:         make([]*worker, runtime.NumCPU()),
		startCh:         make(chan *worker, runtime.NumCPU()),
		notifyStartedCh: make(chan struct{}),
		errCh:           make(chan error),
		ctx:             ctx,
		cancel:          cancel,
	}

}

func (p *WorkerPool) Start() (err error) {
	defer func() {
		close(p.errCh)
		close(p.notifyStartedCh)
		close(p.startCh)
	}()
	workerCount := runtime.NumCPU()

	fmt.Printf("starting %d worker threads\n", workerCount)
	go p.reconcile()

	for range workerCount {
		p.startCh <- NewWorker(p.ctx)
	}

	for workerCount > 0 {
		select {
		case <-p.notifyStartedCh:
			workerCount--
			if workerCount == 0 {
				return nil
			}
		case err = <-p.errCh:
			return err
		}
	}
	fmt.Println("all worker threads started")
	return nil
}

func (p *WorkerPool) ForwardRequestFunc() func(dispatch.DispatchContext) {
	// Used to round robin requests to all the workers
	i := 0
	return func(ctx dispatch.DispatchContext) {
		fmt.Printf("sending request to worker %d\n", i)
		p.workers[i].handle(ctx)
		i = (i + 1) % len(p.workers)
	}
}

func (p *WorkerPool) reconcile() {
	for w := range p.startCh {
		go func(w *worker) {
			defer p.wg.Done()
			p.wg.Add(1)
			go func() {
				p.notifyStartedCh <- struct{}{}
			}()
			if err := w.start(); err != nil {
				fmt.Println("error occurred: ", err)
				p.errCh <- err
				return
			}
		}(w)
	}
}

func (p *WorkerPool) Stop() {
	fmt.Println("gracefully shutting down workers...")
	p.cancel()
	p.wg.Wait()
}

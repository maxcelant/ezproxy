package workers

import (
	"context"
	"fmt"
	"runtime"
	"sync"
)

// Handles the lifecycle of the worker threads
type WorkerPool struct {
	workers         []*Worker
	startCh         chan *Worker
	notifyStartedCh chan struct{}
	errCh           chan error
	ctx             context.Context
	cancel          context.CancelFunc
	wg              sync.WaitGroup
}

func NewWorkerPool() *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())
	return &WorkerPool{
		workers:         make([]*Worker, runtime.NumCPU()),
		startCh:         make(chan *Worker, runtime.NumCPU()),
		notifyStartedCh: make(chan struct{}),
		errCh:           make(chan error),
		ctx:             ctx,
		cancel:          cancel,
	}

}

func (p *WorkerPool) Start() (err error) {
	workerCount := runtime.NumCPU()

	fmt.Printf("starting %d worker threads", workerCount)
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

func (p *WorkerPool) reconcile() {
	for w := range p.startCh {
		go func(w *Worker) {
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

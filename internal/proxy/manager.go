package proxy

import (
	"context"
	"fmt"
	"net/url"
	"sync"
)

type HTTPProxy struct {
	listenerGroup ListenerGroup
	endpoints     []*url.URL
	ctx           context.Context
	cancel        context.CancelFunc
	wg            sync.WaitGroup
}

func NewProxyFromScratch() *HTTPProxy {
	return &HTTPProxy{
		listenerGroup: ListenerGroup{},
	}
}

func (p *HTTPProxy) AddListener(URL string) {
	if err := p.listenerGroup.Add(URL); err != nil {
		fmt.Println("error occured while adding listener", err)
	}
}

func (p *HTTPProxy) AddEndpoint(URL string) {
	e, err := url.Parse(URL)
	if err != nil {
		fmt.Println("bad downstream URL:", err)
		return
	}
	p.endpoints = append(p.endpoints, e)
}

func (p *HTTPProxy) Start() {
	p.ctx, p.cancel = context.WithCancel(context.Background())
	p.listenerGroup.Start()
	// TODO: Stop when all listeners are started
	fmt.Println("proxy has started")
}

// Gracefully handle shutdown when sigterm signal is triggered
func (p *HTTPProxy) Stop() {
	fmt.Println("gracefully shutting down proxy...")

	p.listenerGroup.Stop()

	// This should block until all goroutines are cleaned up
	p.wg.Wait()
	fmt.Println("proxy shutdown complete.")
}

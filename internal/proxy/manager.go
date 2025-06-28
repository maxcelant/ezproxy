package proxy

import (
	"fmt"
	"net/url"
	"sync"
)

type HTTPProxy struct {
	lg        listenerGroup
	endpoints []*url.URL
	wg        sync.WaitGroup
}

func NewProxyFromScratch() *HTTPProxy {
	return &HTTPProxy{
		lg: listenerGroup{},
	}
}

func (p *HTTPProxy) AddListener(URL string) {
	if err := p.lg.add(URL); err != nil {
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
	// Start the listener group
	p.lg.start()
	// TODO: Return when all listeners are started
	fmt.Println("proxy has started")
}

// Gracefully handle shutdown when sigterm signal is triggered
func (p *HTTPProxy) Stop() {
	fmt.Println("gracefully shutting down proxy...")

	// Will block until all listeners are cleaned up
	p.lg.stop()

	fmt.Println("proxy shutdown complete.")
}

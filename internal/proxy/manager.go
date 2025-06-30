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
		lg: listenerGroup{
			startCh:         make(chan *httpListener),
			errCh:           make(chan error),
			notifyStartedCh: make(chan struct{}),
			started:         false,
		},
	}
}

func (p *HTTPProxy) AddListener(URL string) error {
	if err := p.lg.add(URL); err != nil {
		return fmt.Errorf("error occured while adding listener: %w", err)
	}
	return nil
}

func (p *HTTPProxy) AddEndpoint(URL string) error {
	e, err := url.Parse(URL)
	if err != nil {
		return fmt.Errorf("bad downstream URL: %w", err)
	}
	p.endpoints = append(p.endpoints, e)
	return nil
}

func (p *HTTPProxy) Start() (err error) {
	fmt.Println("starting ezproxy")
	// Start the listener group
	err = p.lg.start()
	// TODO: Return when all listeners are started

	return err
}

// Gracefully handle shutdown when sigterm signal is triggered
func (p *HTTPProxy) Stop() {
	fmt.Println("gracefully shutting down ezproxy")

	// Will block until all listeners are cleaned up
	p.lg.stop()

	fmt.Println("proxy shutdown complete")
}

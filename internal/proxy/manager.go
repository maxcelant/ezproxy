package proxy

import (
	"fmt"
	"github.com/maxcelant/ezproxy/internal/listener"
	"net/url"
	"sync"
)

type HTTPProxy struct {
	lg        *listener.ListenerGroup
	endpoints []*url.URL
	wg        sync.WaitGroup
	log       Logger
}

func NewProxyFromScratch(log Logger) *HTTPProxy {
	return &HTTPProxy{
		lg:  listener.NewListenerGroup(),
		log: log,
	}
}

func (p *HTTPProxy) AddListener(URL string) error {
	if err := p.lg.Add(URL); err != nil {
		return fmt.Errorf("error occured while adding listener: %w", err)
	}
	p.log.Info("starting listener", "url", URL)
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
	p.log.Info("starting ezproxy")
	// Start the listener group
	err = p.lg.Start()
	return err
}

// Gracefully handle shutdown when sigterm signal is triggered
func (p *HTTPProxy) Stop() {
	p.log.Info("gracefully shutting down ezproxy")

	// Will block until all listeners are cleaned up
	p.lg.Stop()

	p.log.Info("proxy shutdown complete")
}

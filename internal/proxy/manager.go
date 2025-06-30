package proxy

import (
	"fmt"
	"log/slog"
	"net/url"
	"sync"
)

type HTTPProxy struct {
	lg        listenerGroup
	endpoints []*url.URL
	wg        sync.WaitGroup
	log       *slog.Logger
}

func NewProxyFromScratch(log *slog.Logger) *HTTPProxy {
	return &HTTPProxy{
		lg: listenerGroup{
			startCh:         make(chan *httpListener),
			errCh:           make(chan error),
			notifyStartedCh: make(chan struct{}),
			started:         false,
		},
		log: log,
	}
}

func (p *HTTPProxy) AddListener(URL string) error {
	if err := p.lg.add(URL); err != nil {
		return fmt.Errorf("error occured while adding listener: %w", err)
	}
	p.log.Info("starting listener ", "url", URL)
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
	err = p.lg.start()
	return err
}

// Gracefully handle shutdown when sigterm signal is triggered
func (p *HTTPProxy) Stop() {
	p.log.Info("gracefully shutting down ezproxy")

	// Will block until all listeners are cleaned up
	p.lg.stop()

	p.log.Info("proxy shutdown complete")
}

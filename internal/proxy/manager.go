package proxy

import (
	"sync"

	"github.com/maxcelant/ezproxy/internal/chain"
)

type HTTPProxy struct {
	chains []*chain.Chain
	wg     sync.WaitGroup
	log    Logger
}

type proxyOpts func(*HTTPProxy)

func WithChain(c *chain.Chain) proxyOpts {
	return func(h *HTTPProxy) {
		h.chains = append(h.chains, c)
	}
}

func NewProxyFromScratch(log Logger, opts ...proxyOpts) *HTTPProxy {
	p := &HTTPProxy{
		log: log,
	}

	for _, opt := range opts {
		opt(p)
	}

	return p
}

func (p *HTTPProxy) Start() error {
	p.log.Info("starting ezproxy")

	for _, c := range p.chains {
		if err := c.Start(); err != nil {
			return err
		}
	}
	return nil
}

// Gracefully handle shutdown when sigterm signal is triggered
func (p *HTTPProxy) Stop() {
	p.log.Info("gracefully shutting down ezproxy")

	for _, c := range p.chains {
		c.Stop()
	}

	p.log.Info("proxy shutdown complete")
}

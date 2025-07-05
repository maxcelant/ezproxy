package proxy

import (
	"fmt"
	"log/slog"

	"github.com/maxcelant/ezproxy/internal/chain"
	"github.com/maxcelant/ezproxy/internal/dispatch"
	"github.com/maxcelant/ezproxy/internal/workers"
)

type HTTPProxy struct {
	chains     []*chain.Chain
	workerPool *workers.WorkerPool
	log        *slog.Logger
}

type proxyOpts func(*HTTPProxy)

func WithChain(c *chain.Chain) proxyOpts {
	return func(h *HTTPProxy) {
		c.InheritLoggerFromManager(h.log)
		h.chains = append(h.chains, c)
	}
}

func NewProxyFromScratch(log *slog.Logger, opts ...proxyOpts) *HTTPProxy {
	p := &HTTPProxy{
		log:        log,
		workerPool: workers.NewWorkerPool(),
	}

	for _, opt := range opts {
		opt(p)
	}

	return p
}

func (p *HTTPProxy) Start() error {
	p.log.Info("starting ezproxy")

	if err := p.workerPool.Start(); err != nil {
		return fmt.Errorf("error starting worker pool: %w", err)
	}

	for _, c := range p.chains {
		d := dispatch.NewDispatcher()
		d.Mount(p.workerPool.ForwardRequestFunc())
		if err := c.Start(d); err != nil {
			return err
		}
	}

	return nil
}

// Gracefully handle shutdown of listeners and workers threads
func (p *HTTPProxy) Stop() {
	p.log.Info("gracefully shutting down ezproxy")

	for _, c := range p.chains {
		c.Stop()
	}

	p.workerPool.Stop()

	p.log.Info("proxy shutdown complete")
}

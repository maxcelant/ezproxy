package proxy

import (
	"fmt"
	"log/slog"

	"github.com/maxcelant/ezproxy/internal/dispatch"
	"github.com/maxcelant/ezproxy/internal/workers"
	"github.com/maxcelant/ezproxy/pkg/chain"
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

// Start initializes the proxy and starts all listeners and worker threads
// also attaches the a dispatcher to each chains
func (p *HTTPProxy) Start() error {
	p.log.Info("starting ezproxy")

	if err := p.workerPool.Start(); err != nil {
		return fmt.Errorf("error starting worker pool: %w", err)
	}

	for _, c := range p.chains {
		// IMPORTANT: We mount the dispatcher func from the worker pool here
		// This is how requests will be allowed to flow from the listener group
		// to the worker pool
		d := dispatch.NewDispatcher().
			Mount(p.workerPool.ForwardRequestFunc())
		if err := c.Start(d); err != nil {
			return err
		}
	}

	return nil
}

// Stop gracefully handle shutdown of listeners and workers threads
// We don't return until both have successfully finished
func (p *HTTPProxy) Stop() {
	p.log.Info("gracefully shutting down ezproxy")

	for _, c := range p.chains {
		c.Stop()
	}

	p.workerPool.Stop()

	p.log.Info("proxy shutdown complete")
}

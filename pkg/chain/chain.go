package chain

import (
	"fmt"
	"log/slog"

	"github.com/maxcelant/ezproxy/internal/dispatch"
	"github.com/maxcelant/ezproxy/internal/listener"
)

type Chain struct {
	lg         *listener.ListenerGroup
	upstreams  []string
	dispatcher *dispatch.Dispatcher
	log        *slog.Logger
}

type chainOpts func(*Chain) error

func WithListener(addr string) chainOpts {
	return func(c *Chain) error {
		if err := c.lg.Add(addr); err != nil {
			return fmt.Errorf("error occured while adding listener: %w", err)
		}
		return nil
	}
}

func WithUpstream(addr string) chainOpts {
	// TODO: Add more URL checking here
	return func(c *Chain) error {
		c.upstreams = append(c.upstreams, addr)
		return nil
	}
}

func NewChain(opts ...chainOpts) *Chain {
	c := &Chain{
		lg:        listener.NewListenerGroup(),
		upstreams: make([]string, 0),
	}

	var errors []error
	for _, opt := range opts {
		err := opt(c)
		if err != nil {
			errors = append(errors, err)
		}
	}

	for _, err := range errors {
		fmt.Println(err)
	}
	return c
}

func (c *Chain) Start(d *dispatch.Dispatcher) error {
	d.AddUpstreams(c.upstreams)
	err := c.lg.Start(d)
	return err
}

func (c *Chain) Stop() {
	c.lg.Stop()
}

func (c *Chain) InheritLoggerFromManager(logger *slog.Logger) {
	c.log = logger
}

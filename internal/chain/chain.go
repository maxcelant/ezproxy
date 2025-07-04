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

func WithListener(host string, port int) chainOpts {
	return func(c *Chain) error {
		return c.addListener(host, port)
	}
}

func WithDownstream(host string, port int) chainOpts {
	// TODO: Add more URL checking here
	return func(c *Chain) error {
		c.upstreams = append(c.upstreams, fmt.Sprintf("%s:%d", host, port))
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

func (c *Chain) addListener(host string, port int) error {
	URL := fmt.Sprintf("%s:%d", host, port)
	if err := c.lg.Add(URL); err != nil {
		return fmt.Errorf("error occured while adding listener: %w", err)
	}
	c.log.Debug("adding listener", "listener", URL)
	return nil
}

func (c *Chain) Start(d *dispatch.Dispatcher) error {
	d.AddUpstreams(c.upstreams)
	c.log.Debug("(chain) starting listener group")
	err := c.lg.Start(d)
	return err
}

func (c *Chain) Stop() {
	c.log.Debug("(chain) stopping listener group")
	c.lg.Stop()
}

func (c *Chain) InheritLoggerFromManager(logger *slog.Logger) {
	c.log = logger
}

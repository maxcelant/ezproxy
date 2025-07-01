package chain

import (
	"fmt"

	"github.com/maxcelant/ezproxy/internal/listener"
)

type Chain struct {
	lg listener.ListenerGroup
}

type chainOpts func(*Chain) error

func WithListener(host string, port int) chainOpts {
	return func(c *Chain) error {
		return c.addListener(host, port)
	}
}

func NewChain(opts ...chainOpts) *Chain {
	c := &Chain{}
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
	return nil
}

func (c *Chain) Start() error {
	err := c.lg.Start()
	return err
}

func (c *Chain) Stop() {
	c.lg.Stop()
}

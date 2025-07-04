package listener

import (
	"errors"
	"fmt"
	"net"

	"github.com/maxcelant/ezproxy/internal/dispatch"
)

// Wrapper over native Listener to add functionality
type httpListener struct {
	net.Listener
	*dispatch.Dispatcher
}

func NewListener(l net.Listener, d *dispatch.Dispatcher) *httpListener {
	return &httpListener{l, d}
}

func (l *httpListener) start() error {
	for {
		conn, err := l.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return err
			}
			return fmt.Errorf("error occured when attempting to connect to %s: %w", l.Addr().String(), err)
		}
		go l.Dispatch(conn)
	}
}

package proxy

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"sync"
)

// Wrapper over native Listener to add functionality
type httpListener struct {
	net.Listener
}

// Handles the lifecycle of all listeners
type listenerGroup struct {
	listeners []*httpListener
	wg        sync.WaitGroup
	// Used to spin off and start each listener
	startCh chan *httpListener
	// Used to notify back to start that a listener has started, to decrement counter
	notifyStartedCh chan struct{}
	// Notify the group of any errors
	errCh chan error
	// Dictates that the listenerGroup has been started and cannot be started again
	started bool
}

func (lg *listenerGroup) add(URL string) error {
	url, err := url.Parse(URL)
	if err != nil {
		return fmt.Errorf("bad upstream URL %w", err)
	}
	l, err := net.Listen("tcp", url.Host)
	if err != nil {
		return fmt.Errorf("unable to start listener at %s : %w", url.Host, err)
	}
	lg.listeners = append(lg.listeners, &httpListener{l})
	return nil
}

func (lg *listenerGroup) start() error {
	defer func() {
		close(lg.errCh)
		close(lg.notifyStartedCh)
		close(lg.startCh)
	}()
	// start has already been run, can't be done again
	if lg.started {
		return nil
	}
	lg.started = true
	// Starts the lifecycle of all the listeners
	go lg.reconcile()

	listenersCount := len(lg.listeners)
	for _, l := range lg.listeners {
		lg.startCh <- l
	}

	for {
		select {
		case err := <-lg.errCh:
			return err
		case <-lg.notifyStartedCh:
			listenersCount -= 1
			// Successfully started all of the listeners
			if listenersCount == 0 {
				return nil
			}
		}
	}
}

func (lg *listenerGroup) reconcile() {
	for l := range lg.startCh {
		lg.wg.Add(1)
		go func(l *httpListener) {
			// Notify back to the start method
			go func() {
				lg.notifyStartedCh <- struct{}{}
			}()
			defer lg.wg.Done()
			if err := l.start(); err != nil {
				// Ignore the error if it's just the listener closing
				if errors.Is(err, net.ErrClosed) {
					return
				}
				lg.errCh <- err
				return
			}
		}(l)
	}
}

func (lg *listenerGroup) stop() {
	for _, l := range lg.listeners {
		// Closing the listener will cause Accept to return an error
		// Eventually I want to find a more "defensive programming" approach
		// to safely exiting, but this works for now
		l.Close()
	}
	lg.wg.Wait()
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
		go handle(conn)
	}
}

// Absolutely temporary handling for now
func handle(c net.Conn) {
	defer c.Close()
	reader := bufio.NewReader(c)
	req, err := http.ReadRequest(reader)
	if err != nil {
		fmt.Println("failed to parse request:", err)
		return
	}

	targetURL, err := url.Parse("http://localhost:8080")
	if err != nil {
		fmt.Println("Bad upstream URL:", err)
		return
	}

	req.URL.Scheme = targetURL.Scheme
	req.URL.Host = targetURL.Host
	req.RequestURI = ""

	upstreamConn, err := net.Dial("tcp4", targetURL.Host)
	if err != nil {
		fmt.Println("failed to connect to upstream: ", err)
		return
	}
	defer upstreamConn.Close()

	err = req.Write(upstreamConn)
	if err != nil {
		fmt.Println("failed to write request upstream:", err)
		return
	}

	respReader := bufio.NewReader(upstreamConn)
	resp, err := http.ReadResponse(respReader, req)
	if err != nil {
		fmt.Println("failed to read upstream response: ", err)
		return
	}
	defer resp.Body.Close()

	err = resp.Write(c)
	if err != nil {
		fmt.Println("failed to write response to client: ", err)
		return
	}
}

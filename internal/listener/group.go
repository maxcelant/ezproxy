package listener

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"runtime"
	"sync"
)

// Handles the lifecycle of all listeners
type ListenerGroup struct {
	listeners []*httpListener
	// Unique socket addresses
	socketAddrs []string
	wg          sync.WaitGroup
	// Used to spin off and start each listener
	startCh chan *httpListener
	// Used to notify back to start that a listener has started, to decrement counter
	notifyStartedCh chan struct{}
	// Notify the group of any errors
	errCh chan error
	// Dictates that the listenerGroup has been started and cannot be started again
	started bool
}

func NewListenerGroup() *ListenerGroup {
	return &ListenerGroup{
		startCh:         make(chan *httpListener),
		errCh:           make(chan error),
		notifyStartedCh: make(chan struct{}),
		started:         false,
	}
}

func (lg *ListenerGroup) Add(URL string) error {
	sa, err := url.Parse(URL)
	if err != nil {
		return fmt.Errorf("bad upstream URL %w", err)
	}
	lg.socketAddrs = append(lg.socketAddrs, sa.Host)
	return nil
}

func (lg *ListenerGroup) Start() error {
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

	// Creates all the listeners
	if err := lg.listen(); err != nil {
		return fmt.Errorf("error occured while trying to start listeners: %w", err)
	}

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

func (lg *ListenerGroup) listen() error {
	// We start a number of listener goroutines equal to
	// (# of unique listeners) // CPU's for each unique socket address
	replicas := runtime.NumCPU() / len(lg.socketAddrs)

	fmt.Println("starting listeners")
	for _, sa := range lg.socketAddrs {
		for i := range replicas {
			fmt.Println("starting listener ", i)
			l, err := Listen("tcp", sa)
			if err != nil {
				return fmt.Errorf("unable to start listener at %s : %w", sa, err)
			}
			lg.listeners = append(lg.listeners, &httpListener{l})
		}
	}
	return nil
}

func (lg *ListenerGroup) reconcile() {
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

func (lg *ListenerGroup) Stop() {
	for _, l := range lg.listeners {
		// Closing the listener will cause Accept to return an error
		// Eventually I want to find a more "defensive programming" approach
		// to safely exiting, but this works for now
		l.Close()
	}
	lg.wg.Wait()
}

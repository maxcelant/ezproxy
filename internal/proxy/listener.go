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
	listeners []httpListener
	wg        sync.WaitGroup
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
	fmt.Printf("starting listener at %s\n", url.Host)
	lg.listeners = append(lg.listeners, httpListener{l})
	return nil
}

func (lg *listenerGroup) start() {
	for _, l := range lg.listeners {
		lg.wg.Add(1)
		l.start()
	}
}

func (lg *listenerGroup) stop() {
	for _, l := range lg.listeners {
		l.Close()
		lg.wg.Done()
	}
	lg.wg.Wait()
}

func (l *httpListener) start() {
	for {
		conn, err := l.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return
			}
			fmt.Printf("error occured when attempting to connect to %s", l.Addr().String())
			return
		}
		go handle(conn)
	}
}

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

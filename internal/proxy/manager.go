package proxy

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"sync"
)

type HTTPProxy struct {
	listeners []net.Listener
	endpoints []*url.URL
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
}

func NewProxyFromScratch() *HTTPProxy {
	return &HTTPProxy{}
}

func (p *HTTPProxy) AddListener(URL string) {
	url, err := url.Parse(URL)
	if err != nil {
		fmt.Println("bad upstream URL:", err)
		return
	}
	l, err := net.Listen("tcp", url.Host)
	if err != nil {
		fmt.Println("unable to start listener at ", url.Host)
		return
	}
	p.listeners = append(p.listeners, l)
}

func (p *HTTPProxy) AddEndpoint(URL string) {
	e, err := url.Parse(URL)
	if err != nil {
		fmt.Println("bad downstream URL:", err)
		return
	}
	p.endpoints = append(p.endpoints, e)
}

func (p *HTTPProxy) Start() {
	p.ctx, p.cancel = context.WithCancel(context.Background())
	p.startListeners()
	// TODO: Stop when all listeners are started
}

// Gracefully handle shutdown when sigterm signal is triggered
func (p *HTTPProxy) Stop() {
	fmt.Println("gracefully shutting down proxy...")
	// p.cancel()
	p.stopListeners()
	// This should block until all goroutines are cleaned up
	p.wg.Wait()
	fmt.Println("proxy shutdown complete.")
}

// TODO: Finish waitgroup for graceful shutdown
func (p *HTTPProxy) startListeners() {
	for _, l := range p.listeners {
		p.wg.Add(1)
		go p.startListener(p.ctx, l)
	}
}

func (p *HTTPProxy) stopListeners() {
	for _, l := range p.listeners {
		l.Close()
	}
}

func (p *HTTPProxy) startListener(ctx context.Context, listener net.Listener) {
	defer p.wg.Done()
	for {
		conn, err := listener.Accept()
		if err != nil {
			// closing the socket (part of shutdown sequence)
			if errors.Is(err, net.ErrClosed) {
				return
			}
			fmt.Printf("error occured when attempting to connect to %s", listener.Addr().String())
			return
		}
		go handle(conn)
	}
	// Decrement the wait group no matter what, so we aren't stuck in idle state
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

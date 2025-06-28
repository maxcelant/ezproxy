package proxy

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"sync"
)

type HTTPProxy struct {
	listeners []*url.URL
	endpoints []*url.URL
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
}

func NewProxyFromScratch() *HTTPProxy {
	return &HTTPProxy{}
}

func (p *HTTPProxy) AddListener(URL string) {
	l, err := url.Parse(URL)
	if err != nil {
		fmt.Println("Bad upstream URL:", err)
		return
	}
	p.listeners = append(p.listeners, l)
}

func (p *HTTPProxy) AddEndpoint(URL string) {
	e, err := url.Parse(URL)
	if err != nil {
		fmt.Println("Bad upstream URL:", err)
		return
	}
	p.endpoints = append(p.endpoints, e)
}

func (p *HTTPProxy) Start() {
	p.ctx, p.cancel = context.WithCancel(context.Background())
}

// Gracefully handle shutdown when sigterm signal is triggered
func (p *HTTPProxy) Stop() {
	p.cancel()
	// This should block until all goroutines are cleaned up
}

// TODO: Finish waitgroup for graceful shutdown
func (p HTTPProxy) startListeners() {
	for _, l := range p.listeners {
		p.wg.Add(1)
		go p.startListener(p.ctx, l)
	}
}

func (p HTTPProxy) startListener(ctx context.Context, url *url.URL) {
	listener, err := net.Listen("tcp4", url.Host)
	if err != nil {
		fmt.Printf("unable to create listener on %s: %s", url.Host, err)
		return
	}
	defer listener.Close()
	for {
		select {
		case <-ctx.Done():
			fmt.Println("shutting down listener ", url.Host)
			return
		default:
			conn, err := listener.Accept()
			if err != nil {
				fmt.Printf("error occured when attempting to connect to %s", listener.Addr().String())
				break
			}

			go handle(conn)
		}
	}
}

func handle(c net.Conn) {
	defer c.Close()
	reader := bufio.NewReader(c)
	req, err := http.ReadRequest(reader)
	if err != nil {
		fmt.Println("Failed to parse request:", err)
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
		fmt.Println("Failed to connect to upstream: ", err)
		return
	}
	defer upstreamConn.Close()

	err = req.Write(upstreamConn)
	if err != nil {
		fmt.Println("Failed to write request upstream:", err)
		return
	}

	respReader := bufio.NewReader(upstreamConn)
	resp, err := http.ReadResponse(respReader, req)
	if err != nil {
		fmt.Println("Failed to read upstream response: ", err)
		return
	}
	defer resp.Body.Close()

	err = resp.Write(c)
	if err != nil {
		fmt.Println("Failed to write response to client: ", err)
		return
	}
}

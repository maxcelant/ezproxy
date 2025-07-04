package workers

import (
	"bufio"
	"context"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/maxcelant/ezproxy/internal/dispatch"
)

type worker struct {
	ctx context.Context
}

func NewWorker(ctx context.Context) *worker {
	return &worker{ctx}
}

func (w *worker) start() error {
	fmt.Println("starting worker....")
	for {
		select {
		case <-w.ctx.Done():
			fmt.Println("shutting down worker")
			return nil
		default:
			time.Sleep(5 * time.Second)
		}
	}
}

func (w *worker) handle(ctx dispatch.DispatchContext) {
	conn, upstreams := ctx.Conn, ctx.Upstreams
	defer conn.Close()
	reader := bufio.NewReader(conn)
	req, err := http.ReadRequest(reader)
	if err != nil {
		fmt.Println("failed to parse request:", err)
		return
	}

	// Random loadbalancing (for testing)
	i := rand.Intn(len(upstreams))
	targetURL, err := url.Parse(upstreams[i])
	if err != nil {
		fmt.Println("Bad upstream URL:", err)
		return
	}

	req.URL.Scheme = targetURL.Scheme
	req.URL.Host = targetURL.Host
	req.RequestURI = ""
	req.Header.Add("x-forward-ezproxy", "true")

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

	resp.Header.Set("Host", targetURL.Host)

	err = resp.Write(conn)
	if err != nil {
		fmt.Println("failed to write response to client: ", err)
		return
	}
}

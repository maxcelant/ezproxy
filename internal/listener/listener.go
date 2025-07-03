package listener

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
)

// Wrapper over native Listener to add functionality
type httpListener struct {
	net.Listener
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

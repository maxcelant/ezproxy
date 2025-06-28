package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"net/url"
)

var (
	ip   = "127.0.0.1"
	port = "5000"
)

func main() {
	listener, err := net.Listen("tcp4", fmt.Sprintf("%s:%s", ip, port))
	if err != nil {
		fmt.Printf("unable to create listener on ip: %s with port: %s: %s", ip, port, err)
		return
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("error occured when attempting to connect to %s", listener.Addr().String())
			break
		}

		go handle(conn)
	}

	err = listener.Close()
	if err != nil {
		fmt.Printf("error occured closing the listener %s: %s", listener.Addr().String(), err)
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

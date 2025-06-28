package main

import (
	"errors"
	"fmt"
	"io"
	"net"
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

		go func(c net.Conn) {
			defer c.Close()
			var request string
			buf := make([]byte, 1024)
			for {
				n, err := c.Read(buf)
				if err != nil {
					if errors.Is(err, io.EOF) {
						return
					}
					fmt.Printf("error reading from buffer: %s", err)
					return
				}
				request += string(buf[:n])
			}
			parse(request)
		}(conn)
	}

	err = listener.Close()
	if err != nil {
		fmt.Printf("error occured closing the listener %s: %s", listener.Addr().String(), err)
	}
}

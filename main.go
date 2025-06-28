package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/maxcelant/ezproxy/internal/proxy"
)

var (
	ip   = "127.0.0.1"
	port = "5000"
)

func main() {
	proxy := proxy.NewProxyFromScratch()
	proxy.AddListener(fmt.Sprintf("%s:%s", ip, port))
	proxy.Start()

	time.Sleep(30 * time.Second)
	proxy.Stop()
}

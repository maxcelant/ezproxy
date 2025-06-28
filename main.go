package main

import (
	"fmt"
	"time"

	"github.com/maxcelant/ezproxy/internal/proxy"
)

var (
	ip   = "127.0.0.1"
	port = "5000"
)

func main() {
	host := fmt.Sprintf("http://%s:%s", ip, port)
	proxy := proxy.NewProxyFromScratch()

	proxy.AddListener(host)
	proxy.Start()

	time.Sleep(5 * time.Second)
	proxy.Stop()
}

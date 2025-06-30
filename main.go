package main

import (
	"time"

	"github.com/maxcelant/ezproxy/internal/proxy"
)

func main() {
	proxy := proxy.NewProxyFromScratch()

	proxy.AddListener("http://localhost:5000")
	proxy.AddListener("http://localhost:5001")
	proxy.Start()

	time.Sleep(30 * time.Second)
	proxy.Stop()
}

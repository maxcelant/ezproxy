package main

import (
	"fmt"
	"time"

	"github.com/maxcelant/ezproxy/internal/proxy"
)

func main() {
	var err error
	proxy := proxy.NewProxyFromScratch()

	if err = proxy.AddListener("http://localhost:5000"); err != nil {
		fmt.Println("error adding listener:", err)
		return
	}

	if err = proxy.AddListener("http://localhost:5001"); err != nil {
		fmt.Println("error adding listener:", err)
		return
	}

	if err = proxy.AddListener("http://localhost:5002"); err != nil {
		fmt.Println("error adding listener:", err)
		return
	}

	if err := proxy.Start(); err != nil {
		fmt.Println("error starting proxy:", err)
		proxy.Stop()
		return
	}

	time.Sleep(30 * time.Second)
	proxy.Stop()
}

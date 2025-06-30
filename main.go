package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/maxcelant/ezproxy/internal/proxy"
)

func main() {
	var err error
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	proxy := proxy.NewProxyFromScratch(log)

	if err = proxy.AddListener("http://localhost:5000"); err != nil {
		log.Error("error adding listener", "err", err)
		return
	}

	if err = proxy.AddListener("http://localhost:5001"); err != nil {
		log.Error("error adding listener", "err", err)
		return
	}

	if err = proxy.AddListener("http://localhost:5002"); err != nil {
		log.Error("error adding listener", "err", err)
		return
	}

	if err := proxy.Start(); err != nil {
		log.Error("error starting proxy", "err", err)
		proxy.Stop()
		return
	}

	time.Sleep(30 * time.Second)
	proxy.Stop()
}

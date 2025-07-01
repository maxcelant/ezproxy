package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/maxcelant/ezproxy/internal/chain"
	"github.com/maxcelant/ezproxy/internal/proxy"
)

func main() {
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	chain := chain.NewChain(
		chain.WithListener("http://localhost", 5000),
		chain.WithListener("http://localhost", 5001),
	)

	proxy := proxy.NewProxyFromScratch(log, proxy.WithChain(chain))

	if err := proxy.Start(); err != nil {
		log.Error("error starting proxy", "err", err)
		proxy.Stop()
		return
	}

	time.Sleep(30 * time.Second)
	proxy.Stop()
}

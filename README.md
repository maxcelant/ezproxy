# EzProxy

Still need to complete the worker goroutine functionality before it can be fully used.

### Basic Use

```go
func main() {
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	chain := chain.NewChain(
		chain.WithListener("http://localhost", 5000),
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
```

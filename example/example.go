package main

import (
	"fmt"
	"html"
	"net/http"
	"syscall"
	"time"

	"github.com/hongliang5316/graceful"
)

type handler struct {
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, port: %v, %q", r.Host, html.EscapeString(r.URL.Path))
}

func main() {
	server := graceful.NewServer(graceful.WithPIDFile("./graceful.pid"))
	server.Register(":9222", &handler{})
	err := server.Run()
	fmt.Printf("error: %v\n", err)
}

func listenMultiAddrs() {
	server := graceful.NewServer()
	server.Register("0.0.0.0:9223", &handler{})
	server.Register("0.0.0.0:9224", &handler{})
	server.Register("0.0.0.0:9225", &handler{})
	server.RegisterUnix("/tmp/test_graceful.sock", &handler{})
	err := server.Run()
	fmt.Printf("error: %v\n", err)
}

func callReload() {
	server := graceful.NewServer()
	server.Register("0.0.0.0:9226", &handler{})
	go func() {
		time.Sleep(time.Second)
		server.Reload()
	}()

	err := server.Run()
	fmt.Printf("error: %v\n", err)
}

func setReloadSignal() {
	server := graceful.NewServer(
		graceful.WithReloadSignals([]syscall.Signal{syscall.SIGUSR2}),
		graceful.WithStopSignals([]syscall.Signal{syscall.SIGINT}),
		graceful.WithStopTimeout(time.Minute),
		graceful.WithWatchInterval(10*time.Second),
	)
	server.Register("0.0.0.0:9226", &handler{})
	server.Run()
}

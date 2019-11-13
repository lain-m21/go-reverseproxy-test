package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"

	"github.com/lain-m21/go-reverseproxy-test/src/http"
)

const (
	exitOK = iota
	exitError

	serverAddr  string = ":9000"
	proxyAddr   string = ":9500"
	endpoint    string = "/say"
	backendHost string = "localhost" + serverAddr
	testFlags   string = "v1,v2,v3"
)

func main() {
	// Create separate main instead of doing the actual code here
	// since os.Exit can not handle `defer`. DON'T call `os.Exit` in the any other place.
	os.Exit(realMain(os.Args))
}

func realMain(_ []string) int {
	server := http.NewBackendServer(endpoint)

	transport := http.NewTransport()
	director := http.NewStreamDirector(backendHost, testFlags)
	proxy := http.NewProxy(endpoint, transport, director)

	lnServer, err := net.Listen("tcp", serverAddr)
	if err != nil {
		panic(fmt.Sprintf("failed to listen to the port %s", serverAddr))
	}

	lnProxy, err := net.Listen("tcp", proxyAddr)
	if err != nil {
		panic(fmt.Sprintf("failed to listen to the port %s", proxyAddr))
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg, ctx := errgroup.WithContext(ctx)

	wg.Go(func() error { return server.Serve(lnServer) })
	wg.Go(func() error { return proxy.Serve(lnProxy) })

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, os.Interrupt)
	select {
	case <-sigCh:
		fmt.Println("received SIGTERM, exiting servers")
	case <-ctx.Done():
		fmt.Println("context cancelled, exiting servers")
	}

	doneCh := make(chan struct{})

	go func() {
		defer close(doneCh)
		fmt.Println("shutdown a proxy server")
		proxy.Shutdown(ctx)
		fmt.Println("completed to shutdown a proxy server")

		fmt.Println("shutdown a backend server")
		server.Shutdown(ctx)
		fmt.Println("completed to shutdown a backend server")
	}()

	<-doneCh

	cancel()
	if err := wg.Wait(); err != nil {
		fmt.Println(errors.Wrap(err, "unhandled error received"))
		return exitError
	}
	return exitOK
}

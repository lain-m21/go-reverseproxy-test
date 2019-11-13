package http

import (
	"net/http"
	"net/http/httputil"
	"time"
)

func NewProxy(endpoint string, transport *http.Transport, director func(r *http.Request)) *http.Server {
	handler := &httputil.ReverseProxy{
		Director:  director,
		Transport: transport,
	}

	mux := http.NewServeMux()
	mux.Handle(endpoint, handler)
	server := &http.Server{
		Handler: mux,
	}
	return server
}

func NewStreamDirector(host string, flags string) func(r *http.Request) {
	return func(r *http.Request) {
		r.URL.Scheme = "http"
		r.URL.Host = host
		flags := flags
		r.Header.Set("flags", flags)
	}
}

func NewTransport() *http.Transport {
	return &http.Transport{
		MaxIdleConns:          200,
		MaxIdleConnsPerHost:   100,
		IdleConnTimeout:       120 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
}

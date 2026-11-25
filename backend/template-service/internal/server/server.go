// Package server builds the net/http server wrapping the Gin router.
package server

import (
	"net/http"
	"time"
)

// NewServer returns an *http.Server with the configured address, handler,
// and read/write timeouts applied. cmd/main.go owns starting and shutting
// it down — this just constructs it.
func NewServer(addr string, handler http.Handler, readTimeout, writeTimeout time.Duration) *http.Server {
	return &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}
}

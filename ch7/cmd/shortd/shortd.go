package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/tsingbx/effective-go/ch7/internal/httpio"
	"github.com/tsingbx/effective-go/ch7/short"
)

func main() {
	const (
		addr    = "localhost:8080"
		timeout = 10 * time.Second
	)
	logger := log.New(os.Stderr, "shortener: ", log.LstdFlags|log.Lmsgprefix)
	logger.Println("starting the server on", addr)

	shortener := short.NewServer()
	server := &http.Server{
		Addr:        addr,
		Handler:     http.TimeoutHandler(shortener, timeout, "timeout"),
		ReadTimeout: timeout,
	}
	if os.Getenv("LINKIT_DEBUG") == "1" {
		server.ErrorLog = logger
		server.Handler = httpio.LoggingMiddleware(server.Handler)
	}

	err := server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		logger.Println("server closed unexpectedly:", err)
	}
}

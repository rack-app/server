package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

func NewServer(handler http.Handler) *http.Server {
	return &http.Server{Addr: ":" + port(), Handler: handler}
}

type Server interface {
	ListenAndServe() error
	Shutdown(context.Context) error
}

func StartServer(s Server) {
	if err := s.ListenAndServe(); err != nil {
		if err != http.ErrServerClosed {
			panic(err)
		}
	}
}

func ShutdownServer(s Server, timeout time.Duration) {
	ctx, _ := context.WithTimeout(context.Background(), timeout)

	if err := s.Shutdown(ctx); err != nil {
		fmt.Println(err.Error())
	}
}

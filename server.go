package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"
)

func NewServer(handler http.Handler) *http.Server {
	return &http.Server{Addr: ":" + port(), Handler: handler}
}

func StartServer(s *http.Server) {
	if err := s.ListenAndServe(); err != nil {
		if err != http.ErrServerClosed {
			panic(err)
		}
	}
}

func ShutdownServer(s *http.Server, timeout time.Duration) {
	ctx, _ := context.WithTimeout(context.Background(), timeout)

	if err := s.Shutdown(ctx); err != nil {
		fmt.Println(err.Error())
	}
}

func port() string {
	var port string

	PortEnvValue, isExist := os.LookupEnv("PORT")

	if !isExist {
		port = "9292"
	} else {
		port = PortEnvValue
	}

	return port
}

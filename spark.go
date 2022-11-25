package spark

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"
)

// Options - spark options
type Options struct {
	Handler http.Handler

	//ServerAddres - 0.0.0.0:8080
	ServerAddres    string
	TimeoutShutdown time.Duration
	Ctx             context.Context
}

type spark struct {
	server          *http.Server
	timeoutShutdown time.Duration
	serverAddress   string
	ctx             context.Context
}

// ISpark - interface of spark interface
type ISpark interface {
	Shutdown() error
}

// New - new instance of spark server
func New(options Options) ISpark {
	timeout := 5 * time.Second
	if len(strings.TrimSpace(options.ServerAddres)) == 0 {
		options.ServerAddres = "0.0.0.0:8080"
	}
	server := &http.Server{
		Handler: options.Handler,
		Addr:    options.ServerAddres,
	}
	if options.TimeoutShutdown != 0 {
		timeout = options.TimeoutShutdown
	}
	spark := &spark{
		server:          server,
		timeoutShutdown: timeout,
		serverAddress:   options.ServerAddres,
	}

	if options.Ctx == nil {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		spark.ctx = ctx

	}
	go server.ListenAndServe()

	fmt.Println("started at : ", spark.serverAddress)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	return spark
}

// Shutdown - shutdown appliation
func (s *spark) Shutdown() error {
	if err := s.server.Shutdown(s.ctx); err != nil {
		fmt.Printf("Could not shut down server correctly: %v\n", err)
		os.Exit(1)
		return err
	}
	fmt.Println("shutdown successfully!")
	return nil
}

package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/scriptdealer/tcp-pow/internal/transport"
)

func main() {
	srv := transport.NewServer() // for the sake of brevity, dependenices are injected right inside the constructor

	interruption := make(chan os.Signal, 1)
	signal.Notify(interruption, os.Interrupt)

	// Start server in a goroutine
	go func() {
		if err := srv.Start(); err != nil {
			srv.Log.Error("running server", "error", err)
		}
	}()

	// Wait for interrupt signal
	<-interruption

	fmt.Println("\nShutting down gracefully...")
	if err := srv.Stop(); err != nil {
		srv.Log.Error("stopping server", "error", err)
	}
}

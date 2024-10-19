package test

import (
	"testing"
	"time"

	"github.com/scriptdealer/tcp-pow/internal/transport"
)

func TestClientServerIntegration(t *testing.T) {
	// Start the server
	server := transport.NewServer()
	go func() {
		if err := server.Start(); err != nil {
			t.Errorf("Failed to start server: %v", err)
		}
	}()
	defer server.Stop()

	// Give the server some time to start
	time.Sleep(time.Second)

	// Create a client
	client := transport.NewClient()

	// Test cases
	t.Run("GetQuote", func(t *testing.T) {
		quote, err := client.GetQuote()
		if err != nil {
			t.Errorf("Failed to get quote: %v", err)
		}
		if quote == "" {
			t.Error("Received empty quote")
		}
	})

	t.Run("MultipleQuotes", func(t *testing.T) {
		for i := 0; i < 5; i++ {
			quote, err := client.GetQuote()
			if err != nil {
				t.Errorf("Failed to get quote on iteration %d: %v", i, err)
			}
			if quote == "" {
				t.Errorf("Received empty quote on iteration %d", i)
			}
		}
	})

	t.Run("ConcurrentClients", func(t *testing.T) {
		numClients := 10
		done := make(chan bool)

		for i := 0; i < numClients; i++ {
			go func(clientNum int) {
				client := transport.NewClient()
				quote, err := client.GetQuote()
				if err != nil {
					t.Errorf("Client %d failed to get quote: %v", clientNum, err)
				}
				if quote == "" {
					t.Errorf("Client %d received empty quote", clientNum)
				}
				t.Log(quote)
				done <- true
			}(i)
		}

		for i := 0; i < numClients; i++ {
			<-done
		}
	})
}

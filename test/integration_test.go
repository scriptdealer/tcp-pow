package test

import (
	"sync"
	"testing"
	"time"

	"github.com/scriptdealer/tcp-pow/internal/transport"
)

func TestClientServerIntegration(t *testing.T) {
	server := transport.NewServer()
	errChan := make(chan error, 1)
	go func() {
		errChan <- server.Start()
	}()
	// Give the server a moment to start up
	time.Sleep(time.Second)

	select {
	case err := <-errChan:
		t.Fatalf("Failed to start server: %v", err)
	default:
		// Server started successfully, continue with tests
	}
	defer server.Stop()

	// Change maxCnn in NewConnectionPool from 5 to 6 to make complexity x10
	t.Run("SingleQuoteRequest", testSingleQuoteRequest)
	t.Run("MultipleQuoteRequests", testMultipleQuoteRequests)
	t.Run("ConcurrentClients", testConcurrentClients)
	t.Run("QuoteUniqueness", testQuoteUniqueness)
}

func testSingleQuoteRequest(t *testing.T) {
	client := transport.NewClient()
	quote, err := client.GetQuote()
	if err != nil {
		t.Fatalf("Failed to get quote: %v", err)
	}
	if quote == "" {
		t.Fatal("Received empty quote")
	}
}

func testMultipleQuoteRequests(t *testing.T) {
	start := time.Now()
	target := time.Millisecond * 900
	n := 6
	client := transport.NewClient()
	for i := 0; i < n; i++ {
		quote, err := client.GetQuote()
		if err != nil {
			t.Fatalf("Failed to get quote on iteration %d: %v", i, err)
		}
		if quote == "" {
			t.Fatalf("Received empty quote on iteration %d", i)
		}
	}
	if time.Since(start) < target {
		t.Fatalf("%d requests from same IP should take more than %v", n, target)
	}
}

func testConcurrentClients(t *testing.T) {
	numClients := 4 // Can't really test it from the same IP address
	var wg sync.WaitGroup
	wg.Add(numClients)

	for i := 0; i < numClients; i++ {
		go func() {
			defer wg.Done()
			client := transport.NewClient()
			quote, err := client.GetQuote()
			if err != nil {
				t.Errorf("Failed to get quote: %v", err)
			}
			if quote == "" {
				t.Error("Received empty quote")
			}
		}()
	}

	wg.Wait()
}

func testQuoteUniqueness(t *testing.T) {
	client := transport.NewClient()
	quotes := make(map[string]bool)
	numQuotes := 5

	for i := 0; i < numQuotes; i++ {
		quote, err := client.GetQuote()
		if err != nil {
			t.Fatalf("Failed to get quote on iteration %d: %v", i, err)
		}
		if quote == "" {
			t.Fatalf("Received empty quote on iteration %d", i)
		}
		quotes[quote] = true
	}

	if len(quotes) <= 1 {
		t.Fatal("Expected multiple unique quotes")
	}
}

package main

import (
	"github.com/scriptdealer/tcp-pow/internal/transport"
)

func main() {
	client := transport.NewClient() // for the sake of brevity, dependenices are injected right inside the constructor
	quote, err := client.GetQuote()
	if err != nil {
		client.Log.Error("Error getting quote", "reason", err)
	} else {
		client.Log.Info("Quote received", "quote", quote)
	}
}

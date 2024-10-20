package services

import (
	"fmt"
	"log/slog"
	"testing"
	"time"
)

func TestNewConnectionPool(t *testing.T) {
	cp := NewConnectionPool(slog.Default(), 100, 5, time.Minute)
	if cp.size != 100 {
		t.Errorf("Expected size 100, got %d", cp.size)
	}
	if cp.window != time.Minute {
		t.Errorf("Expected window 1 minute, got %v", cp.window)
	}
}

func TestCountRequests(t *testing.T) {
	cp := NewConnectionPool(slog.Default(), 10, 5, time.Minute)
	cp.maxCnn = 5

	// Test single IP
	for i := 1; i <= 6; i++ {
		count := cp.CountRequests("192.168.1.1")
		if i <= 5 && count != i {
			t.Errorf("Expected count %d, got %d", i, count)
		}
		if i > 5 && count != 5 {
			t.Errorf("Expected count to be capped at 5, got %d", count)
		}
	}

	// Test multiple IPs
	for i := 0; i < cp.size; i++ {
		ip := fmt.Sprintf("192.168.2.%d", i)
		count := cp.CountRequests(ip)
		if i < cp.size-1 {
			if count != 1 {
				t.Errorf("Expected count 1 for IP %s, got %d", ip, count)
			}
		} else {
			if count != ConnLimitReached {
				t.Errorf("Expected ConnLimitReached, got %d", count)
			}
		}
	}
}

func TestRemoveOldRecords(t *testing.T) {
	cp := NewConnectionPool(slog.Default(), 10, 5, time.Minute)

	// Add some records
	cp.CountRequests("192.168.1.1")
	cp.CountRequests("192.168.1.2")

	// Simulate time passing
	time.Sleep(time.Second)

	// Manually call removeOldRecords with a future time
	cp.removeOldRecords(time.Now().Add(2 * time.Minute))

	if len(cp.elMap) != 0 {
		t.Errorf("Expected all records to be removed, got %d", len(cp.elMap))
	}
}

func TestConcurrentAccess(t *testing.T) {
	cp := NewConnectionPool(slog.Default(), 100, 5, time.Minute)
	done := make(chan bool)

	for i := 0; i < 10; i++ {
		go func(id int) {
			for j := 0; j < 100; j++ {
				ip := fmt.Sprintf("192.168.1.%d", id)
				cp.CountRequests(ip)
			}
			done <- true
		}(i)
	}

	for i := 0; i < 10; i++ {
		<-done
	}

	if len(cp.elMap) != 10 {
		t.Errorf("Expected 10 unique IPs, got %d", len(cp.elMap))
	}
}

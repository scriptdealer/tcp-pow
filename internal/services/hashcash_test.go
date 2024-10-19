package services

import (
	"strings"
	"testing"
	"time"
)

func TestNewSimple(t *testing.T) {
	hc := NewHC(0, 0)
	if hc.difficulty != 7 {
		t.Errorf("Expected default difficulty to be 7, got %d", hc.difficulty)
	}
	if hc.timeout != time.Second*10 {
		t.Errorf("Expected default timeout to be 10 seconds, got %v", hc.timeout)
	}
}

func TestNewChallenge(t *testing.T) {
	hc := NewHC(0, 0)
	resourceID := "testResource"
	challenge := hc.NewChallenge(resourceID)
	parts := strings.Split(challenge, "|")
	if len(parts) != 3 {
		t.Errorf("Expected challenge to have 3 parts, got %d", len(parts))
	}
	if parts[0] != "7" {
		t.Errorf("Expected difficulty to be 7, got %s", parts[0])
	}
	if parts[1] != resourceID {
		t.Errorf("Expected resourceID to be %s, got %s", resourceID, parts[1])
	}
}

func TestSolveAndVerify(t *testing.T) {
	hc := NewHC(5, time.Second*10)
	challenge := hc.NewChallenge("testResource")
	solution := hc.Solve(challenge)
	if !hc.Verify(challenge, solution) {
		t.Errorf("Solution %s for challenge %s failed verification", solution, challenge)
	}
}

func TestVerifyInvalidSolution(t *testing.T) {
	hc := NewHC(0, 0)
	challenge := hc.NewChallenge("testResource")
	if hc.Verify(challenge, "invalidSolution") {
		t.Errorf("Invalid solution passed verification")
	}
}

func TestVerifyExpiredChallenge(t *testing.T) {
	hc := NewHC(0, 0)
	hc.timeout = time.Millisecond // Set a very short timeout for testing
	challenge := hc.NewChallenge("testResource")
	time.Sleep(time.Millisecond * 2) // Wait for the challenge to expire
	solution := hc.Solve(challenge)
	if hc.Verify(challenge, solution) {
		t.Errorf("Expired challenge passed verification")
	}
}

func TestParseInvalidChallenge(t *testing.T) {
	hc := NewHC(0, 0)
	_, _, _, err := hc.parse("invalid:challenge:format")
	if err == nil {
		t.Errorf("Expected error for invalid challenge format, got nil")
	}
}

func TestParseInvalidDifficulty(t *testing.T) {
	hc := NewHC(0, 0)
	_, _, _, err := hc.parse("0:resource:123456789")
	if err == nil {
		t.Errorf("Expected error for invalid difficulty, got nil")
	}
}

func TestParseInvalidTimestamp(t *testing.T) {
	hc := NewHC(0, 0)
	_, _, _, err := hc.parse("4:resource:-1")
	if err == nil {
		t.Errorf("Expected error for invalid timestamp, got nil")
	}
}

func TestPrintSolutionTimes(t *testing.T) {
	resourceID := "testResource"
	for difficulty := 2; difficulty < 7; difficulty++ {
		hc := NewHC(difficulty, time.Hour)
		challenge := hc.NewChallenge(resourceID)

		start := time.Now()
		solution := hc.Solve(challenge)
		duration := time.Since(start)

		t.Logf("Difficulty %d: Solution time %v, Solution: %s", difficulty, duration, solution)

		if !hc.Verify(challenge, solution) {
			t.Errorf("Solution verification failed for difficulty %d", difficulty)
		}
	}
}

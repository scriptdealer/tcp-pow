package services

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// https://github.com/umahmood/hashcash
// https://github.com/catalinc/hashcash

type BasicHashcash struct {
	difficulty int
	timeout    time.Duration
}

func NewHC(diffLevel int, timeout time.Duration) *BasicHashcash {
	if diffLevel < 1 || diffLevel > 64 {
		diffLevel = 7
	}
	if timeout < time.Second || timeout > time.Hour {
		timeout = time.Second * 10
	}

	return &BasicHashcash{
		difficulty: diffLevel,
		timeout:    timeout,
	}
}

func (hc *BasicHashcash) NewChallenge(resourceID string) string {
	return fmt.Sprintf("%d|%s|%d", hc.difficulty, resourceID, time.Now().UnixNano())
}

func (hc *BasicHashcash) NewChallengeWithDifficulty(resourceID string, difficulty int) string {
	return fmt.Sprintf("%d|%s|%d", difficulty, resourceID, time.Now().UnixNano())
}

func (hc *BasicHashcash) parse(challenge string) (int, string, int64, error) {
	parts := strings.Split(challenge, "|")
	if len(parts) != 3 {
		return 0, "", 0, errors.New("invalid challenge formatting")
	}
	difficulty, err := strconv.Atoi(parts[0])
	if err != nil || difficulty < 1 || difficulty > 64 {
		return 0, "", 0, fmt.Errorf("invalid difficulty level in challenge: %d", difficulty)
	}
	timestamp, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil || timestamp <= 0 || time.Now().UnixNano()-timestamp > int64(hc.timeout) {
		return 0, "", 0, fmt.Errorf("invalid or expired timestamp in challenge: %d", timestamp)
	}

	return difficulty, parts[1], timestamp, nil
}

func (hc *BasicHashcash) Solve(challenge string) string {
	difficulty, _, _, err := hc.parse(challenge)
	if err != nil {
		return err.Error()
	}
	nonce := 0
	for {
		attempt := fmt.Sprintf("%s%d", challenge, nonce)
		hash := sha256.Sum256([]byte(attempt))
		hashStr := hex.EncodeToString(hash[:])
		if strings.HasPrefix(hashStr, strings.Repeat("0", difficulty)) {
			return strconv.Itoa(nonce)
		}
		nonce++
	}
}

func (hc *BasicHashcash) Verify(challenge, solution string) bool {
	difficulty, _, timestamp, err := hc.parse(challenge)
	if err != nil || time.Now().UnixNano()-timestamp > int64(hc.timeout) {
		return false
	}
	hash := sha256.Sum256([]byte(challenge + solution))
	hashStr := hex.EncodeToString(hash[:])

	return strings.HasPrefix(hashStr, strings.Repeat("0", difficulty))
}

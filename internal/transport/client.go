package transport

import (
	"bufio"
	"fmt"
	"log/slog"
	"net"
	"strings"

	"github.com/scriptdealer/tcp-pow/internal/config"
	"github.com/scriptdealer/tcp-pow/internal/observability"
	"github.com/scriptdealer/tcp-pow/internal/services"
)

type Client struct {
	address string
	powAlg  *services.BasicHashcash
	Log     *slog.Logger
}

func NewClient() *Client {
	cfg := config.New()
	logger, _ := observability.NewLogger(slog.LevelDebug)

	return &Client{
		address: fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		powAlg:  services.NewHC(0, 0),
		Log:     logger,
	}
}

func (c *Client) GetQuote() (string, error) {
	c.Log.Info("connecting to server", "address", c.address)
	conn, err := net.Dial("tcp", c.address)
	if err != nil {
		return "", fmt.Errorf("connecting to server: %w", err)
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)
	challenge, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("reading challenge: %w", err)
	}

	challenge = strings.TrimSpace(challenge)
	nonce := c.powAlg.Solve(challenge)

	conn.Write([]byte(nonce + "\n")) //nolint:errcheck

	return reader.ReadString('\n')
}

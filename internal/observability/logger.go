package observability

import (
	"log/slog"
	"os"
)

func NewLogger(level slog.Leveler) (*slog.Logger, error) {
	var handler slog.Handler
	opts := &slog.HandlerOptions{
		AddSource: false,
		Level:     level,
	}
	handler = slog.NewJSONHandler(os.Stdout, opts)

	return slog.New(handler), nil
}

package logger

import (
	"context"
	"fmt"
	"log/slog"
)

const DefaultLogFile = "kubelet-credential-provider-vault.log"

type Logger interface {
	Log(ctx context.Context, level slog.Level, msg string, args ...any)
	Close() error
}

func ParseLogLevel(level string) (slog.Level, error) {
	switch level {
	case "debug":
		return slog.LevelDebug, nil
	case "info":
		return slog.LevelInfo, nil
	case "warn":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	default:
		return 0, fmt.Errorf("invalid log level: %s", level)
	}
}

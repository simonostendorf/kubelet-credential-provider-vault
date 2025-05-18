package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
)

type FileLogger struct {
	logger *slog.Logger
	file   *os.File
}

func NewFileLogger(enabled bool, logFile string, logLevel string) (Logger, error) {
	if enabled {
		// setup logger to log to file
		f, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644) //gosec:disable G302 G304
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %v", err)
		}

		// parse log level
		level, err := ParseLogLevel(logLevel)
		if err != nil {
			return nil, fmt.Errorf("failed to parse log level: %v", err)
		}

		// setup logger
		logger := slog.New(slog.NewJSONHandler(f, &slog.HandlerOptions{
			Level: level,
		}))

		return &FileLogger{
			logger: logger,
			file:   f,
		}, nil
	}
	return &FileLogger{
		logger: nil,
		file:   nil,
	}, nil
}

func (l *FileLogger) Close() error {
	if l.file != nil {
		if err := l.file.Close(); err != nil {
			return fmt.Errorf("failed to close log file: %v", err)
		}
		l.file = nil
	}
	if l.logger != nil {
		l.logger = nil
	}
	return nil
}

func (l *FileLogger) Log(ctx context.Context, level slog.Level, msg string, args ...any) {
	if l.logger != nil {
		l.logger.Log(ctx, level, msg, args...)
	}
}

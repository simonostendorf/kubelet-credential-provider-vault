package logger

import (
	"log/slog"
	"testing"
)

func TestParseLogLevel(t *testing.T) {
	tests := []struct {
		name           string
		logLevel       string
		expected       slog.Level
		expectedErrMsg string
	}{
		{
			name:     "debug",
			logLevel: "debug",
			expected: slog.LevelDebug,
		},
		{
			name:     "ingo",
			logLevel: "info",
			expected: slog.LevelInfo,
		},
		{
			name:     "warn",
			logLevel: "warn",
			expected: slog.LevelWarn,
		},
		{
			name:     "error",
			logLevel: "error",
			expected: slog.LevelError,
		},
		{
			name:           "invalid",
			logLevel:       "invalid",
			expected:       0,
			expectedErrMsg: "invalid log level: invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseLogLevel(tt.logLevel)
			if err != nil && err.Error() != tt.expectedErrMsg {
				t.Errorf("unexpected error: got %v, want %v", err, tt.expectedErrMsg)
			}
			if err == nil && result != tt.expected {
				t.Errorf("unexpected result: got %v, want %v", result, tt.expected)
			}
		})
	}
}

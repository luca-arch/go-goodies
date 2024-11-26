package logger_test

import (
	"log/slog"
	"testing"

	"github.com/luca-arch/go-goodies/logger"
	"github.com/stretchr/testify/assert"
)

// TestParseLevelString ensures that the ParseLevel function correctly parses string representations of log levels.
func TestParseLevelString(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		args string
		want slog.Level
	}{
		"Debug level - lowercase string": {
			args: "debug",
			want: slog.LevelDebug,
		},
		"Info level - lowercase string": {
			args: "info",
			want: slog.LevelInfo,
		},
		"Warn level - lowercase string": {
			args: "warn",
			want: slog.LevelWarn,
		},
		"Warn level - lowercase full string": {
			args: "warning",
			want: slog.LevelWarn,
		},
		"Error level - lowercase string": {
			args: "error",
			want: slog.LevelError,
		},
		"Debug level - uppercase string": {
			args: "DEBUG",
			want: slog.LevelDebug,
		},
		"Info level - uppercase string": {
			args: "INFO",
			want: slog.LevelInfo,
		},
		"Warn level - uppercase string": {
			args: "WARN",
			want: slog.LevelWarn,
		},
		"Warn level - uppercase full string": {
			args: "WARNING",
			want: slog.LevelWarn,
		},
		"Error level - uppercase string": {
			args: "ERROR",
			want: slog.LevelError,
		},
		"Unknown level": {
			args: "unknown",
			want: slog.LevelInfo,
		},
		"Empty string": {
			args: "",
			want: slog.LevelInfo,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			actual := logger.ParseLevel(test.args)
			assert.Equal(t, test.want, actual)
		})
	}
}

// TestParseLevel ensures that the ParseLevel function correctly keeps log levels.
func TestParseLevel(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		field slog.Level
	}{
		"Debug level": {
			field: slog.LevelDebug,
		},
		"Info level": {
			field: slog.LevelInfo,
		},
		"Warn level": {
			field: slog.LevelWarn,
		},
		"Error level": {
			field: slog.LevelError,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			actual := logger.ParseLevel(test.field)
			assert.Equal(t, test.field, actual)
		})
	}
}

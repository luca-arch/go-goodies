package logger

import (
	"io"
	"log/slog"
	"os"
	"strings"
)

// addSource causes the handler to compute the source code position of the log statement and add a SourceKey attribute to the output.
const addSource = false

type Lvl interface {
	slog.Level | string
}

// NewDev returns a text-formatted structured logger with debug level enabled.
func NewDev() *slog.Logger {
	return New(slog.LevelDebug, false)
}

// New returns a structured logger with the specified options.
func New[L Lvl](level L, json bool) *slog.Logger {
	lvl := new(slog.LevelVar)
	lvl.Set(ParseLevel(level))

	opts := &slog.HandlerOptions{
		AddSource:   addSource,
		Level:       lvl,
		ReplaceAttr: nil,
	}

	if json {
		return slog.New(slog.NewJSONHandler(os.Stdout, opts))
	}

	return slog.New(slog.NewTextHandler(os.Stdout, opts))
}

// NewNop returns a silent logger.
func NewNop() *slog.Logger {
	// TODO: replace with slog.DiscardHandler when available.
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

// ParseLevel returns the log level from the provided string/level.
func ParseLevel[L Lvl](level L) slog.Level {
	switch v := any(level).(type) {
	case slog.Level:
		return v
	case string:
		s := strings.ToLower(v)

		switch s {
		case "debug":
			return slog.LevelDebug
		case "warn", "warning":
			return slog.LevelWarn
		case "error":
			return slog.LevelError
		}
	}

	return slog.LevelInfo
}

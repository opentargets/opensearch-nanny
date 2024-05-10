// Package log provides a simple logging package.
package log

import (
	"log/slog"
	"os"

	"github.com/opentargets/opensearch-nanny/internal/config"
)

func levelFromString(l string) slog.Level {
	switch l {
	case "Debug":
		return slog.LevelDebug
	case "Info":
		return slog.LevelInfo
	case "Warn":
		return slog.LevelWarn
	case "Error":
		return slog.LevelError
	default:
		return slog.LevelDebug
	}
}

// InitLogger inits the logger with the configuration.
func InitLogger(sc config.ServerConfig) {
	ho := &slog.HandlerOptions{
		Level: levelFromString(sc.LogLevel),
	}

	var lh slog.Handler
	if sc.LogHandler == "json" {
		lh = slog.NewJSONHandler(os.Stdout, ho)
	} else if sc.LogHandler == "text" {
		lh = slog.NewTextHandler(os.Stdout, ho)
	}

	slog.SetDefault(slog.New(lh))
}

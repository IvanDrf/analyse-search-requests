package adapters

import (
	"log"
	"log/slog"
	"os"

	"github.com/IvanDrf/analyse-search-requests/internal/config"
)

func InitLogger(config *config.Config) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     selectLoggerLevel(config.App.LoggerLevel),
		AddSource: config.App.AddSource,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				t := a.Value.Time()
				a.Value = slog.StringValue(t.Format("2006-01-02-15:04"))
			}

			return a
		},
	}))

	slog.SetDefault(logger)
}

func selectLoggerLevel(level string) slog.Level {
	switch level {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError

	default:
		log.Fatalf("invalid logger level in config level=%s", level)
	}

	panic("")
}

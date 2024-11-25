package config

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

const (
	logFormatDev     = "dev"
	logFormatJSON    = "json"
	logFormatJSONGCP = "json-gcp"

	gcpLevelKey   = "severity"
	gcpSourceKey  = "logging.googleapis.com/sourceLocation"
	gcpMessageKey = "message"
)

func Die(msg string, args ...any) {
	slog.Error(msg, args...)
	os.Exit(1)
}

func DieLog(log *slog.Logger, msg string, args ...any) {
	log.Error(msg, args...)
	os.Exit(1)
}

func initLog(logLevel, logFormat string) {
	var level slog.Level
	if err := level.UnmarshalText([]byte(logLevel)); err != nil {
		Die(invalidConfigValue, "LOG_LEVEL", logLevel, "err", err)
	}

	var handler slog.Handler
	opts := slog.HandlerOptions{Level: level}

	switch logFormat {
	case logFormatJSON:
		handler = slog.NewJSONHandler(os.Stdout, &opts)
	case logFormatJSONGCP:
		opts.ReplaceAttr = replaceAttrGCP
		handler = slog.NewJSONHandler(os.Stdout, &opts)
	case logFormatDev:
		handler = tint.NewHandler(os.Stdout, &tint.Options{Level: level})
	default:
		handler = slog.NewTextHandler(os.Stdout, &opts)
	}

	slog.SetDefault(slog.New(handler))
}

func replaceAttrGCP(_ []string, attr slog.Attr) slog.Attr {
	switch attr.Key {
	case slog.LevelKey:
		attr.Key = gcpLevelKey
	case slog.SourceKey:
		attr.Key = gcpSourceKey
	case slog.MessageKey:
		attr.Key = gcpMessageKey
	}
	return attr
}

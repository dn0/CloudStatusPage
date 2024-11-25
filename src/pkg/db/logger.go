package db

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/tracelog"
)

const (
	logMessagePrepare = "Prepare"
	logMessagePrefix  = "pgx:"
)

type pgxLogger struct {
	logger *slog.Logger
	debug  bool
}

func newLogger(logger *slog.Logger, debug bool) pgx.QueryTracer {
	var level tracelog.LogLevel
	if debug {
		level = tracelog.LogLevelTrace
	} else {
		level = tracelog.LogLevelInfo
	}
	return &tracelog.TraceLog{
		LogLevel: level,
		Logger: &pgxLogger{
			logger: logger,
			debug:  debug,
		},
	}
}

func (l *pgxLogger) Log(ctx context.Context, level tracelog.LogLevel, msg string, data map[string]any) {
	if msg == logMessagePrepare && !l.debug {
		return
	}
	attrs := make([]slog.Attr, 0, len(data))
	for k, v := range data {
		if k == "time" {
			k = "took" // to avoid duplicate time field in the log
		}
		attrs = append(attrs, slog.Any(k, v))
	}
	l.logger.LogAttrs(ctx, translateLogLevel(level), logMessagePrefix+msg, attrs...)
}

func translateLogLevel(level tracelog.LogLevel) slog.Level {
	switch level {
	case tracelog.LogLevelTrace:
		return slog.LevelDebug
	case tracelog.LogLevelDebug:
		return slog.LevelDebug
	case tracelog.LogLevelInfo:
		return slog.LevelDebug
	case tracelog.LogLevelWarn:
		return slog.LevelInfo
	case tracelog.LogLevelError:
		return slog.LevelWarn
	case tracelog.LogLevelNone:
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

package db

import (
	"context"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"cspage/pkg/config"
)

const (
	initialPingTimeout = 3 * time.Second
)

type Conn interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type Client interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
	Ping(ctx context.Context) error
	Close()
}

type Tx = pgx.Tx

type Batch = pgx.Batch

var ErrTxClosed = pgx.ErrTxClosed

func newClient(ctx context.Context, name, connString string, debug bool) Client {
	if connString == "-" {
		return nil // We don't need this client.
	}

	var client Client
	if connString == "" {
		slog.Warn("Empty database URL => using dummy DB client", "db_conn", name)
		client = newDummyClient()
	} else {
		cfg, err := pgxpool.ParseConfig(connString)
		if err != nil {
			config.Die("Could not parse the DB config", "db_conn", name, "err", err)
		}
		cfg.ConnConfig.Tracer = newLogger(slog.Default(), debug)
		var errPool error
		client, errPool = pgxpool.NewWithConfig(ctx, cfg)
		if errPool != nil {
			config.Die("Could not create a DB pool", "db_conn", name, "err", errPool)
		}
	}

	tctx, cancel := context.WithTimeout(ctx, initialPingTimeout)
	defer cancel()
	if err := client.Ping(tctx); err != nil {
		config.Die("Could not connect to DB", "db_conn", name, "err", err)
	}

	return client
}

func closeClient(name string, client Client) {
	if client == nil {
		return
	}

	slog.Debug("Closing DB client...", "db_conn", name)
	client.Close()
}

package main

import (
	"context"
	"log/slog"

	"cspage/pkg/db"
	"cspage/pkg/http"
	"cspage/pkg/mon/web"
	"cspage/pkg/worker"
)

const svcName = "mon-web"

func main() {
	ctx := context.Background()
	cfg := web.NewConfig()
	dbClients := db.NewClients(ctx, &cfg.DatabaseConfig, cfg.Debug)

	worker.Run(
		ctx,
		slog.With("name", svcName, "env", cfg.Env),
		func() {
			db.CloseClients(dbClients)
		},
		http.NewServer(ctx, &cfg.HTTPConfig, cfg.Debug, web.Handlers(cfg, dbClients)),
	)
}

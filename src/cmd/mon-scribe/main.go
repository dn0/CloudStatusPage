package main

import (
	"context"
	"log/slog"

	"cspage/pkg/db"
	"cspage/pkg/http"
	"cspage/pkg/mon/scribe"
	"cspage/pkg/msg"
	"cspage/pkg/worker"
)

const svcName = "mon-scribe"

func main() {
	ctx := context.Background()
	cfg := scribe.NewConfig()
	dbClients := db.NewClients(ctx, &cfg.DatabaseConfig, cfg.Debug)
	pubsubClient := msg.NewPubsubClient(ctx, cfg.PubsubProjectID, msg.NewPubsubConfig(nil, &cfg.PubsubSubscriberConfig))

	workers := []worker.Worker{
		http.NewSimpleServer(ctx, &cfg.HTTPConfig, cfg.Debug),
	}
	for _, id := range cfg.Env.PubsubSubscriptions {
		workers = append(workers, scribe.NewConsumer(ctx, cfg, dbClients, msg.NewPubsubSubscription(pubsubClient, id)))
	}

	worker.Run(
		ctx,
		slog.With("name", svcName, "env", cfg.Env),
		func() {
			msg.ClosePubsubClient(pubsubClient)
			db.CloseClients(dbClients)
		},
		workers...,
	)
}

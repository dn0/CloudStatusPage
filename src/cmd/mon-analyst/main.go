package main

import (
	"context"
	"log/slog"

	"cspage/pkg/db"
	"cspage/pkg/http"
	"cspage/pkg/mon/analyst"
	"cspage/pkg/mon/speaker"
	"cspage/pkg/msg"
	"cspage/pkg/worker"
)

const svcName = "mon-analyst"

func main() {
	ctx := context.Background()
	cfg := analyst.NewConfig()
	speakerCfg := speaker.NewConfigFromAnalyst(cfg)
	dbClients := db.NewClients(ctx, &cfg.DatabaseConfig, cfg.Debug)
	pubsubClient := msg.NewPubsubClient(ctx, cfg.PubsubProjectID, msg.NewPubsubConfig(&cfg.PubsubPublisherConfig, nil))
	pubsubAlertTopic := msg.NewPubsubTopic(pubsubClient, cfg.PubsubAlertTopic)

	workers := []worker.Worker{
		http.NewSimpleServer(ctx, &cfg.HTTPConfig, cfg.Debug),
	}
	for _, cloud := range cfg.Clouds {
		incidentChan := speaker.NewIncidentChannel()
		workers = append(
			workers,
			analyst.NewIncidentTicker(ctx, cfg, dbClients, cloud, pubsubAlertTopic, incidentChan),
			speaker.NewIncidentReceiver(ctx, speakerCfg, dbClients, cloud, incidentChan),
			analyst.NewPingTicker(ctx, cfg, dbClients, cloud, pubsubAlertTopic),
		)
		for _, w := range analyst.NewProbeTickers(ctx, cfg, dbClients, cloud, pubsubAlertTopic) {
			workers = append(workers, w)
		}
	}

	worker.Run(
		ctx,
		slog.With("name", svcName, "env", cfg.Env),
		func() {
			msg.ClosePubsubTopic(pubsubAlertTopic)
			msg.ClosePubsubClient(pubsubClient)
			db.CloseClients(dbClients)
		},
		workers...,
	)
}

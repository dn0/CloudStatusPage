package scribe

import (
	"context"

	"cspage/pkg/db"
	"cspage/pkg/msg"
	"cspage/pkg/worker"
)

func NewConsumer(
	ctx context.Context,
	cfg *Config,
	dbc *db.Clients,
	sub *msg.PubsubSubscription,
) *worker.Daemon {
	j := &job{
		cfg:  cfg,
		name: sub.ID,
		dbc:  dbc,
	}
	return worker.NewConsumer(ctx, &cfg.BaseConfig, j, sub)
}

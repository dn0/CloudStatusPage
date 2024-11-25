package speaker

import (
	"context"
	"log/slog"

	"cspage/pkg/data"
	"cspage/pkg/db"
)

type incidentJob struct {
	cfg     *Config
	dbc     *db.Clients
	log     *slog.Logger
	cloud   string
	name    string
	twitter *twitterClient
}

func (j *incidentJob) String() string {
	return j.name
}

func (j *incidentJob) Enabled() bool {
	return j.cfg.TwitterAPIKey != "" && j.cfg.TwitterAPISecret != ""
}

func (j *incidentJob) PreStart(_ context.Context) {}

func (j *incidentJob) Start(_ context.Context) {}

func (j *incidentJob) Stop(_ context.Context) {}

func (j *incidentJob) Shutdown(_ error) {}

func (j *incidentJob) Process(ctx context.Context, inc *data.Incident) error {
	res, err := j.twitter.createTweetFromIncident(ctx, inc)
	if err != nil {
		j.log.Error("Could not create tweet from incident", "incident", inc, "err", err)
		return err
	}

	if res != nil {
		j.log.Info("Created tweet from incident", "incident", inc, "tweet", res)
		if j.cfg.DryRun {
			return nil
		}
		inc.Data.TwitterId = res.Data.ID
		if err = data.UpdateIncidentData(ctx, j.dbc.Write, j.cloud, inc); err != nil {
			j.log.Error("Could not save tweet ID to incident data", "incident", inc, "tweet", res, "err", err)
		}
	}

	return nil
}
